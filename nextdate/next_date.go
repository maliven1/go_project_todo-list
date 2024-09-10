package nextdate

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

const Layout = "20060102"

func WeekComparison(newDate time.Time, daysToMove int) int {
	var Time = true
	n := 0
	for Time {
		n++
		newDate = newDate.AddDate(0, 0, 1)
		if int(newDate.Weekday()) == daysToMove {
			Time = false
		}
	}
	return n
}
func nextDateW(currentDate time.Time, possibleDays []int, possibleMonths []int) (time.Time, error) {
	sort.Ints(possibleMonths)

	currentDay := currentDate.Day()
	currentMonth := int(currentDate.Month())
	currentYear := currentDate.Year()

	var targetDay, targetMonth int
	found := false

	for _, month := range possibleMonths {
		normalizedDays := normalizeDays(possibleDays, currentYear, month)
		sort.Ints(normalizedDays)
		if month > currentMonth || (month == currentMonth && normalizedDays[len(normalizedDays)-1] > currentDay) {
			targetMonth = month
			for _, day := range normalizedDays {

				if month > currentMonth || day > currentDay {
					targetDay = day
					found = true
					break
				}
			}
			if found {
				break
			}
		}
	}

	if !found {
		targetMonth = possibleMonths[0]
		normalizedDays := normalizeDays(possibleDays, currentYear+1, targetMonth)
		sort.Ints(normalizedDays)
		targetDay = normalizedDays[0]
		currentYear++
	}

	for {
		lastDayOfMonth := time.Date(currentYear, time.Month(targetMonth+1), 0, 0, 0, 0, 0, currentDate.Location()).Day()
		if targetDay <= lastDayOfMonth {
			break
		}
		targetDay = normalizeDays(possibleDays, currentYear, targetMonth)[0]
		targetMonth++
		if targetMonth > 12 {
			targetMonth = 1
			currentYear++
		}
	}

	return time.Date(currentYear, time.Month(targetMonth), targetDay, 0, 0, 0, 0, currentDate.Location()), nil
}
func normalizeDays(days []int, year int, month int) []int {
	lastDayOfMonth := time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC).Day()
	normalizedDays := make([]int, len(days))
	for i, day := range days {

		if day < 0 {
			normalizedDays[i] = lastDayOfMonth + day + 1
		} else {
			normalizedDays[i] = day
		}
	}
	return normalizedDays
}

func NextDate(now time.Time, date string, repeat string) (string, error) {

	if repeat == "" {
		return "", fmt.Errorf("не указан repeat")
	}

	startDate, err := time.Parse(Layout, date)
	if err != nil {
		return "", fmt.Errorf("неверный формат даты: %v", err)
	}

	ruleSplited := strings.Split(repeat, " ")
	ruleType := ruleSplited[0]

	switch ruleType {
	case "d":
		if len(ruleSplited) < 2 {
			return "", fmt.Errorf("не указано количество дней")
		}

		daysToMove, err := strconv.Atoi(ruleSplited[1])

		if err != nil {
			return "", err
		}
		if daysToMove > 400 {

			return "", fmt.Errorf("количество дней не должно превышать 400")
		}
		newDate := startDate.AddDate(0, 0, daysToMove)
		for newDate.Before(now) {
			newDate = newDate.AddDate(0, 0, daysToMove)
		}
		return newDate.Format(Layout), nil

	case "y":
		newDate := startDate.AddDate(1, 0, 0)
		for newDate.Before(now) {
			newDate = newDate.AddDate(1, 0, 0)
		}
		return newDate.Format(Layout), nil

	case "w":

		if len(ruleSplited) < 2 {
			return "", fmt.Errorf("не указано количество дней")
		}
		SplitRule := strings.Split(ruleSplited[1], ",")
		for startDate.Year() < now.Year() {
			startDate = time.Date(now.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)
		}
		timeDur := startDate.AddDate(0, 0, 7) //edasdsa

		for _, v := range SplitRule {
			newDate := startDate
			daysToMove, err := strconv.Atoi(v)
			if err != nil {
				return "", err
			}

			if daysToMove == 7 {
				daysToMove = 0
			}
			if daysToMove >= 7 || daysToMove < 0 {
				return "", fmt.Errorf("Нет такого дня недели")
			}
			n := WeekComparison(newDate, daysToMove)
			newDate = startDate.AddDate(0, 0, n)
			timeDurint, _ := strconv.Atoi(timeDur.Format(Layout))
			newDateint, _ := strconv.Atoi(newDate.Format(Layout))
			if timeDurint > newDateint {
				timeDur = newDate
			}
		}

		return timeDur.Format(Layout), nil
	case "m":
		// now = time.Date(2024, time.January, 20, 0, 0, 0, 0, time.Local) //Для тестов, которые работают только с января 2024 года.
		// if startDate.Before(now) {
		// 	startDate = now
		// }
		if len(ruleSplited) > 3 {
			return "", fmt.Errorf("Не верный формат")
		}
		if len(ruleSplited) == 2 {
			var Month = "1,2,3,4,5,6,7,8,9,10,11,12"
			ruleSplited = append(ruleSplited, Month)
		}
		SplitRuleMonth := strings.Split(ruleSplited[2], ",") //месяца
		possibleMonths := make([]int, len(SplitRuleMonth))
		for i, v := range SplitRuleMonth {
			m, err := strconv.Atoi(v)
			if err != nil {
				return "", err
			}
			possibleMonths[i] = m
		}
		SplitRule := strings.Split(ruleSplited[1], ",") //дни
		possibleDays := make([]int, len(SplitRule))
		for i, v := range SplitRule {
			m, err := strconv.Atoi(v)
			if err != nil {
				return "", err
			}
			if m > 31 || m < -2 {
				return "", fmt.Errorf("Число больше 31")
			}
			possibleDays[i] = m
		}

		if len(ruleSplited) == 2 {
			res, err := nextDateW(startDate, possibleDays, possibleMonths)
			if err != nil {
				return "", err
			}
			return res.Format(Layout), nil
		}
		if len(ruleSplited) == 3 {
			res, err := nextDateW(startDate, possibleDays, possibleMonths)
			if err != nil {
				return "", err
			}
			return res.Format(Layout), nil

		}
		return "", nil
	default:
		return "", fmt.Errorf("некорректный тип правила")
	}

}
