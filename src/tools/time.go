package tools

import (
	"context"
	"go.uber.org/zap"
	"time"
)

// 通用时间格式模板
var TnoLayout = "20060102150405"
var TimeLayout = "2006-01-02 15:04:05"
var DateLayout = "2006-01-02"
var DateTimeLayout = "2006-01-02 15:04:05"
var DateCompactLayout = "20060102"
var DateTimeCompactLayout = "20060102150405"

func GetCurrentDate() string {
	return time.Now().Format(TimeLayout)
}

func GetCurrentUnix() int64 {
	return time.Now().Unix()
}

func GetCurrentMilliUnix() int64 {
	return time.Now().UnixNano() / 1000000
}

func GetCurrentNanoUnix() int64 {
	return time.Now().UnixNano()
}

func GetDateRange(rawTime time.Time) (time.Time, time.Time) {
	sTime, err := time.ParseInLocation(DateLayout, rawTime.Format(DateLayout), time.Local)
	if err != nil {
		logs.Error(nil, "GetDateRange", zap.Any("rawTime", rawTime), zap.Error(err))
	}
	eTime := sTime.Add(time.Second * (86400 - 1))

	return sTime, eTime
}

func GetDateStart(rawTime time.Time) time.Time {
	sTime, err := time.ParseInLocation(DateLayout, rawTime.Format(DateLayout), time.Local)
	if err != nil {
		logs.Error(nil, "GetDateStart", zap.Any("rawTime", rawTime), zap.Error(err))
	}
	return sTime
}

func GetDateEnd(rawTime time.Time) time.Time {
	sTime, err := time.ParseInLocation(DateLayout, rawTime.Format(DateLayout), time.Local)
	if err != nil {
		logs.Error(nil, "GetDateEnd", zap.Any("rawTime", rawTime), zap.Error(err))
	}
	eTime := sTime.Add(time.Second * (86400 - 1))

	return eTime
}

// parse local time with standard
func ParseLocalTime(ctx context.Context, localTime time.Time) time.Time {
	sTime := time.Unix(LocalTimeZoneOffset(ctx, localTime), 0)
	return sTime
}

func LocalTimeZoneOffset(ctx context.Context, localTime time.Time) int64 {
	// ZoneOffset...
	zone, offset := time.Now().Zone()
	logs.Debugw(ctx, "[ZoneOffset]", map[string]interface{}{"zone": zone, "offset": offset})

	localTs := localTime.Unix()
	localTs += int64(offset)

	return localTs
}

func TimeSerials(ctx context.Context, startTs int64, endTs int64, timeArea int64) map[string]int64 {
	// TimeLayout
	var appTimeLayout = "2006-01-02 15:04:05"

	// ZoneOffset...
	_, offset := time.Now().Zone()
	//logs.Debugw(ctx, "[ZoneOffset]", map[string]interface{}{"zone": zone, "offset": offset})

	//AtTs
	startTs += int64(offset)
	atTs := startTs - startTs%timeArea
	atTs -= int64(offset)

	// TestSerials
	timeSerials := make(map[string]int64)
	for ; atTs <= endTs; atTs += timeArea {
		atTime := time.Unix(atTs, 0).Format(appTimeLayout)
		//logs.Debugw(ctx, "[AtData][atTs: %d][atTime: %s]", map[string]interface{}{"atTs": atTs, "atTime": atTime})
		timeSerials[atTime] = atTs
	}
	//logs.Debugw(ctx, "时间片段", map[string]interface{}{"startTs": startTs, "endTs": endTs, "atTs": atTs})

	return timeSerials
}

// GetBetweenDates 根据开始日期和结束日期计算出时间段内所有日期
// 参数为日期格式，如：2020-01-01
func GetBetweenDates(ctx context.Context, startDate, endDate string) ([]string, error) {
	if startDate == endDate {
		return []string{startDate}, nil
	}

	var d []string
	timeFormatTpl := "2006-01-02 15:04:05"
	if len(timeFormatTpl) != len(startDate) {
		timeFormatTpl = timeFormatTpl[0:len(startDate)]
	}
	date, err := time.Parse(timeFormatTpl, startDate)
	if err != nil {
		// 时间解析，异常
		logs.Error(ctx, "timeParse Error", zap.Any("rawTime", startDate), zap.Error(err))
		return d, err
	}
	date2, err := time.Parse(timeFormatTpl, endDate)
	if err != nil {
		// 时间解析，异常
		logs.Error(ctx, "timeParse Error", zap.Any("rawTime", endDate), zap.Error(err))
		return d, err
	}
	if date2.Before(date) {
		// 如果结束时间小于开始时间，异常
		logs.Error(ctx, "EndDate is smaller than StartDate", zap.Error(err))
		return d, err
	}
	// 输出日期格式固定
	timeFormatTpl = "2006-01-02"
	date2Str := date2.Format(timeFormatTpl)
	d = append(d, date.Format(timeFormatTpl))
	for {
		date = date.AddDate(0, 0, 1)
		dateStr := date.Format(timeFormatTpl)
		d = append(d, dateStr)
		if dateStr == date2Str {
			break
		}
	}
	return d, err
}
