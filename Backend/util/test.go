package util

import (
	"Backend/constant"
	"Backend/dto"
	"time"

	"github.com/shopspring/decimal"
)

// Date generator object to instantly generate 'dummy' dates
type DateGenerator time.Time

func NewDateGenerator(date string) *DateGenerator {
	timeValue, _ := time.Parse(constant.LayoutISO, date)
	dateGen := DateGenerator(timeValue)
	return &dateGen
}

func (d *DateGenerator) Current() dto.DateOnly {
	return dto.DateOnly(*d)
}

func (d *DateGenerator) Next() dto.DateOnly {
	*d = DateGenerator(d.Current().AddDate(0, 0, 1))
	return d.Current()
}

// OHLCV generator object to instantly generate 'dummy' data
type OHLCVGenerator struct {
	DateGen *DateGenerator
	Value   int64
	Volume  int
}

func NewOHLCVGenerator(dateGen *DateGenerator, value, volume int) *OHLCVGenerator {
	return &OHLCVGenerator{
		DateGen: dateGen,
		Value:   int64(value),
		Volume:  volume,
	}
}

func (o *OHLCVGenerator) Next() dto.DailyOHLCVRes {
	var res dto.DailyOHLCVRes

	res.Day = o.DateGen.Next()

	res.Volume = o.Volume
	o.Volume++

	res.OHLC = make(map[string]decimal.Decimal)
	res.OHLC["open"] = decimal.NewFromInt(o.Value + 1)
	res.OHLC["high"] = decimal.NewFromInt(o.Value + 2)
	res.OHLC["low"] = decimal.NewFromInt(o.Value + 3)
	res.OHLC["close"] = decimal.NewFromInt(o.Value + 4)
	o.Value += 100

	return res
}
