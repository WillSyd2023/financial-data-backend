package repo

import (
	"Backend/dto"
	"database/sql"

	"github.com/gin-gonic/gin"
)

type RepoItf interface {
	CheckSymbolExists(*gin.Context, *dto.CollectSymbolReq) (int, error)
	InsertNewSymbolData(*gin.Context, int, *dto.StockDataRes) error
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

func (rp *Repo) InsertNewSymbolData(ctx *gin.Context, id int, stockData *dto.StockDataRes) error {
	return nil
}
