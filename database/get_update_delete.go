package database

import (
	"database/sql"
	"time"

	"github.com/maliven1/go_final_project/entity"
	"github.com/maliven1/go_final_project/nextdate"
)

func (db DB) Delete(id string) error {

	query := "DELETE FROM scheduler WHERE id = :id"
	_, err := db.db.Query(query, sql.Named("id", id))
	if err != nil {
		return err
	}

	return nil
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
	if err = rows.Close(); err != nil {
		return entity.AddTask{}, err
	}
	if task.Repeat == "" {
		query := "DELETE FROM scheduler WHERE id = :id"
		_, err := db.db.Query(query, sql.Named("id", id))
		if err != nil {
			return entity.AddTask{}, err
		}
		return entity.AddTask{}, nil
	}

	task.Date, err = nextdate.NextDate(now, task.Date, task.Repeat)
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
