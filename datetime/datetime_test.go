package datetime_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/alphatr/go-utils/datetime"
)

type JsonData struct {
	CreateTime *datetime.DateTime `json:"create_time"`
	UpdateTime *datetime.DateTime `json:"update_time"`
	AddTime    *datetime.DateTime `json:"add_time"`
}

type OutputData struct {
	CreateTime int64  `json:"create_time"`
	UpdateTime int64  `json:"update_time"`
	AddTime    string `json:"add_time"`
}

func TestTime(test *testing.T) {
	now := time.Now()
	value := JsonData{
		CreateTime: datetime.New(now).SetDisplay(datetime.UnixTimestamp),
		AddTime:    datetime.New(now).SetFormat(time.RFC3339),
		UpdateTime: datetime.Now().SetDisplay(datetime.MillisecondsTimestamp),
	}

	result := OutputData{}
	data, err := json.Marshal(value)
	if err != nil {
		test.Fatal(err.Error())
	}

	if err := json.Unmarshal(data, &result); err != nil {
		test.Fatal(err.Error())
	}

	if result.AddTime != now.Format(time.RFC3339) {
		test.Fatalf("%s != %s", result.AddTime, now.Format(time.RFC3339))
	}

	if result.CreateTime != now.Unix() {
		test.Fatalf("%d != %d", result.CreateTime, now.Unix())
	}

	if result.UpdateTime <= 0 {
		test.Fatal("UpdateTime == 0")
	}
}
