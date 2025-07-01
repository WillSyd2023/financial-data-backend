package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Symbol struct {
	Id            primitive.ObjectID `json:"id,omitempty"`
	Name          string             `json:"name,omitempty"`
	LastRefreshed time.Time          `json:"last_refreshed,omitempty"`
}
