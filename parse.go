package smartime

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type BaseTime time.Time

func NowBase() BaseTime {
	return BaseTime(time.Now())
}

var relativeRegex = regexp.MustCompile(`^(now|today|thisMonth|nextMonth|lastMonth)([\+|\-]\d+(ns|us|µs|ms|s|m|h))?$`)

func (bt BaseTime) ParseTime(s string) (t time.Time, err error) {
	var (
		ts  int64
		loc = time.Time(bt).Location()
	)
	if s == "" {
		err = errors.New("time string cannot be empty")
		return
	} else if strings.HasPrefix(s, "+") { // Relative time: duration after now
		if du, err := time.ParseDuration(s[1:]); err != nil {
			return t, err
		} else {
			t = time.Time(bt).Add(du)
		}
	} else if strings.HasPrefix(s, "-") { // Relative time: duration before now
		if du, err := time.ParseDuration(s[1:]); err != nil {
			return t, err
		} else {
			t = time.Time(bt).Add(-du)
		}
	} else if s == "now" {
		t = time.Time(bt)
	} else if f := s[0]; f == 'n' || f == 't' || f == 'l' {
		if m := relativeRegex.FindStringSubmatch(s); len(m) > 2 {
			tt := time.Time(bt)
			switch m[1] {
			case "now":
			case "today":
				tt = time.Date(tt.Year(), tt.Month(), tt.Day(), 0, 0, 0, 0, tt.Location())
			case "thisMonth":
				tt = time.Date(tt.Year(), tt.Month(), 1, 0, 0, 0, 0, tt.Location())
			case "lastMonth":
				y := tt.Year()
				m := int(tt.Month() - 1)
				if m < 1 {
					m = 12
					y--
				}
				tt = time.Date(y, time.Month(m), 1, 0, 0, 0, 0, tt.Location())
			case "nextMonth":
				y := tt.Year()
				m := int(tt.Month() + 1)
				if m > 12 {
					m = 1
					y++
				}
				tt = time.Date(y, time.Month(m), 1, 0, 0, 0, 0, tt.Location())
			default:
				err = fmt.Errorf("unsupported time format: %s", s)
				return
			}
			if offset := m[2]; offset != "" {
				var du time.Duration
				if du, err = time.ParseDuration(offset[1:]); err != nil {
					return
				} else if offset[0] == '+' {
					tt = tt.Add(du)
				} else {
					tt = tt.Add(-du)
				}
			}
			t = tt
		} else {
			err = fmt.Errorf("unsupported time format: %s", s)
		}
	} else { // Absolute time
		switch len(s) {
		case 6: // yymmdd
			t, err = time.ParseInLocation("060102", s, loc)
		case 8: // YYYYmmdd yy-mm-dd
			if t, err = time.ParseInLocation("20060102", s, loc); err == nil {
			} else if t, err = time.ParseInLocation("06-01-02", s, loc); err == nil {
			}
		case 10: // YYYY-mm-dd timestamp(to second)
			if t, err = time.ParseInLocation("2006-01-02", s, loc); err == nil {
			} else if ts, err = strconv.ParseInt(s, 10, 64); err == nil {
				t = time.Unix(ts, 0)
			}
		case 13: // timestamp(to millisecond)
			if ts, err = strconv.ParseInt(s, 10, 64); err == nil {
				t = time.Unix(ts/1000, (ts%1000)*1e6)
			}
		case 14: // YYYYmmddHHMMSS
			t, err = time.ParseInLocation("20060102150405", s, loc)
		case 19: // YYYY-mm-dd HH:MM:SS
			t, err = time.ParseInLocation("2006-01-02 15:04:05", s, loc)
		case 22: // YYYY-mm-dd HH:MM:SS+XX
			t, err = time.Parse("2006-01-02 15:04:05-07", s)
		case 24: // YYYY-mm-dd HH:MM:SS+XXXX
			t, err = time.Parse("2006-01-02 15:04:05-0700", s)
		default:
			// complex relative time format
			err = fmt.Errorf("unsupported time format: %s", s)
		}
	}
	return
}

func ParseTime(s string) (time.Time, error) {
	return NowBase().ParseTime(s)
}
