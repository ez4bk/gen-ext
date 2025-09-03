package ezgen

import (
	"strconv"
	"time"
)

type TimeRange struct {
	Start time.Time
	End   time.Time
}

func StartEndStr2TimeRange(startStr, endStr string) (timeRange TimeRange) {
	if startStr == "" && endStr == "" {
		return TimeRange{}
	}
	startMs, _ := strconv.ParseInt(startStr, 10, 64)
	start := time.Time{}
	if startMs > 0 {
		start = time.UnixMilli(startMs)
	}
	endMs, _ := strconv.ParseInt(endStr, 10, 64)
	end := time.Now()
	if endMs > 0 {
		end = time.UnixMilli(endMs)
	}
	timeRange = TimeRange{
		Start: start,
		End:   end,
	}
	return timeRange
}
