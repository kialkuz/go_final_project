package date

import "time"

const DateFormat = "20060102"

func IsDateAfter(firstDate, secondDate time.Time) bool {
	return firstDate.Format(DateFormat) > secondDate.Format(DateFormat)
}
