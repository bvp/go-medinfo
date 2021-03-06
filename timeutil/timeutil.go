// timeutil
package timeutil

import (
	"errors"
	"time"
)

// StartOfDay returns time at start of day of t.
func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// StartOfWeek returns time at start of week of t.
func StartOfWeek(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day()-int(t.Weekday()-1), 0, 0, 0, 0, t.Location())
}

func StartOfNextWeek(t time.Time) time.Time {
	return StartOfWeek(t).AddDate(0, 0, 7)
}

func WeekNumber(t time.Time) int {
	_, weekNum := t.ISOWeek()
	return weekNum
}

func NumberOfTheWeekInMonth(t time.Time) int {
	beginningOfTheMonth := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	_, thisWeek := t.ISOWeek()
	_, beginningWeek := beginningOfTheMonth.ISOWeek()
	return 1 + thisWeek - beginningWeek
}

func FormatDate(t time.Time) string {
	return t.Format("02-01-2006")
}

type Week struct {
	Days   []time.Time
	Year   int
	Number int
}

// NewWeek constructs new Week entity from given parameters (year and ISO-8601-compatible week number)
func NewWeek(params ...int) (*Week, error) {
	if len(params) < 2 {
		return &Week{}, errors.New("NewWeek(): too few arguments, specify year and number of week")
	} else if params[0] < 0 {
		return &Week{}, errors.New("NewWeek(): year can't be less than zero")
	} else if params[1] < 1 || params[1] > 53 {
		return &Week{}, errors.New("NewWeek(): number of week can't be less than 1 or greater than 53")
	}

	var (
		week = initWeek(params...)
		day  = 1
		fd   = time.Date(week.Year, 1, day, 0, 0, 0, 0, time.UTC)
		y, w = fd.ISOWeek()
	)

	for y != week.Year && w > 1 {
		day++
		fd = time.Date(week.Year, 1, day, 0, 0, 0, 0, time.UTC)
		y, w = fd.ISOWeek()
	}

	// getting Monday of the 1st week
	for fd.Weekday() > 1 {
		day--
		fd = time.Date(week.Year, 1, day, 0, 0, 0, 0, time.UTC)
	}

	// getting first day of the given week
	fd = fd.AddDate(0, 0, (week.Number-1)*7)
	for fd.Year() > y {
		fd = fd.AddDate(0, 0, -7)
	}

	// getting dates for whole week
	for i := 0; i < 7; i++ {
		week.Days = append(week.Days, fd.Add(time.Duration(i)*24*time.Hour))
	}

	return &week, nil
}

// Next calculates and returns information (year, week number and dates) about next week
func (week *Week) Next() (*Week, error) {
	var newYear, newWeek int
	if week.Number+1 > 53 {
		newYear = week.Year + 1
		newWeek = 1
	} else {
		newYear = week.Year
		newWeek = week.Number + 1
	}
	w, e := NewWeek(newYear, newWeek)

	return w, e
}

// Previous calculates and returns information (year, week number and dates) about previous week
func (week *Week) Previous() (*Week, error) {
	var newYear, newWeek int
	if week.Number-1 < 1 {
		newYear = week.Year - 1
		newWeek = 53
	} else {
		newYear = week.Year
		newWeek = week.Number - 1
	}
	w, e := NewWeek(newYear, newWeek)

	return w, e
}

func initWeek(params ...int) Week {
	var week = Week{
		Year:   params[0],
		Number: params[1],
	}
	return week
}
