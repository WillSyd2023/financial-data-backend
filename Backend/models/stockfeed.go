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
	Id primitive.ObjectID `bson:"_id,omitempty"`
}
