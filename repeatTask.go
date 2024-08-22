package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

func conv(i []string) []int {
	ints := make([]int, len(i))
	if i[0] == "d" || i[0] == "m" && len(i) == 2 {
		b, err := strconv.Atoi(i[1])
		if err != nil {
			fmt.Errorf("Не правильно заданы правила повторения дней")
		}
		ints[0] = b
		return ints
	} else if i[0] == "w" {
		str := i[1]
		s := strings.Split(str, ",")
		for b, v := range s {
			ints[b], _ = strconv.Atoi(v)
			return ints
		}
	} else if i[0] == "m" {
		str := i[1]
		s := strings.Split(str, ",")
		for b, v := range s {
			ints[b], _ = strconv.Atoi(v)
		}
		str = i[2]
		s = strings.Split(str, ",")
		for b, v := range s {
			ints[b], _ = strconv.Atoi(v)
		}
		return ints
	}
	return nil
}

const Layout = "20060102"

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if len(repeat) == 0 {
		return "", nil
	}
	i := strings.Split(repeat, " ")
	if i[0] == "d" && len(i) < 3 && len(i) != 1 {
		ints := conv(i)
		val := ints[0]
		if val > 400 {
			return "", fmt.Errorf("Не правильно заданы правила повторения дней")
		}
		nextTime := now.AddDate(0, 0, val)
		nextTimeString := nextTime.Format(Layout)
		return nextTimeString, nil
	} else if i[0] == "y" && len(i) == 1 {
		nextTime := now.AddDate(1, 0, 0)
		nextTimeString := nextTime.Format(Layout)
		return nextTimeString, nil
	} else if i[0] == "w" && len(i) < 3 && len(i) != 1 {

	} else if i[0] == "m" && len(i) <= 3 && len(i) != 1 {

	} else {
		log.Fatalf("Не правильный формат данных")
	}

	return repeat, nil
}
