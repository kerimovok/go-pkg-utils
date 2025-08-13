package datetime

import (
	"fmt"
	"time"
)

// Common time constants
const (
	DayInSeconds   = 24 * 60 * 60
	WeekInSeconds  = 7 * DayInSeconds
	MonthInSeconds = 30 * DayInSeconds
	YearInSeconds  = 365 * DayInSeconds

	DayInMinutes   = 24 * 60
	WeekInMinutes  = 7 * DayInMinutes
	MonthInMinutes = 30 * DayInMinutes
	YearInMinutes  = 365 * DayInMinutes
)

// Common time formats
const (
	DateFormat          = "2006-01-02"
	TimeFormat          = "15:04:05"
	DateTimeFormat      = "2006-01-02 15:04:05"
	ISO8601Format       = "2006-01-02T15:04:05Z07:00"
	RFC3339Format       = time.RFC3339
	TimestampFormat     = "20060102150405"
	HumanDateFormat     = "January 2, 2006"
	HumanTimeFormat     = "3:04 PM"
	HumanDateTimeFormat = "January 2, 2006 at 3:04 PM"
)

// Now returns the current time
func Now() time.Time {
	return time.Now()
}

// NowUTC returns the current time in UTC
func NowUTC() time.Time {
	return time.Now().UTC()
}

// Today returns today's date at midnight
func Today() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

// TodayUTC returns today's date at midnight in UTC
func TodayUTC() time.Time {
	now := time.Now().UTC()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
}

// Yesterday returns yesterday's date at midnight
func Yesterday() time.Time {
	return Today().AddDate(0, 0, -1)
}

// Tomorrow returns tomorrow's date at midnight
func Tomorrow() time.Time {
	return Today().AddDate(0, 0, 1)
}

// StartOfWeek returns the start of the week (Monday) for the given time
func StartOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday = 7
	}
	return t.AddDate(0, 0, -weekday+1).Truncate(24 * time.Hour)
}

// EndOfWeek returns the end of the week (Sunday) for the given time
func EndOfWeek(t time.Time) time.Time {
	return StartOfWeek(t).AddDate(0, 0, 6).Add(24*time.Hour - time.Nanosecond)
}

// StartOfMonth returns the start of the month for the given time
func StartOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// EndOfMonth returns the end of the month for the given time
func EndOfMonth(t time.Time) time.Time {
	return StartOfMonth(t).AddDate(0, 1, 0).Add(-time.Nanosecond)
}

// StartOfYear returns the start of the year for the given time
func StartOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
}

// EndOfYear returns the end of the year for the given time
func EndOfYear(t time.Time) time.Time {
	return StartOfYear(t).AddDate(1, 0, 0).Add(-time.Nanosecond)
}

// IsToday checks if the given time is today
func IsToday(t time.Time) bool {
	now := time.Now()
	return t.Year() == now.Year() && t.Month() == now.Month() && t.Day() == now.Day()
}

// IsYesterday checks if the given time is yesterday
func IsYesterday(t time.Time) bool {
	yesterday := time.Now().AddDate(0, 0, -1)
	return t.Year() == yesterday.Year() && t.Month() == yesterday.Month() && t.Day() == yesterday.Day()
}

// IsTomorrow checks if the given time is tomorrow
func IsTomorrow(t time.Time) bool {
	tomorrow := time.Now().AddDate(0, 0, 1)
	return t.Year() == tomorrow.Year() && t.Month() == tomorrow.Month() && t.Day() == tomorrow.Day()
}

// IsWeekend checks if the given time is a weekend (Saturday or Sunday)
func IsWeekend(t time.Time) bool {
	weekday := t.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

// IsWeekday checks if the given time is a weekday (Monday to Friday)
func IsWeekday(t time.Time) bool {
	return !IsWeekend(t)
}

// IsPast checks if the given time is in the past
func IsPast(t time.Time) bool {
	return t.Before(time.Now())
}

// IsFuture checks if the given time is in the future
func IsFuture(t time.Time) bool {
	return t.After(time.Now())
}

// Age calculates the age in years from the given birth date
func Age(birthDate time.Time) int {
	now := time.Now()
	age := now.Year() - birthDate.Year()

	// Adjust if birthday hasn't occurred this year yet
	if now.Month() < birthDate.Month() ||
		(now.Month() == birthDate.Month() && now.Day() < birthDate.Day()) {
		age--
	}

	return age
}

// DaysBetween calculates the number of days between two dates
func DaysBetween(start, end time.Time) int {
	duration := end.Sub(start)
	return int(duration.Hours() / 24)
}

// BusinessDaysBetween calculates the number of business days between two dates
func BusinessDaysBetween(start, end time.Time) int {
	if start.After(end) {
		start, end = end, start
	}

	days := 0
	current := start.Truncate(24 * time.Hour)
	endDate := end.Truncate(24 * time.Hour)

	for current.Before(endDate) || current.Equal(endDate) {
		if IsWeekday(current) {
			days++
		}
		current = current.AddDate(0, 0, 1)
	}

	return days
}

// ParseDate parses a date string using common formats
func ParseDate(dateStr string) (time.Time, error) {
	formats := []string{
		DateFormat,
		DateTimeFormat,
		ISO8601Format,
		RFC3339Format,
		TimestampFormat,
		time.RFC3339,
		time.RFC822,
		time.RFC1123,
		"2006/01/02",
		"02/01/2006",
		"01/02/2006",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04",
		"02-01-2006",
		"01-02-2006",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date string: %s", dateStr)
}

// FormatDuration formats a duration in human-readable format
func FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.0fs", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.0fm", d.Minutes())
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%.1fh", d.Hours())
	}
	return fmt.Sprintf("%.1fd", d.Hours()/24)
}

// TimeAgo returns a human-readable string representing how long ago the time was
func TimeAgo(t time.Time) string {
	now := time.Now()
	duration := now.Sub(t)

	if duration < time.Minute {
		return "just now"
	}
	if duration < time.Hour {
		minutes := int(duration.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	}
	if duration < 24*time.Hour {
		hours := int(duration.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	}
	if duration < 7*24*time.Hour {
		days := int(duration.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	}
	if duration < 30*24*time.Hour {
		weeks := int(duration.Hours() / (7 * 24))
		if weeks == 1 {
			return "1 week ago"
		}
		return fmt.Sprintf("%d weeks ago", weeks)
	}
	if duration < 365*24*time.Hour {
		months := int(duration.Hours() / (30 * 24))
		if months == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	}

	years := int(duration.Hours() / (365 * 24))
	if years == 1 {
		return "1 year ago"
	}
	return fmt.Sprintf("%d years ago", years)
}

// TimeUntil returns a human-readable string representing how long until the time
func TimeUntil(t time.Time) string {
	now := time.Now()
	if t.Before(now) {
		return "in the past"
	}

	duration := t.Sub(now)

	if duration < time.Minute {
		return "in less than a minute"
	}
	if duration < time.Hour {
		minutes := int(duration.Minutes())
		if minutes == 1 {
			return "in 1 minute"
		}
		return fmt.Sprintf("in %d minutes", minutes)
	}
	if duration < 24*time.Hour {
		hours := int(duration.Hours())
		if hours == 1 {
			return "in 1 hour"
		}
		return fmt.Sprintf("in %d hours", hours)
	}
	if duration < 7*24*time.Hour {
		days := int(duration.Hours() / 24)
		if days == 1 {
			return "in 1 day"
		}
		return fmt.Sprintf("in %d days", days)
	}
	if duration < 30*24*time.Hour {
		weeks := int(duration.Hours() / (7 * 24))
		if weeks == 1 {
			return "in 1 week"
		}
		return fmt.Sprintf("in %d weeks", weeks)
	}
	if duration < 365*24*time.Hour {
		months := int(duration.Hours() / (30 * 24))
		if months == 1 {
			return "in 1 month"
		}
		return fmt.Sprintf("in %d months", months)
	}

	years := int(duration.Hours() / (365 * 24))
	if years == 1 {
		return "in 1 year"
	}
	return fmt.Sprintf("in %d years", years)
}

// AddBusinessDays adds business days to a date (skipping weekends)
func AddBusinessDays(t time.Time, days int) time.Time {
	result := t
	remaining := days

	if days > 0 {
		for remaining > 0 {
			result = result.AddDate(0, 0, 1)
			if IsWeekday(result) {
				remaining--
			}
		}
	} else {
		for remaining < 0 {
			result = result.AddDate(0, 0, -1)
			if IsWeekday(result) {
				remaining++
			}
		}
	}

	return result
}

// GetQuarter returns the quarter (1-4) for the given time
func GetQuarter(t time.Time) int {
	month := int(t.Month())
	return (month-1)/3 + 1
}

// StartOfQuarter returns the start of the quarter for the given time
func StartOfQuarter(t time.Time) time.Time {
	quarter := GetQuarter(t)
	month := (quarter-1)*3 + 1
	return time.Date(t.Year(), time.Month(month), 1, 0, 0, 0, 0, t.Location())
}

// EndOfQuarter returns the end of the quarter for the given time
func EndOfQuarter(t time.Time) time.Time {
	return StartOfQuarter(t).AddDate(0, 3, 0).Add(-time.Nanosecond)
}

// IsLeapYear checks if the given year is a leap year
func IsLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// DaysInMonth returns the number of days in the given month and year
func DaysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// ToUnixTimestamp converts time to Unix timestamp
func ToUnixTimestamp(t time.Time) int64 {
	return t.Unix()
}

// FromUnixTimestamp converts Unix timestamp to time
func FromUnixTimestamp(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}

// ToUnixMillis converts time to Unix timestamp in milliseconds
func ToUnixMillis(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

// FromUnixMillis converts Unix timestamp in milliseconds to time
func FromUnixMillis(millis int64) time.Time {
	return time.Unix(0, millis*int64(time.Millisecond))
}

// Truncate truncates time to the specified duration
func Truncate(t time.Time, d time.Duration) time.Time {
	return t.Truncate(d)
}

// Round rounds time to the nearest duration
func Round(t time.Time, d time.Duration) time.Time {
	return t.Round(d)
}

// Max returns the later of two times
func Max(t1, t2 time.Time) time.Time {
	if t1.After(t2) {
		return t1
	}
	return t2
}

// Min returns the earlier of two times
func Min(t1, t2 time.Time) time.Time {
	if t1.Before(t2) {
		return t1
	}
	return t2
}

// InTimeRange checks if a time is within the given range (inclusive)
func InTimeRange(t, start, end time.Time) bool {
	return (t.Equal(start) || t.After(start)) && (t.Equal(end) || t.Before(end))
}
