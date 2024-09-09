package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/maliven1/go_final_project/entity"
	_ "modernc.org/sqlite"
)

const Layout = "20060102"
const limit = 40
const searchLayout = "02.01.2006"

type DB struct {
	db *sql.DB
}

func GetTask(db DB) ([]entity.AddTask, error) {
	var task entity.AddTask
	var tasks []entity.AddTask
	rows, err := db.db.Query(`SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT :limit`, sql.Named("limit", limit))
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
	if err = rows.Close(); err != nil {
		return []entity.AddTask{}, err
	}
	if len(tasks) == 0 {
		return []entity.AddTask{}, nil
	}
	return tasks, nil
}
func GetTaskSearch(db DB, search string) ([]entity.AddTask, error) {
	var task entity.AddTask
	var tasks []entity.AddTask

	searchData, err := time.Parse(searchLayout, search)
	if err == nil {
		search = searchData.Format(Layout)
	}

	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE :search OR date LIKE :search ORDER BY date LIMIT :limit`
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
	if err = rows.Close(); err != nil {
		return []entity.AddTask{}, err
	}
	if len(tasks) == 0 {
		return []entity.AddTask{}, nil
	}
	return tasks, nil

}
func (db DB) Close() {
	db.db.Close()
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

func (db DB) AddTask(task entity.AddTask) (int64, error) {
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

func New() (DB, error) {
	var install bool
	if _, err := os.Stat(os.Getenv("TODO_DBFILE")); err != nil {
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
			date VARCHAR(8) ,
			title VARCHAR(128) NOT NULL,
			comment VARCHAR(256) ,
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
