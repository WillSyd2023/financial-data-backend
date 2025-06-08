package util

import (
	"Backend/constant"
	"Backend/dto"
	"time"
)

type DateGenerator time.Time

func NewDateGenerator(date string) DateGenerator {
	timeValue, _ := time.Parse(constant.LayoutISO, date)
	return DateGenerator(timeValue)
}

func (d *DateGenerator) Current() dto.DateOnly {
	return dto.DateOnly(*d)
}

func (d *DateGenerator) Next() dto.DateOnly {
	return (*d).Current().AddDate(0, 0, 1)
}
