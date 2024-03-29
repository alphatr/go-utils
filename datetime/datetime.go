package datetime

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"
)

// display 时间输出类型
type display uint

// DateTime 时间类型
type DateTime struct {
	time.Time
	Display display
	Format  string
}

var timezoneCache = map[string]*time.Location{}

// 默认支持的显示类型
const (
	SecondsISO8601 display = iota
	MillisecondsISO8601
	ISO8601
	UnixTimestamp
	MillisecondsTimestamp
	Custom
)

// New 创建时间
func New(value time.Time) *DateTime {
	return &DateTime{Time: value, Display: SecondsISO8601}
}

// Now 当前时间
func Now() *DateTime {
	return &DateTime{Time: time.Now(), Display: SecondsISO8601}
}

// Parse 解析时间字符串
func Parse(value string, formatters ...string) (*DateTime, error) {
	formatter := time.RFC3339
	if len(formatters) == 1 {
		formatter = formatters[0]
	}

	res, err := time.Parse(formatter, value)
	if err != nil {
		return nil, err
	}

	return New(res), nil
}

// Parse 解析时间字符串
func ParseInTimezone(value string, timezone string, formatters ...string) (*DateTime, error) {
	formatter := time.RFC3339
	if len(formatters) == 1 {
		formatter = formatters[0]
	}

	if timezone == "" {
		timezone = "Local"
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, err
	}

	res, err := time.ParseInLocation(formatter, value, loc)
	if err != nil {
		return nil, err
	}

	return New(res), nil
}

// SetDisplay 设置输出类型
func (ins *DateTime) SetDisplay(option display) *DateTime {
	ins.Display = option
	return ins
}

// SetFormat 设置输出格式
func (ins *DateTime) SetFormat(format string) *DateTime {
	ins.Display = Custom
	ins.Format = format
	return ins
}

// SetTimezone 设置时区
func (ins *DateTime) SetTimezone(name string) *DateTime {
	if zone, ok := timezoneCache[name]; ok {
		ins.Time = ins.Time.In(zone)
		return ins
	}

	zone, err := time.LoadLocation(name)
	if err != nil {
		return ins
	}

	timezoneCache[name] = zone
	ins.Time = ins.Time.In(zone)
	return ins
}

// UnmarshalJSON 转换 JSON
func (ins *DateTime) UnmarshalJSON(input []byte) error {
	var timestamp float64
	if err := json.Unmarshal(input, &timestamp); err == nil {
		sec := math.Floor(timestamp)
		nsec := int64(math.Floor((timestamp - sec) * 1000000000))
		ins.Display = SecondsISO8601
		ins.Time = time.Unix(int64(sec), nsec).Local()
		return nil
	}

	var timeIns time.Time
	if err := json.Unmarshal(input, &timeIns); err != nil {
		return err
	}

	ins.Display = SecondsISO8601
	ins.Time = timeIns.Local()
	return nil
}

// MarshalJSON 转
func (ins DateTime) MarshalJSON() ([]byte, error) {
	if ins.Display == UnixTimestamp {
		return json.Marshal(ins.Unix())
	}

	if ins.Display == MillisecondsTimestamp {
		return json.Marshal(ins.UnixNano() / 1000000)
	}

	formatMap := map[display]string{
		SecondsISO8601:      time.RFC3339,
		MillisecondsISO8601: "2006-01-02T15:04:05.999Z07:00",
		ISO8601:             time.RFC3339Nano,
		Custom:              ins.Format,
	}

	return json.Marshal(ins.Time.Format(formatMap[ins.Display]))
}

// Scan database/sql Scan 接口实现
func (ins *DateTime) Scan(value interface{}) error {
	timeValue := value.(time.Time)
	ins.Time = timeValue.Local()
	ins.Display = SecondsISO8601
	return nil
}

// Value database/sql Value 接口实现
func (ins DateTime) Value() (driver.Value, error) {
	return ins.Time, nil
}

func (ins *DateTime) UnmarshalText(input []byte) error {
	timestamp, err := strconv.ParseFloat(string(input), 64)
	if err == nil {
		sec := math.Floor(timestamp)
		nsec := int64(math.Floor((timestamp - sec) * 1000000000))
		ins.Display = SecondsISO8601
		ins.Time = time.Unix(int64(sec), nsec).Local()
		return nil
	}

	content, _ := json.Marshal(string(input))
	var timeIns time.Time
	if err := json.Unmarshal(content, &timeIns); err != nil {
		return err
	}

	ins.Display = SecondsISO8601
	ins.Time = timeIns.Local()
	return nil
}

func (ins DateTime) MarshalText() ([]byte, error) {
	if ins.Display == UnixTimestamp {
		return []byte(fmt.Sprintf("%d", ins.Unix())), nil
	}

	if ins.Display == MillisecondsTimestamp {
		return []byte(fmt.Sprintf("%d", ins.UnixNano()/1000000)), nil
	}

	formatMap := map[display]string{
		SecondsISO8601:      time.RFC3339,
		MillisecondsISO8601: "2006-01-02T15:04:05.999Z07:00",
		ISO8601:             time.RFC3339Nano,
		Custom:              ins.Format,
	}

	return []byte(ins.Time.Format(formatMap[ins.Display])), nil
}
