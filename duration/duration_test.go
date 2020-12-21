package duration_test

import (
	"testing"
	"time"

	"github.com/alphatr/go-utils/duration"
)

type shiftCase struct {
	from     string
	duration duration.Duration
	want     string
}

type parseCase struct {
	from string
	want duration.Duration
}

const timeFormat = "2006-01-02 03:04:05"

func buildTime(test *testing.T, input string) time.Time {
	result, err := time.Parse(timeFormat, input)
	if err != nil {
		test.Fatal(err)
	}

	return result
}

func TestCanParse(test *testing.T) {
	cases := []parseCase{
		{"P1Y", duration.Duration{Year: 1}},
		{"P1M", duration.Duration{Month: 1}},
		{"P1W", duration.Duration{Week: 1}},
		{"P1D", duration.Duration{Day: 1}},
		{"PT1H", duration.Duration{Hour: 1}},
		{"PT1M", duration.Duration{Minute: 1}},
		{"PT1S", duration.Duration{Second: 1}},
		{"P6Y1M7W2DT3H4M5S", duration.Duration{Year: 6, Month: 1, Week: 7, Day: 2, Hour: 3, Minute: 4, Second: 5}},
	}

	for index, item := range cases {
		got, err := duration.Parse(item.from)
		if err != nil {
			test.Fatal(err)
		}
		if item.want != got {
			test.Fatalf("Case %d: want=%+v, got=%+v", index, item.want, got)
		}
	}
}

func TestCanShift(test *testing.T) {
	cases := []shiftCase{
		{"2020-01-01 00:00:00", duration.Duration{}, "2020-01-01 00:00:00"},
		{"2020-01-01 00:00:00", duration.Duration{Year: 1}, "2021-01-01 00:00:00"},
		{"2020-01-01 00:00:00", duration.Duration{Month: 1}, "2020-02-01 00:00:00"},
		{"2020-01-01 00:00:00", duration.Duration{Month: 2}, "2020-03-01 00:00:00"},
		{"2020-01-01 00:00:00", duration.Duration{Week: 1}, "2020-01-08 00:00:00"},
		{"2020-01-01 00:00:00", duration.Duration{Day: 1}, "2020-01-02 00:00:00"},
		{"2020-01-01 00:00:00", duration.Duration{Hour: 1}, "2020-01-01 01:00:00"},
		{"2020-01-01 00:00:00", duration.Duration{Minute: 1}, "2020-01-01 00:01:00"},
		{"2020-01-01 00:00:00", duration.Duration{Second: 1}, "2020-01-01 00:00:01"},
		{
			"2020-01-01 00:00:00",
			duration.Duration{Year: 6, Month: 1, Day: 2, Week: 7, Hour: 3, Minute: 4, Second: 5},
			"2026-03-24 03:04:05",
		},
	}

	for index, item := range cases {
		from := buildTime(test, item.from)
		want := buildTime(test, item.want)

		got := item.duration.Shift(from)
		if !want.Equal(got) {
			test.Fatalf("Case %d: want=%s, got=%s", index, want, got)
		}
	}
}
