package smartime_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/lennon-guan/smartime"
)

func assertParseTime(t *testing.T, bt smartime.BaseTime, s string, expected func(time.Time, error) string) {
	if errMsg := expected(bt.ParseTime(s)); errMsg != "" {
		t.Error(errMsg)
	}
}

func assertParseTimeEqualTo(t *testing.T, bt smartime.BaseTime, s string, expected time.Time) {
	assertParseTime(t, bt, s, func(v time.Time, err error) string {
		if err != nil {
			return fmt.Sprintf("ParseTime(%#v) returns error: %+v", s, err)
		} else if !v.Equal(expected) {
			return fmt.Sprintf("ParseTime(%#v) returns %+v, not equals to %+v", s, v, expected)
		}
		return ""
	})
}

func TestParseTimeAbsolute(t *testing.T) {
	var (
		now = time.Now()
		loc = now.Location()
		bt  = smartime.BaseTime(now)
	)
	assertParseTimeEqualTo(t, bt, "241204", time.Date(2024, 12, 4, 0, 0, 0, 0, loc))
	assertParseTimeEqualTo(t, bt, "20241204", time.Date(2024, 12, 4, 0, 0, 0, 0, loc))
	assertParseTimeEqualTo(t, bt, "2024-12-04", time.Date(2024, 12, 4, 0, 0, 0, 0, loc))
	assertParseTimeEqualTo(t, bt, "1400010000", time.Unix(1400010000, 0))
	assertParseTimeEqualTo(t, bt, "1400010056123", time.Unix(1400010056, 123000000))
	assertParseTimeEqualTo(t, bt, "20241204154931", time.Date(2024, 12, 4, 15, 49, 31, 0, loc))
	assertParseTimeEqualTo(t, bt, "2024-12-04 15:49:31", time.Date(2024, 12, 4, 15, 49, 31, 0, loc))
	assertParseTimeEqualTo(t, bt, "2024-12-04 15:49:31-03", time.Date(2024, 12, 4, 18, 49, 31, 0, time.UTC))
	assertParseTimeEqualTo(t, bt, "2024-12-04 15:49:31+0300", time.Date(2024, 12, 4, 12, 49, 31, 0, time.UTC))
}

func TestParseTimeRelative(t *testing.T) {
	var (
		loc = time.Local
		now = time.Date(2024, 12, 4, 11, 22, 33, 0, loc)
		bt  = smartime.BaseTime(now)
	)
	assertParseTimeEqualTo(t, bt, "+1h", now.Add(time.Hour))
	assertParseTimeEqualTo(t, bt, "-1h", now.Add(-time.Hour))
	assertParseTimeEqualTo(t, bt, "-1h", now.Add(-time.Hour))
	assertParseTimeEqualTo(t, bt, "now", now)
	assertParseTimeEqualTo(t, bt, "today", time.Date(2024, 12, 4, 0, 0, 0, 0, loc))
	assertParseTimeEqualTo(t, bt, "thisMonth", time.Date(2024, 12, 1, 0, 0, 0, 0, loc))
	assertParseTimeEqualTo(t, bt, "lastMonth", time.Date(2024, 11, 1, 0, 0, 0, 0, loc))
	assertParseTimeEqualTo(t, bt, "nextMonth", time.Date(2025, 1, 1, 0, 0, 0, 0, loc))
	assertParseTimeEqualTo(t, bt, "now+1h", time.Date(2024, 12, 4, 12, 22, 33, 0, loc))
	assertParseTimeEqualTo(t, bt, "now-1m", time.Date(2024, 12, 4, 11, 21, 33, 0, loc))
	assertParseTimeEqualTo(t, bt, "today+2s", time.Date(2024, 12, 4, 0, 0, 2, 0, loc))
	assertParseTimeEqualTo(t, bt, "today-2s", time.Date(2024, 12, 3, 23, 59, 58, 0, loc))
	assertParseTimeEqualTo(t, bt, "thisMonth", time.Date(2024, 12, 1, 0, 0, 0, 0, loc))
	assertParseTimeEqualTo(t, bt, "lastMonth", time.Date(2024, 11, 1, 0, 0, 0, 0, loc))
	assertParseTimeEqualTo(t, bt, "nextMonth", time.Date(2025, 1, 1, 0, 0, 0, 0, loc))
}
