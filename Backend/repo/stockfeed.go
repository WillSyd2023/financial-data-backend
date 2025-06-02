package repo

import (
	"Backend/dto"
	"database/sql"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type RepoItf interface {
	CheckSymbolExists(*gin.Context, *dto.CollectSymbolReq) (bool, error)
	InsertNewSymbolData(*gin.Context, *dto.StockDataRes) error
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

func (rp *Repo) InsertNewSymbolData(ctx *gin.Context, stockData *dto.StockDataRes) error {
	// Insert new symbol and last-refreshed data
	var id int
	err := rp.db.QueryRowContext(
		ctx,
		"INSERT INTO symbols (symbol, last_refreshed) VALUES ($1, $2) RETURNING symbol_id",
		stockData.MetaData.Symbol,
		stockData.MetaData.LastRefreshed.Format(time.RFC3339),
	).Scan(&id)
	if err != nil {
		return err
	}

	// Data to insert and numbering on SQL code
	data := make([]any, 0)
	pos := 1
	query := "INSERT INTO ohlcv_per_day " +
		"(record_day, open_price, high_price, low_price, close_price, volume, symbol_id) " +
		"VALUES "
	for i, ohlcv := range stockData.TimeSeries {
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
		data = append(data,
			ohlcv.Day.Format(time.RFC3339), ohlcv.OHLC["open"], ohlcv.OHLC["high"], ohlcv.OHLC["low"],
			ohlcv.OHLC["close"], ohlcv.Volume, id)
	}

	// Insert data
	_, err = rp.db.ExecContext(ctx, query, data...)
	return err
}
