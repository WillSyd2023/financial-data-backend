package repo

import (
	"Backend/dto"
	"database/sql"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type RepoItf interface {
	CheckSymbolExists(*gin.Context, *dto.CollectSymbolReq) (bool, error)
	InsertNewSymbolData(*gin.Context, *dto.DataPerSymbol) error
	DeleteSymbol(*gin.Context, *dto.DeleteSymbolReq) error
	StoredData(*gin.Context) ([]dto.DataPerSymbol, error)
}

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (rp *Repo) CheckSymbolExists(ctx *gin.Context, req *dto.CollectSymbolReq) (bool, error) {
	var exists bool
	err := rp.db.QueryRowContext(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM symbols WHERE symbol = $1)",
		req.Symbol).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (rp *Repo) InsertNewSymbolData(ctx *gin.Context, data *dto.DataPerSymbol) error {
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
