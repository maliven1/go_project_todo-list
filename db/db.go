package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/maliven1/go_final_project/entity"
	_ "modernc.org/sqlite"
)

const Layout = "20060102"

type DB struct {
	db *sql.DB
}

func GetTaskSearch(db DB, search string) ([]entity.AddTask, error) {
	var task entity.AddTask
	var tasks []entity.AddTask
	limit := 40
	searchLayout := "02.01.2006"
	searchData, err := time.Parse(searchLayout, search)
	if err == nil {
		search = searchData.Format(Layout)
	}

	query := `SELECT * FROM scheduler WHERE title LIKE :search OR date LIKE :search ORDER BY date LIMIT :limit`
	rows, err := db.db.Query(query, sql.Named("search", `%`+search+`%`), sql.Named("limit", limit))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for rows.Next() {
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if len(task.Comment) == 0 {
			task.Comment = ""
		}

		if err != nil {
			log.Println(err)
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if len(tasks) == 0 {
		tasks = []entity.AddTask{}
		return tasks, nil
	}
	return tasks, nil

}
func GetTaskParam(db DB, param string) (entity.AddTask, error) {
	var task entity.AddTask
	_, err := strconv.Atoi(param)
	if err != nil {
		log.Println(err)
		return entity.AddTask{}, err
	}
	query := `SELECT * FROM scheduler WHERE id = :param`
	rows := db.db.QueryRow(query, sql.Named("param", param))
	if err != nil {
		log.Println(err)
		return entity.AddTask{}, err
	}
	if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.AddTask{}, fmt.Errorf("не найдена таска %w", err)
		}
	}

	return task, nil
}

func GetTask(db DB) ([]entity.AddTask, error) {
	var task entity.AddTask
	var tasks []entity.AddTask
	rows, err := db.db.Query(`SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT 40`)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for rows.Next() {
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if len(task.Comment) == 0 {
			task.Comment = ""
		}

		if err != nil {
			log.Println(err)
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if len(tasks) == 0 {
		tasks = []entity.AddTask{}
		return tasks, nil
	}
	return tasks, nil
}

func (db DB) AddTask(task entity.Task) (int64, error) {
	res, err := db.db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)", task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (db DB) DataBaseID(task entity.AddTask) (int, error) {
	query := "SELECT id FROM scheduler WHERE id = :id"
	rows, err := db.db.Query(query, sql.Named("id", task.ID))
	if err != nil {
		return 0, err
	}
	var id int
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return 0, nil
		}
	}
	return id, nil
}

func (db DB) DeleteID(id string) (entity.AddTask, error) {

	query := "DELETE FROM scheduler WHERE id = :id"
	_, err := db.db.Query(query, sql.Named("id", id))
	if err != nil {
		return entity.AddTask{}, err
	}

	return entity.AddTask{}, nil
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
		// now = time.Date(2024, time.January, 20, 0, 0, 0, 0, time.Local) //Для тество, которые работают только с января 2024 года.
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

func MonthComparison(startDate time.Time, SplitRule []string) (time.Time, error) {
	var res time.Time
	intSplitRule := make([]int, len(SplitRule))
	for i, rule := range SplitRule {
		m, err := strconv.Atoi(rule)
		if err != nil {
			return time.Time{}, err
		}
		intSplitRule[i] = m
	}
	sort.Ints(intSplitRule)
	startMonth := startDate.Month()
	var targetMonth int
	found := false
	for _, rule := range intSplitRule {
		m := rule
		if m > int(startMonth) {
			targetMonth = m
			found = true
			break
		}
	}

	if !found {
		targetMonth = intSplitRule[0]
	}

	res = time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, time.UTC)

	for int(res.Month()) != targetMonth {
		res = res.AddDate(0, 1, 0)

	}
	return res, nil
}

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

func (db DB) ConfirmTask(id string) (entity.AddTask, error) {
	var task entity.AddTask
	now := time.Now()
	query := "SELECT * FROM scheduler WHERE id = :id"
	rows, err := db.db.Query(query, sql.Named("id", id))
	if err != nil {
		return entity.AddTask{}, err
	}
	for rows.Next() {
		err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return entity.AddTask{}, err
		}
	}
	if task.Repeat == "" {
		query := "DELETE FROM scheduler WHERE id = :id"
		_, err := db.db.Query(query, sql.Named("id", id))
		if err != nil {
			return entity.AddTask{}, err
		}
		return entity.AddTask{}, nil
	}

	task.Date, err = NextDate(now, task.Date, task.Repeat)
	if err != nil {
		return entity.AddTask{}, err
	}
	_, err = db.UpdateTask(task)
	if err != nil {
		return entity.AddTask{}, err
	}
	return task, err
}

func (db DB) UpdateTask(task entity.AddTask) (int64, error) {
	query := "UPDATE  scheduler SET id = :id, date = :date, title = :title, comment= :comment, repeat= :repeat WHERE id = :id"
	res, err := db.db.Exec(query, sql.Named("id", task.ID), sql.Named("date", task.Date), sql.Named("title", task.Title), sql.Named("comment", task.Comment), sql.Named("repeat", task.Repeat))
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (db DB) Close() {
	db.db.Close()
}

func New() (DB, error) {
	var install bool
	if _, err := os.Stat("scheduler.db"); err != nil {
		if os.IsNotExist(err) {
			log.Println("База данных будет создана")
			install = true
		} else {
			log.Println("не получилось проверить файл")
			log.Fatal(err)
		}
	}
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		log.Fatal(err)
	}
	if install {
		_, err := db.Exec(`CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TEXT NOT NULL,
			title TEXT NOT NULL,
			comment TEXT ,
			repeat VARCHAR(128) NOT NULL
			);`)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := db.Exec(`CREATE INDEX index_date ON scheduler(date);`); err != nil {
			log.Fatal(err)
		} else {
			log.Println("База данных создана")
		}
	} else {
		fmt.Println("База данных уже существует")
	}

	return DB{db: db}, nil
}
