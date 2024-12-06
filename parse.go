package smartime

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

type BaseTime struct {
	t time.Time
	c func(string, *time.Location) (time.Time, error)
}

func NewBaseTime(t time.Time) *BaseTime {
	return &BaseTime{t: t}
}
func NowBase() *BaseTime {
	return NewBaseTime(time.Now())
}

func (bt *BaseTime) Time() time.Time {
	return bt.t
}

func (bt *BaseTime) WithCustomParser(f func(string, *time.Location) (time.Time, error)) *BaseTime {
	bt.c = f
	return bt
}

var (
	relativeRegex = regexp.MustCompile(`^(now|today|thisMonth|nextMonth|lastMonth)?([\+|\-](\d+(ns|us|µs|ms|s|m|h))+)?$`)
	zeroTime      = time.Time{}
)

func (bt *BaseTime) ParseTime(s string) (t time.Time, err error) {
	var (
		ts  int64
		loc = bt.Time().Location()
	)
	if bt.c != nil {
		if t, e := bt.c(s, loc); e == nil {
			return t, nil
		}
	}
	if s == "" {
		err = errors.New("time string cannot be empty")
		return
	} else if s == "zero" || s == "0" {
		t = zeroTime
	} else if s == "now" {
		t = bt.Time()
	} else if f := s[0]; f == '-' || f == '+' || f == 'n' || f == 't' || f == 'l' {
		if m := relativeRegex.FindStringSubmatch(s); len(m) > 2 {
			tt := bt.Time()
			switch m[1] {
			case "":
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
				}
				if offset[0] == '+' {
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
		case 23: // YYYY-mm-dd HH:MM:SS.XXX
			t, err = time.Parse("2006-01-02 15:04:05.000", s)
		case 24: // YYYY-mm-dd HH:MM:SS+XXXX
			t, err = time.Parse("2006-01-02 15:04:05-0700", s)
		case 26: // YYYY-mm-dd HH:MM:SS.XXX+XX
			t, err = time.Parse("2006-01-02 15:04:05.000-07", s)
		case 28: // YYYY-mm-dd HH:MM:SS.XXX+XXXX
			t, err = time.Parse("2006-01-02 15:04:05.000-0700", s)
		default:
			// complex relative time format
			err = fmt.Errorf("unsupported time format: %s", s)
		}
	}
	return
}

func (bt *BaseTime) MustParseTime(s string) time.Time {
	if t, err := bt.ParseTime(s); err != nil {
		panic(err)
	} else {
		return t
	}
}

func ParseTime(s string) (time.Time, error) {
	return NowBase().ParseTime(s)
}

func MustParseTime(s string) time.Time {
	return NowBase().MustParseTime(s)
}
