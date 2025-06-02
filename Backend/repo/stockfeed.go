package repo

import (
	"Backend/dto"
	"database/sql"

	"github.com/gin-gonic/gin"
)

type RepoItf interface {
	CheckSymbolExists(*gin.Context, *dto.CollectSymbolReq) (int, error)
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

func (rp *Repo) CheckSymbolExists(ctx *gin.Context, req *dto.CollectSymbolReq) (int, error) {
	var id int
	err := rp.db.QueryRowContext(
		ctx,
		"SELECT symbol_id FROM symbols WHERE symbol = $1",
		req.Symbol).Scan(&id)
	return id, err
}

func (rp *Repo) InsertNewSymbolData(ctx *gin.Context, stockData *dto.StockDataRes) error {
	// Insert new symbol and last-refreshed data
	var id int
	err := rp.db.QueryRowContext(
		ctx,
		"INSERT INTO symbols (symbol, last_refreshed) VALUES ($1, $2) RETURNING symbol_id",
		stockData.MetaData.Symbol,
		stockData.MetaData.LastRefreshed,
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

	return nil
}
