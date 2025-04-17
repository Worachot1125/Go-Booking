package enum

import (
	"time"
)

type Day int

const (
	DAY_SUNDAY Day = iota + 1
	DAY_MONDAY
	DAY_TUESDAY
	DAY_WEDNESDAY
	DAY_THURSDAY
	DAY_FRIDAY
	DAY_SATURDAY
)

var (
	dayName = map[Day]string{
		DAY_SUNDAY:    "Sunday",
		DAY_MONDAY:    "Monday",
		DAY_TUESDAY:   "Tuesday",
		DAY_WEDNESDAY: "Wednesday",
		DAY_THURSDAY:  "Thursday",
		DAY_FRIDAY:    "Friday",
		DAY_SATURDAY:  "Saturday",
	}
)

func (s Day) List() map[Day]string {
	return dayName
}

func (s Day) String() string {
	return dayName[s]
}

func GetDay(text string) Day {
	switch text {
	case DAY_SUNDAY.String():
		return DAY_SUNDAY
	case DAY_MONDAY.String():
		return DAY_MONDAY
	case DAY_TUESDAY.String():
		return DAY_TUESDAY
	case DAY_WEDNESDAY.String():
		return DAY_WEDNESDAY
	case DAY_THURSDAY.String():
		return DAY_THURSDAY
	case DAY_FRIDAY.String():
		return DAY_FRIDAY
	case DAY_SATURDAY.String():
		return DAY_SATURDAY
	default:
		return DAY_SUNDAY
	}
}

// ListDay returns a list of days
func ListDay() []Day {
	var list []Day
	for i := DAY_SUNDAY; i <= DAY_SATURDAY; i++ {
		list = append(list, i)
	}
	return list
}

// DateToDay converts time.Weekday to enum.Day
func DateToDay(date time.Time) Day {
	switch date.Weekday() {
	case time.Sunday:
		return DAY_SUNDAY
	case time.Monday:
		return DAY_MONDAY
	case time.Tuesday:
		return DAY_TUESDAY
	case time.Wednesday:
		return DAY_WEDNESDAY
	case time.Thursday:
		return DAY_THURSDAY
	case time.Friday:
		return DAY_FRIDAY
	case time.Saturday:
		return DAY_SATURDAY
	default:
		return DAY_SUNDAY
	}
}
