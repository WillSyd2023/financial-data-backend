package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Symbol struct {
	Id            primitive.ObjectID `json:"id"`
	Name          string             `json:"name"`
	LastRefreshed time.Time          `json:"last_refreshed"`
}
