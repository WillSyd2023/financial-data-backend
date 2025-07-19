package entity

import (
	"time"
)

type Symbol struct {
	Id            int
	Name          string
	LastRefreshed time.Time
}
