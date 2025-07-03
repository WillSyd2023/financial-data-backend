package repo

import (
	"Backend/configs"
	"Backend/dto"
	"Backend/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RepoItf interface {
	CheckSymbolExists(*gin.Context, *dto.CollectSymbolReq) (bool, error)
	InsertNewSymbolData(*gin.Context, *dto.DataPerSymbol) error
	DeleteSymbol(*gin.Context, *dto.DeleteSymbolReq) error
	StoredData(*gin.Context) ([]dto.DataPerSymbol, error)
}

type Repo struct {
	symbolCollection *mongo.Collection
	ohlcvCollection  *mongo.Collection
}

func NewRepo() *Repo {
	return &Repo{
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
	c := ctx.Request.Context()

	// Insert new symbol and last-refreshed date
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

func (rp *Repo) DeleteSymbol(ctx *gin.Context, req *dto.DeleteSymbolReq) error {
	c := ctx.Request.Context()
	if _, err := rp.symbolCollection.DeleteOne(c, bson.M{"name": bson.M{"$eq": req.Symbol}}); err != nil {
		return err
	}
	_, err := rp.ohlcvCollection.DeleteMany(c, bson.M{"ticker": bson.M{"$eq": req.Symbol}})
	return err
}

func (rp *Repo) StoredData(ctx *gin.Context) ([]dto.DataPerSymbol, error) {
	c := ctx.Request.Context()

	results, err := rp.symbolCollection.Find(c, bson.D{}, options.Find().SetSort(
		bson.D{{Key: "name", Value: 1}}))
	if err != nil {
		return nil, err
	}

	data := make([]dto.DataPerSymbol, 0)
	defer results.Close(c)
	for results.Next(c) {
		var symbol models.Symbol
		if err = results.Decode(&symbol); err != nil {
			return nil, err
		}
		data = append(data, dto.DataPerSymbol{
			MetaData: &dto.SymbolDataMeta{
				Symbol:        symbol.Name,
				LastRefreshed: dto.DateOnly(symbol.LastRefreshed)},
		})
	}

	results, err = rp.ohlcvCollection.Find(c, bson.D{}, options.Find().SetSort(
		bson.D{{Key: "ticker", Value: 1}, {Key: "date", Value: 1}}))
	if err != nil {
		return nil, err
	}

	ix := 0
	defer results.Close(c)
	for results.Next(c) {
		var ohlcv models.DailyOHLCV
		if err = results.Decode(&ohlcv); err != nil {
			return nil, err
		}

		open, err := decimal.NewFromString(ohlcv.OpenPrice.String())
		if err != nil {
			return nil, err
		}
		high, err := decimal.NewFromString(ohlcv.HighPrice.String())
		if err != nil {
			return nil, err
		}
		low, err := decimal.NewFromString(ohlcv.LowPrice.String())
		if err != nil {
			return nil, err
		}
		close, err := decimal.NewFromString(ohlcv.ClosePrice.String())
		if err != nil {
			return nil, err
		}

		res := dto.DailyOHLCVRes{
			Day: dto.DateOnly(ohlcv.Date),
			OHLC: map[string]decimal.Decimal{
				"open":  open,
				"high":  high,
				"low":   low,
				"close": close,
			},
			Volume: int(ohlcv.Volume),
		}

		if data[ix].MetaData.Symbol != ohlcv.Ticker {
			ix++
		}

		data[ix].TimeSeries = append(data[ix].TimeSeries, res)
	}

	// Remember to figure out number of data for each stock
	for _, datum := range data {
		datum.MetaData.Size = len(datum.TimeSeries)
	}

	return data, nil
}
