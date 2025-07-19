package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Symbol struct {
	Id            primitive.ObjectID `bson:"_id,omitempty"`
	Name          string             `bson:"name"`
	LastRefreshed time.Time          `bson:"last_refreshed"`
}

type DailyOHLCV struct {
	Id         primitive.ObjectID   `bson:"_id,omitempty"`
	Date       time.Time            `bson:"date"`
	Ticker     string               `bson:"ticker"`
	OpenPrice  primitive.Decimal128 `bson:"open_price"`
	HighPrice  primitive.Decimal128 `bson:"high_price"`
	LowPrice   primitive.Decimal128 `bson:"low_price"`
	ClosePrice primitive.Decimal128 `bson:"close_price"`
	Volume     int64                `bson:"volume"`
}
