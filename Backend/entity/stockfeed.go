package entity

import (
	"encoding/json"
	"time"
)

type Symbol struct {
	Id            int
	Name          string
	LastRefreshed time.Time
}

type DateOnly time.Time

func (d DateOnly) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(d).Format("2006-01-02"))
}

func (d DateOnly) Weekday() time.Weekday {
	return time.Time(d).Weekday()
}

func (d DateOnly) AddDate(years, months, days int) DateOnly {
	return DateOnly(time.Time(d).AddDate(years, months, days))
}

func (d DateOnly) Before(e DateOnly) bool {
	return time.Time(d).Before(time.Time(e))
}

func (d DateOnly) After(e DateOnly) bool {
	return time.Time(d).After(time.Time(e))
}
