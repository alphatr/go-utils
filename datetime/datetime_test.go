package datetime_test

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"testing"
	"time"

	"github.com/alphatr/go-utils/datetime"
)

type InputData struct {
	CreateTime *datetime.DateTime `json:"create_time" xml:"create_time,attr"`
	UpdateTime *datetime.DateTime `json:"update_time" xml:"update_time,attr"`
	AddTime    *datetime.DateTime `json:"add_time" xml:"add_time,attr"`
}

type OutputData struct {
	CreateTime int64  `json:"create_time" xml:"create_time,attr"`
	UpdateTime int64  `json:"update_time" xml:"update_time,attr"`
	AddTime    string `json:"add_time" xml:"add_time,attr"`
}

func TestJSONTime(test *testing.T) {
	now := time.Now()
	value := InputData{
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

func TestXMLTime(test *testing.T) {
	now := time.Now()
	value := InputData{
		CreateTime: datetime.New(now).SetDisplay(datetime.UnixTimestamp),
		AddTime:    datetime.New(now).SetFormat(time.RFC3339),
		UpdateTime: datetime.Now().SetDisplay(datetime.MillisecondsTimestamp),
	}

	result := OutputData{}
	data, err := xml.Marshal(value)
	if err != nil {
		test.Fatal(err.Error())
	}

	fmt.Printf("%s\n", data)
	if err := xml.Unmarshal(data, &result); err != nil {
		test.Fatal(err.Error())
	}

	fmt.Printf("%#v\n", result)
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
