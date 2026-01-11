package nextdate

import (
	"errors"
	"strconv"
	"strings"
	"time"

	dateService "github.com/Yandex-Practicum/final/pkg/services/date"
)

const (
	countWeekDays = 7
)

func CheckIntervalTypeWeekDays(formatParts []string) error {
	if len(formatParts) < 2 || formatParts[1] == "" {
		return errors.New("неверный формат интервала")
	}

	weekDays := strings.Split(formatParts[1], ",")
	for _, weekDay := range weekDays {
		weekDay, err := strconv.Atoi(weekDay)
		if err != nil {
			return errors.New("неверный формат интервала")
		}
		if weekDay <= 0 {
			return errors.New("передан нулевой или отрицательный интервал")
		}

		if weekDay > countWeekDays {
			return errors.New("недопустимый день недели")
		}
	}

	return nil
}

func GetNextDateByWeekDays(now, date time.Time, days string) time.Time {
	weekDays := strings.Split(days, ",")

	var weekDaysNumbers []int

	for _, weekDay := range weekDays {
		monthDayNumber, _ := strconv.Atoi(weekDay)
		weekDaysNumbers = append(weekDaysNumbers, monthDayNumber)
	}

	nextDate, nextDateFound := getNextDateByWeekDaysList(now, date, int(now.Weekday()), weekDaysNumbers)
	if !nextDateFound {
		nextDate = getNextDateByFirstListDay(weekDaysNumbers[0], now, date)
	}

	return nextDate
}

// поиск даты по алгоритму, если текущий день недели меньше какого-либо дня из правила повторения
func getNextDateByWeekDaysList(now, date time.Time, currentWeekdayNumber int, weekDays []int) (time.Time, bool) {
	var nextDate time.Time
	nextDateFound := false

	for _, day := range weekDays {
		if day > currentWeekdayNumber {
			nextDate = getNextDateByInterval(now, date, 0, day-currentWeekdayNumber)
			nextDateFound = true
		}
	}

	return nextDate, nextDateFound
}

// если текущий день недели больше всех дней из правил повторения, то ищем дату с первым днем
// недели из правила повторения
func getNextDateByFirstListDay(firstListDay int, now, date time.Time) time.Time {
	for {
		date = date.AddDate(0, 0, 1)
		if int(date.Weekday()) == firstListDay && dateService.IsDateAfter(date, now) {
			return date
		}
	}
}
