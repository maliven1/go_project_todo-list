package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
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
	Layout := "20060102"
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

	default:
		return "", fmt.Errorf("некорректный тип правила")
	}

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
