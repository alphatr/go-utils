package duration

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// Duration 时长对象
type Duration struct {
	Year   int
	Month  int
	Week   int
	Day    int
	Hour   int
	Minute int
	Second int
}

var pattern = regexp.MustCompile(`^P((?P<year>\d+)Y)?((?P<month>\d+)M)?((?P<week>\d+)W)?((?P<day>\d+)D)?(T((?P<hour>\d+)H)?((?P<minute>\d+)M)?((?P<second>\d+)S)?)?$`)

// Parse 解析 Duration 字符串
func Parse(input string) (Duration, error) {
	match := []string{}
	result := Duration{}

	if !pattern.MatchString(input) {
		return result, errors.New("could not parse duration string")
	}

	match = pattern.FindStringSubmatch(input)
	for i, name := range pattern.SubexpNames() {
		part := match[i]
		if i == 0 || name == "" || part == "" {
			continue
		}

		value, err := strconv.ParseInt(part, 10, 32)
		if err != nil {
			return result, err
		}

		switch name {
		case "year":
			result.Year = int(value)
		case "month":
			result.Month = int(value)
		case "week":
			result.Week = int(value)
		case "day":
			result.Day = int(value)
		case "hour":
			result.Hour = int(value)
		case "minute":
			result.Minute = int(value)
		case "second":
			result.Second = int(value)
		default:
			return result, fmt.Errorf("unknown field %s", name)
		}
	}

	return result, nil
}

// IsDateZero 是否日期零值
func (ins *Duration) IsDateZero() bool {
	return ins.Year == 0 && ins.Month == 0 && ins.Week == 0 && ins.Day == 0
}

// IsTimeZero 是否时间零值
func (ins *Duration) IsTimeZero() bool {
	return ins.Hour == 0 && ins.Minute == 0 && ins.Second == 0
}

// IsZero 是否零值
func (ins *Duration) IsZero() bool {
	return ins.IsDateZero() && ins.IsTimeZero()
}

// Shift 时间间隔计算
func (ins *Duration) Shift(input time.Time) time.Time {
	if !ins.IsDateZero() {
		days := ins.Week*7 + ins.Day
		input = input.AddDate(ins.Year, ins.Month, days)
	}

	if !ins.IsTimeZero() {
		input = input.Add(ins.timeDuration())
	}

	return input
}

func (ins *Duration) timeDuration() time.Duration {
	var (
		hour   = time.Duration(ins.Hour) * time.Hour
		minute = time.Duration(ins.Minute) * time.Minute
		second = time.Duration(ins.Second) * time.Second
	)

	return hour + minute + second
}
