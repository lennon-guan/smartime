package smartime_test

import (
	"testing"
	"time"

	"github.com/lennon-guan/smartime"
)

func assertParseTimeEqualTo(t *testing.T, bt *smartime.BaseTime, s string, expected time.Time) {
	if v, err := bt.ParseTime(s); err != nil {
		t.Errorf("ParseTime(%#v) returns error: %+v", s, err)
	} else if !v.Equal(expected) {
		t.Errorf("ParseTime(%#v) returns %+v, not equals to %+v", s, v, expected)
	}
}

func assertParseTimeSatisfied(t *testing.T, bt *smartime.BaseTime, s string, check func(time.Time) bool) {
	if v, err := bt.ParseTime(s); err != nil {
		t.Errorf("ParseTime(%#v) returns error: %+v", s, err)
	} else if !check(v) {
		t.Errorf("ParseTime(%#v) returns %+v, not satisfied", s, v)
	}
}

func TestParseTimeAbsolute(t *testing.T) {
	var (
		now = time.Now()
		loc = now.Location()
		bt  = smartime.NewBaseTime(now)
	)
	assertParseTimeSatisfied(t, bt, "0", func(t time.Time) bool { return t.IsZero() })
	assertParseTimeSatisfied(t, bt, "zero", func(t time.Time) bool { return t.IsZero() })
	assertParseTimeEqualTo(t, bt, "241204", time.Date(2024, 12, 4, 0, 0, 0, 0, loc))
	assertParseTimeEqualTo(t, bt, "20241204", time.Date(2024, 12, 4, 0, 0, 0, 0, loc))
	assertParseTimeEqualTo(t, bt, "2024-12-04", time.Date(2024, 12, 4, 0, 0, 0, 0, loc))
	assertParseTimeEqualTo(t, bt, "1400010000", time.Unix(1400010000, 0))
	assertParseTimeEqualTo(t, bt, "1400010056123", time.Unix(1400010056, 123*1e6))
	assertParseTimeEqualTo(t, bt, "20241204154931", time.Date(2024, 12, 4, 15, 49, 31, 0, loc))
	assertParseTimeEqualTo(t, bt, "2024-12-04 15:49:31", time.Date(2024, 12, 4, 15, 49, 31, 0, loc))
	assertParseTimeEqualTo(t, bt, "2024-12-04 15:49:31-03", time.Date(2024, 12, 4, 18, 49, 31, 0, time.UTC))
	assertParseTimeEqualTo(t, bt, "2024-12-04 15:49:31+0300", time.Date(2024, 12, 4, 12, 49, 31, 0, time.UTC))
	assertParseTimeEqualTo(t, bt, "2024-12-04 15:49:31.123-03", time.Date(2024, 12, 4, 18, 49, 31, 123*1e6, time.UTC))
	assertParseTimeEqualTo(t, bt, "2024-12-04 15:49:31.123+0300", time.Date(2024, 12, 4, 12, 49, 31, 123*1e6, time.UTC))
}

func TestParseTimeRelative(t *testing.T) {
	var (
		loc = time.Local
		now = time.Date(2024, 12, 4, 11, 22, 33, 0, loc)
		bt  = smartime.NewBaseTime(now)
	)
	assertParseTimeEqualTo(t, bt, "+1h", now.Add(time.Hour))
	assertParseTimeEqualTo(t, bt, "-1h", now.Add(-time.Hour))
	assertParseTimeEqualTo(t, bt, "now", now)
	assertParseTimeEqualTo(t, bt, "today", time.Date(2024, 12, 4, 0, 0, 0, 0, loc))
	assertParseTimeEqualTo(t, bt, "thisMonth", time.Date(2024, 12, 1, 0, 0, 0, 0, loc))
	assertParseTimeEqualTo(t, bt, "lastMonth", time.Date(2024, 11, 1, 0, 0, 0, 0, loc))
	assertParseTimeEqualTo(t, bt, "nextMonth", time.Date(2025, 1, 1, 0, 0, 0, 0, loc))
	assertParseTimeEqualTo(t, bt, "now+1h", time.Date(2024, 12, 4, 12, 22, 33, 0, loc))
	assertParseTimeEqualTo(t, bt, "now-1m", time.Date(2024, 12, 4, 11, 21, 33, 0, loc))
	assertParseTimeEqualTo(t, bt, "today+2s", time.Date(2024, 12, 4, 0, 0, 2, 0, loc))
	assertParseTimeEqualTo(t, bt, "today+1m2s", time.Date(2024, 12, 4, 0, 1, 2, 0, loc))
	assertParseTimeEqualTo(t, bt, "today-2s", time.Date(2024, 12, 3, 23, 59, 58, 0, loc))
	assertParseTimeEqualTo(t, bt, "thisMonth+30ms", time.Date(2024, 12, 1, 0, 0, 0, 30*1e6, loc))
	assertParseTimeEqualTo(t, bt, "lastMonth+30us", time.Date(2024, 11, 1, 0, 0, 0, 30*1e3, loc))
	assertParseTimeEqualTo(t, bt, "lastMonth+30µs", time.Date(2024, 11, 1, 0, 0, 0, 30*1e3, loc))
	assertParseTimeEqualTo(t, bt, "nextMonth+30ns", time.Date(2025, 1, 1, 0, 0, 0, 30, loc))
}

func TestParseTimeCustomParser(t *testing.T) {
	var (
		now = time.Now()
		loc = now.Location()
		bt  = smartime.NewBaseTime(now).WithCustomParser(func(s string, loc *time.Location) (t time.Time, err error) {
			if t, err = time.ParseInLocation("2006/01/02", s, loc); err == nil {
			} else if t, err = time.ParseInLocation("2006/1/2", s, loc); err == nil {
			} else if t, err = time.ParseInLocation("2006.1.2", s, loc); err == nil {
			} else if t, err = time.ParseInLocation("2006.01.02", s, loc); err == nil {
			}
			return
		})
	)
	assertParseTimeEqualTo(t, bt, "20241204", time.Date(2024, 12, 4, 0, 0, 0, 0, loc))
	assertParseTimeEqualTo(t, bt, "2024/12/4", time.Date(2024, 12, 4, 0, 0, 0, 0, loc))
	assertParseTimeEqualTo(t, bt, "2024/12/04", time.Date(2024, 12, 4, 0, 0, 0, 0, loc))
	assertParseTimeEqualTo(t, bt, "2024/2/4", time.Date(2024, 2, 4, 0, 0, 0, 0, loc))
	assertParseTimeEqualTo(t, bt, "2024.12.4", time.Date(2024, 12, 4, 0, 0, 0, 0, loc))
	assertParseTimeEqualTo(t, bt, "2024.12.04", time.Date(2024, 12, 4, 0, 0, 0, 0, loc))
	assertParseTimeEqualTo(t, bt, "2024.2.4", time.Date(2024, 2, 4, 0, 0, 0, 0, loc))
}
