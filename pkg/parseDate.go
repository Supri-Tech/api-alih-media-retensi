package pkg

import "time"

func ParseDate(date string) time.Time {
	formats := []string{
		"2006-01-02",
		"02/01/2006",
		"02-01-2006",
		"2006/01/02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, date); err == nil {
			return t
		}
	}

	return time.Time{}
}
