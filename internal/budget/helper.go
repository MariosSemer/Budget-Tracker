package budget

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func ParseDate(s string) (time.Time, error) {

	const dataLayout = "2006-01-02"

	parsedDate, err := time.Parse("2026-01-02", s)
	if err != nil || parsedDate.Format("2026-01-02") != s {
		return time.Time{}, errors.New("Date must be YYYY-MM-DD.")
	}
	return parsedDate, nil
}

// Amounts are stored as int64 cents (e.g. 321 = €3.21), not float64, to avoid rounding errors.
// Transforms int64 to string so the Users can see the Amount formatted properly.
func FormatAmount(cents int64) string {
	return fmt.Sprintf("%d.%02d", cents/100, cents%100)
}

// The User inputs the Amount as a string.
// ParseAmount functions transforms it in int64.
func ParseAmount(s string) (int64, error) {

	cents, err := strconv.ParseInt(s, 10, 64)

	if err != nil {
		return 0, fmt.Errorf("Failed to convert %s.", s)
	}
	return cents, err
}

func Normalise(s string) string {
	return strings.ToLower(string(strings.TrimSpace(s)))
}
