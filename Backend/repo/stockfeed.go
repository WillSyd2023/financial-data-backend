package repo

import (
	"Backend/configs"
	"Backend/dto"
	"Backend/models"
	"database/sql"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RepoItf interface {
	CheckSymbolExists(*gin.Context, *dto.CollectSymbolReq) (bool, error)
	InsertNewSymbolData(*gin.Context, *dto.DataPerSymbol) error
	DeleteSymbol(*gin.Context, *dto.DeleteSymbolReq) error
	StoredData(*gin.Context) ([]dto.DataPerSymbol, error)
}

type Repo struct {
	db               *sql.DB
	symbolCollection *mongo.Collection
	ohlcvCollection  *mongo.Collection
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{
		db:               db,
		symbolCollection: configs.GetCollection(configs.DB, "symbols"),
		ohlcvCollection:  configs.GetCollection(configs.DB, "daily_ohlcv"),
	}
}

func (rp *Repo) CheckSymbolExists(ctx *gin.Context, req *dto.CollectSymbolReq) (bool, error) {
	c := ctx.Request.Context()

	filter := bson.M{"symbol": req.Symbol}
	err := rp.symbolCollection.FindOne(c, filter).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (rp *Repo) InsertNewSymbolData(ctx *gin.Context, data *dto.DataPerSymbol) error {
	// Insert new symbol and last-refreshed date
	c := ctx.Request.Context()
	if _, err := rp.symbolCollection.InsertOne(c, models.Symbol{
		Id:            primitive.NewObjectID(),
		Name:          data.MetaData.Symbol,
		LastRefreshed: time.Time(data.MetaData.LastRefreshed),
	}); err != nil {
		return err
	}

	// Insert time-series data
	timeSeries := make([]any, len(data.TimeSeries))
	for i, ohlcv := range data.TimeSeries {
		openPrice, err := primitive.ParseDecimal128(ohlcv.OHLC["open"].String())
		if err != nil {
			return err
		}
		highPrice, err := primitive.ParseDecimal128(ohlcv.OHLC["high"].String())
		if err != nil {
			return err
		}
		lowPrice, err := primitive.ParseDecimal128(ohlcv.OHLC["low"].String())
		if err != nil {
			return err
		}
		closePrice, err := primitive.ParseDecimal128(ohlcv.OHLC["close"].String())
		if err != nil {
			return err
		}
		timeSeries[i] = models.DailyOHLCV{
			Date:       time.Time(ohlcv.Day),
			Ticker:     data.MetaData.Symbol,
			OpenPrice:  openPrice,
			HighPrice:  highPrice,
			LowPrice:   lowPrice,
			ClosePrice: closePrice,
			Volume:     int64(ohlcv.Volume),
		}
	}

	_, err := rp.ohlcvCollection.InsertMany(ctx, timeSeries)
	return err
}

func (rp *Repo) InsertNewSymbolDataPostgres(ctx *gin.Context, data *dto.DataPerSymbol) error {
	// Insert new symbol and last-refreshed data
	var id int
	err := rp.db.QueryRowContext(
		ctx,
		"INSERT INTO symbols (symbol, last_refreshed) VALUES ($1, $2) RETURNING symbol_id",
		data.MetaData.Symbol,
		time.Time(data.MetaData.LastRefreshed).Format(time.RFC3339),
	).Scan(&id)
	if err != nil {
		return err
	}

	// Data to insert and numbering on SQL code
	timeSeries := make([]any, 0)
	pos := 1
	query := "INSERT INTO ohlcv_per_day " +
		"(record_day, open_price, high_price, low_price, close_price, volume, symbol_id) " +
		"VALUES "
	for i, ohlcv := range data.TimeSeries {
		// Comma for SQL syntax
		if i != 0 {
			query += ", "
		}

		// Numbering for SQL code
		query += fmt.Sprintf(
			"($%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			pos, pos+1, pos+2, pos+3, pos+4, pos+5, pos+6)
		pos += 7

		// Data corresponding to numbering
		timeSeries = append(timeSeries,
			time.Time(ohlcv.Day).Format(time.RFC3339),
			ohlcv.OHLC["open"], ohlcv.OHLC["high"], ohlcv.OHLC["low"], ohlcv.OHLC["close"],
			ohlcv.Volume, id)
	}

	// Insert data
	_, err = rp.db.ExecContext(ctx, query, timeSeries...)
	return err
}

func (rp *Repo) DeleteSymbol(ctx *gin.Context, req *dto.DeleteSymbolReq) error {
	res, err := rp.db.ExecContext(ctx,
		"DELETE FROM symbols WHERE symbol=$1", req.Symbol)
	if err == nil {
		count, err := res.RowsAffected()
		if err == nil && count > 0 {
			return nil
		}
	}
	return err
}

func (rp *Repo) StoredData(ctx *gin.Context) ([]dto.DataPerSymbol, error) {
	query := "SELECT " +
		"symbol, last_refreshed, record_day, open_price, high_price, low_price, close_price, volume " +
		"FROM symbols INNER JOIN ohlcv_per_day ON symbols.symbol_id = ohlcv_per_day.symbol_id " +
		"ORDER BY symbol, record_day ASC"
	rows, err := rp.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := make([]dto.DataPerSymbol, 0)
	ix := -1
	for rows.Next() {
		meta := dto.SymbolDataMeta{}
		ohlcv := dto.DailyOHLCVRes{}
		var open, high, low, close decimal.Decimal
		if err := rows.Scan(
			&meta.Symbol, &meta.LastRefreshed, &ohlcv.Day,
			&open, &high, &low, &close,
			&ohlcv.Volume,
		); err != nil {
			return nil, err
		}
		ohlcv.OHLC = map[string]decimal.Decimal{
			"open":  open,
			"high":  high,
			"low":   low,
			"close": close,
		}

		if len(data) == 0 ||
			data[ix].MetaData.Symbol != meta.Symbol {
			ix++
			data = append(data, dto.DataPerSymbol{
				MetaData:   &meta,
				TimeSeries: []dto.DailyOHLCVRes{ohlcv},
			})
		} else {
			data[ix].TimeSeries = append(
				data[ix].TimeSeries,
				ohlcv,
			)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Remember to figure out number of data for each stock
	for _, datum := range data {
		datum.MetaData.Size = len(datum.TimeSeries)
	}

	return data, nil
}
