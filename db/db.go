package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/maliven1/go_final_project/entity"
	_ "modernc.org/sqlite"
)

type DB struct {
	db *sql.DB
}

func GetTaskSearch(db DB, search string) {

	// limit := 40

	// query := `SELECT * FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT :limit`
	// rows, err := db.db.Query(query, sql.Named("search", search), sql.Named("limit", limit))
	// if err != nil {
	// 	log.Println(err)
	// 	return nil, err
	// }

	// for rows.Next() {
	// 	i := 0
	// 	i++
	// 	err := rows.Scan(&task.Date, &task.Title)
	// 	if len(task.Comment) == 0 {
	// 		task.Comment = ""
	// 	}

	// 	if err != nil {
	// 		log.Println(err)
	// 		return nil, err
	// 	}
	// 	tasks = append(tasks, task)
	// }
	// if len(tasks) == 0 {
	// 	tasks = []GetSearchTask{}
	// 	return tasks, nil
	// }
	return

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
		i := 0
		i++
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
