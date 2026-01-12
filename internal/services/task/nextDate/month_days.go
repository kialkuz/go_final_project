package nextDate

import (
	"errors"
	"slices"
	"strconv"
	"strings"
	"time"

	datePkg "github.com/Yandex-Practicum/final/pkg/date"
)

func checkIntervalTypeMonthDays(now time.Time, formatParts []string) error {
	monthDays := strings.Split(formatParts[1], ",")
	lastMonthDay := getMonthLastDay(now.Year(), int(now.Month()))
	hasMonthesList := len(formatParts) == 3

	for _, monthDay := range monthDays {
		monthDayNumber, err := strconv.Atoi(monthDay)
		if err != nil {
			return errors.New("неверный формат интервала")
		}
		if monthDayNumber < -2 || (!hasMonthesList && monthDayNumber > lastMonthDay) {
			return errors.New("недопустимый день месяца")
		}
	}

	if hasMonthesList {
		monthes := strings.Split(formatParts[2], ",")
		decemberNumber := int(time.December)

		for _, month := range monthes {
			monthNumber, err := strconv.Atoi(month)
			if err != nil {
				return errors.New("неверный формат списка месяцев")
			}

			if monthNumber < 1 || monthNumber > decemberNumber {
				return errors.New("недопустимый месяц")
			}

			lastMonthDay := getMonthLastDay(now.Year(), monthNumber)

			for _, monthDay := range monthDays {
				monthDayNumber, _ := strconv.Atoi(monthDay)
				if monthDayNumber > lastMonthDay {
					return errors.New("недопустимый день месяца")
				}
			}
		}
	}

	return nil
}

func GetNextDateByMonthDays(now, date time.Time, days, monthes string) time.Time {
	monthDays := strings.Split(days, ",")

	var monthDaysNumbers []int

	for _, monthDay := range monthDays {
		monthDayNumber, _ := strconv.Atoi(monthDay)
		monthDaysNumbers = append(monthDaysNumbers, monthDayNumber)
	}

	monthDaysNumbers = sortMonthDaysNumbers(monthDaysNumbers)

	if monthes == "" {
		var nextDate time.Time
		var nextDateFound bool
		if datePkg.IsDateAfter(date, now) {
			nextDate, nextDateFound = getNextDateMonthOfDateStart(now, date, monthDaysNumbers)
		}

		// если месяц даты старта не включает нужную дату, то продолжаем поиск
		if !nextDateFound {
			for {
				date = date.AddDate(0, 0, 1)

				if datePkg.IsDateAfter(date, now) && date.Day() == monthDaysNumbers[0] {
					nextDate = getMonthDay(date.Year(), int(date.Month()), monthDaysNumbers[0])
					break
				}
			}
		}

		return nextDate
	} else {
		var monthNumbers []int

		monthList := strings.Split(monthes, ",")
		for _, month := range monthList {
			monthNumber, _ := strconv.Atoi(month)
			monthNumbers = append(monthNumbers, monthNumber)
		}

		slices.Sort(monthNumbers)

		dateStart := getDateStartForMonthesList(now, date, monthNumbers)
		nextDate, nextDateFound := getNextDateMonthOfDateStart(now, dateStart, monthDaysNumbers)

		if !nextDateFound {
			nextDate = getMonthDay(dateStart.Year(), int(dateStart.Month()), monthDaysNumbers[0])
		}

		return nextDate
	}
}

// сортировка дней из правила повторения
// сначала по возрастанию положительные числа
// затем отрицательные, т.к. -2, -1 - последние дни месяца
func sortMonthDaysNumbers(monthDaysNumbers []int) []int {
	var positiveNumbers []int
	var negativeNumbers []int

	for _, monthDaysNumber := range monthDaysNumbers {
		if monthDaysNumber > 0 {
			positiveNumbers = append(positiveNumbers, monthDaysNumber)
		} else {
			negativeNumbers = append(negativeNumbers, monthDaysNumber)
		}
	}

	slices.Sort(positiveNumbers)
	slices.Sort(negativeNumbers)

	return append(positiveNumbers, negativeNumbers...)
}

// получение даты начала поиска при передаче списка месяцев
func getDateStartForMonthesList(now, dateStart time.Time, monthes []int) time.Time {
	currentMonthNumber := int(now.Month())
	dateStartMonthNumber := int(dateStart.Month())
	isDateAfter := datePkg.IsDateAfter(dateStart, now)

	for _, monthNumber := range monthes {
		if currentMonthNumber == monthNumber && dateStartMonthNumber == monthNumber && isDateAfter {
			return dateStart
		} else if currentMonthNumber < monthNumber {
			return getMonthDay(now.Year(), monthNumber, 1)
		}
	}

	return getMonthDay(now.Year()+1, monthes[0], 1)
}

// получение следующей даты в рамках месяца даты старта
func getNextDateMonthOfDateStart(now, date time.Time, monthDaysNumbers []int) (time.Time, bool) {
	var nextDate time.Time

	dateStartMonthDay := date.Day()
	nextDateFound := false

	lastMonthDay := getMonthLastDay(date.Year(), int(date.Month()))

	for _, dayNumber := range monthDaysNumbers {
		if lastMonthDay < dayNumber {
			continue
		}

		if dateStartMonthDay < dayNumber {
			nextDate = getNextDateByInterval(now, date, 0, dayNumber-dateStartMonthDay)
			nextDateFound = true
			break
		} else if dayNumber == -2 || dayNumber == -1 {
			prevMonthDay := lastMonthDay - 1

			if dayNumber == -2 && dateStartMonthDay < prevMonthDay {
				nextDate = getMonthDay(date.Year(), int(date.Month()), prevMonthDay)
				nextDateFound = true
				break
			} else if dayNumber == -1 && dateStartMonthDay < lastMonthDay {
				nextDate = getMonthDay(date.Year(), int(date.Month()), lastMonthDay)
				nextDateFound = true
				break
			}
		}
	}

	return nextDate, nextDateFound
}
