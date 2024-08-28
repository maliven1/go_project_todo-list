package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	database "github.com/maliven1/go_final_project/db"
	"github.com/maliven1/go_final_project/entity"
)

// сократить повторяющийся код (addTask)
func UpdateTaskHandler(db database.DB) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		var task entity.AddTask
		var buf bytes.Buffer
		res.Header().Set("Content-Type", "application/json; charset=UTF-8")
		now := time.Now()
		_, err := buf.ReadFrom(req.Body)
		if err != nil {
			responseWhithError(res, "Ошибка чтения")
			return
		}
		if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
			responseWhithError(res, "Ошибка чтения")
			return
		}
		dbID, err := db.DataBaseID(task)
		if err != nil {
			responseWhithError(res, "Ошибка работы db")
			return
		}
		reqID, err := strconv.Atoi(task.ID)
		if err != nil {
			responseWhithError(res, "Ошибка получения id")
			return
		}
		if reqID != dbID {
			responseWhithError(res, "нет такого id")
			return
		}
		if task.Title == "" {
			responseWhithError(res, "Не указан заголовок задачи")
			return
		}
		if task.Date == "" {
			task.Date = now.Format(Layout)
		}
		if _, err = time.Parse(Layout, task.Date); err != nil {
			responseWhithError(res, "Не верный формат времени")
			return
		}
		if task.Date < now.Format(Layout) {
			task.Date = now.Format(Layout)
		}
		if task.Repeat != "" {
			_, err := database.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				responseWhithError(res, "Не верное условие повторения")
				return
			}
		}

		id, err := db.UpdateTask(task)
		if err != nil {
			r, _ := json.Marshal(ErrorResponse{Message: "Ошибка при получение id"})
			res.WriteHeader(http.StatusBadRequest)
			res.Write(r)
			return
		}
		idRes := strconv.Itoa(int(id))

		responseWhithOk(res, idRes)

	}
}
