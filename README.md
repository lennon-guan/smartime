# smartime Package Documentation

The `smartime` package is designed to provide flexible time parsing functionality. It introduces the `BaseTime` type, which extends Go's native `time.Time` type, and offers utilities for parsing time strings with both relative and absolute formats.

## Package smartime

```go
import (
    "errors"
    "fmt"
    "regexp"
    "strconv"
    "strings"
    "time"
)
```

### Types

#### `type BaseTime`

`BaseTime` is a custom type that extends the `time.Time` type. 

### Functions

#### `func NowBase() BaseTime`

`NowBase` returns the current time as a `BaseTime` type. It serves as a starting point for parsing relative time strings.

```go
func NowBase() BaseTime {
    return BaseTime(time.Now())
}
```

#### `func (bt BaseTime) ParseTime(s string) (time.Time, error)`

Parses a time string `s` based on the `BaseTime` object `bt` and returns a `time.Time` object or an error. This method supports parsing both relative and absolute time formats.

- **Relative time formats** include:
  - `+{duration}`: Duration after the current time, e.g., "+2h" for 2 hours in the future.
  - `-{duration}`: Duration before the current time, e.g., "-30m" for 30 minutes in the past.
  - Keywords such as `now`, `today`, `thisMonth`, `nextMonth`, `lastMonth` with optional offsets (e.g., `thisMonth+2d` for two days after the start of this month).

- **Absolute time formats** include:
  - "yymmdd"
  - "YYYYmmdd", "yy-mm-dd"
  - "YYYY-mm-dd"
  - "timestamp" (numeric, up to millisecond precision)
  - "YYYYmmddHHMMSS"
  - "YYYY-mm-dd HH:MM:SS"
  - "YYYY-mm-dd HH:MM:SS"
  - Timezone notations ("YYYY-mm-dd HH:MM:SS+XX", "YYYY-mm-dd HH:MM:SS+XXXX")

You can find more formats and examples in `parse_test.go`

```go
func (bt BaseTime) ParseTime(s string) (time.Time, error)
```

#### `func (bt BaseTime) MustParseTime(s string) time.Time`

Similar to `ParseTime`, but panics if the parsing fails. Useful when the format is expected to be valid and you prefer to handle parsing errors with exceptions.

```go
func (bt BaseTime) MustParseTime(s string) time.Time
```

#### `func ParseTime(s string) (time.Time, error)`

Utility function that parses the given time string `s` using the current time as the base. Equivalent to calling `NowBase().ParseTime(s)`.

```go
func ParseTime(s string) (time.Time, error)
```

#### `func MustParseTime(s string) time.Time`

Utility function that parses the given time string `s`, panicking if the parsing fails. Equivalent to `NowBase().MustParseTime(s)`.

```go
func MustParseTime(s string) time.Time
```

## Usage Examples

```go
package main

import (
    "fmt"
    "smartime"
)

func main() {
    // Parse a simple absolute date
    t, err := smartime.ParseTime("20220101")
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Parsed time:", t)
    }

    // Parse a relative date
    t, err = smartime.ParseTime("-2h")
    if err != nil {
        fmt.Println("Error:", err)
    } else {
        fmt.Println("Parsed time, 2 hours ago:", t)
    }

    // Parse using MustParseTime to allow program to panic on error
    fmt.Println("Parsed time at the start of this month:", smartime.MustParseTime("thisMonth"))
}
```

This package is particularly suitable for applications needing dynamic date calculations, flexible date parsing, or domain-specific date handling logic.
