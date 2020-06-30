package datetime

import (
	"database/sql/driver"
	"encoding/json"
	"math"
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

// UnmarshalJSON 转换 JSON
func (ins *DateTime) UnmarshalJSON(input []byte) error {
	var timestamp float64
	if err := json.Unmarshal(input, &timestamp); err == nil {
		sec := math.Floor(timestamp)
		nsec := int64(math.Floor((timestamp - sec) * 1000000000))
		ins.Display = SecondsISO8601
		ins.Time = time.Unix(int64(sec), nsec)
		return nil
	}

	var timeIns time.Time
	if err := json.Unmarshal(input, &timeIns); err != nil {
		return err
	}

	ins.Display = SecondsISO8601
	ins.Time = timeIns
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
	ins.Time = value.(time.Time)
	ins.Display = SecondsISO8601
	return nil
}

// Value database/sql Value 接口实现
func (ins DateTime) Value() (driver.Value, error) {
	return ins.Time, nil
}
