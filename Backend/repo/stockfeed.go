package repo

import (
	"Backend/dto"
	"database/sql"

	"github.com/gin-gonic/gin"
)

type RepoItf interface {
	CheckSymbolExists(*gin.Context, *dto.CollectSymbolReq) (bool, error)
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
	err := rp.db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM symbols WHERE symbol = $1);",
		req.Symbol).Scan(&exists)
	return exists, err
}
