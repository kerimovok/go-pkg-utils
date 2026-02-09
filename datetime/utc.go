package datetime

import "time"

// ParseRFC3339ToUTC parses an RFC3339 timestamp string and returns a time in UTC.
// If parsing fails, it returns a zero time and the parse error.
func ParseRFC3339ToUTC(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}, err
	}
	return t.UTC(), nil
}

// ParseRFC3339PtrToUTC parses an RFC3339 timestamp pointer and returns a *time.Time in UTC.
// Nil or empty input returns nil.
func ParseRFC3339PtrToUTC(s *string) (*time.Time, error) {
	if s == nil || *s == "" {
		return nil, nil
	}
	t, err := ParseRFC3339ToUTC(*s)
	if err != nil {
		return nil, err
	}
	if t.IsZero() {
		return nil, nil
	}
	return &t, nil
}

// ToUTC converts a time to UTC (no-op for zero time).
func ToUTC(t time.Time) time.Time {
	if t.IsZero() {
		return t
	}
	return t.UTC()
}

// ToUTCOrZero converts a time to UTC, but returns zero if the input is zero.
func ToUTCOrZero(t time.Time) time.Time {
	if t.IsZero() {
		return time.Time{}
	}
	return t.UTC()
}

// FormatRFC3339UTC formats a time in RFC3339 with UTC zone.
func FormatRFC3339UTC(t time.Time) string {
	return ToUTC(t).Format(time.RFC3339)
}

// ExcelDateTimeLayout is the shared Excel-friendly date-time layout.
const ExcelDateTimeLayout = "2006-01-02 15:04:05"

// FormatExcelUTC formats a time in the shared Excel layout, normalized to UTC.
// Zero times format to an empty string.
func FormatExcelUTC(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return ToUTC(t).Format(ExcelDateTimeLayout)
}
