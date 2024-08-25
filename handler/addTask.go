package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/maliven1/go_final_project/db"
	"github.com/maliven1/go_final_project/entity"
)

const Layout = "20060102"

type ErrorResponse struct {
	Message string `json:"error"`
}
type ScResponse struct {
	ID string `json:"id"`
}
type TaskResponse struct {
	Task []entity.AddTask `json:"tasks"`
}

func NewTaskHandler(db db.DB) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		var task entity.Task
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
		if task.Repeat == "" {
			task.Date = now.Format(Layout)
		} else {
			_, err := nextDate(now, task.Date, task.Repeat)
			if err != nil {
				responseWhithError(res, "Не верное условие повторения")
				return
			}
		}

		id, err := db.AddTask(task)
		if err != nil {
			log.Panicln(err)
			r, _ := json.Marshal(ErrorResponse{Message: "Ошибка при получение id"})
			res.WriteHeader(http.StatusBadRequest)
			res.Write(r)
			return
		}
		idRes := strconv.Itoa(int(id))

		responseWhithOk(res, idRes)

	}
}
func NextDateHandler(res http.ResponseWriter, req *http.Request) {
	now := req.FormValue("now")
	date := req.FormValue("date")
	repeat := req.FormValue("repeat")

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")

	nowTime, err := time.Parse(Layout, now)
	if err != nil {
		http.Error(res, "Некорректный формат даты", http.StatusBadRequest)
		return
	}
	nextDate, err := nextDate(nowTime, date, repeat)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	// Возвращаем ответ
	_, err = res.Write([]byte(nextDate))
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
}

func nextDate(now time.Time, date string, repeat string) (string, error) {
	//проверка на пустой repeat
	if repeat == "" {
		return "", fmt.Errorf("не указан repeat")
	}

	startDate, err := time.Parse(Layout, date)
	//проверка на неверный формат даты
	if err != nil {
		return "", fmt.Errorf("неверный формат даты: %v", err)
	}
	//разделяем правило на тип правила и его аргумент
	ruleSplited := strings.Split(repeat, " ")
	//тип правила
	ruleType := ruleSplited[0]

	switch ruleType {
	case "d":
		if len(ruleSplited) < 2 {
			return "", fmt.Errorf("не указано количество дней")
		}
		//количество дней для переноса задачи
		daysToMove, err := strconv.Atoi(ruleSplited[1])

		if err != nil {
			return "", err
		}
		if daysToMove > 400 {

			return "", fmt.Errorf("количество дней не должно превышать 400")
		}
		newDate := startDate.AddDate(0, 0, daysToMove)
		//проверяем, что возвращаемая дата не меньше текущей, если меньше - сдвигаем на указанное количество дней
		for newDate.Before(now) {
			newDate = newDate.AddDate(0, 0, daysToMove)
		}
		return newDate.Format(Layout), nil

	case "y":
		newDate := startDate.AddDate(1, 0, 0)
		//проверяем, что возвращаемая дата не меньше текущей, если меньше - сдвигаем еще на год
		for newDate.Before(now) {
			newDate = newDate.AddDate(1, 0, 0)
		}
		return newDate.Format(Layout), nil

	default:
		return "", fmt.Errorf("некорректный тип правила")
	}

}
func responseWhithOk(res http.ResponseWriter, s string) {
	r, _ := json.Marshal(ScResponse{ID: s})
	res.WriteHeader(http.StatusOK)
	res.Write(r)
}
func responseWhithError(res http.ResponseWriter, s string) {
	r, _ := json.Marshal(ErrorResponse{Message: s})
	res.WriteHeader(http.StatusBadRequest)
	res.Write(r)
}
func responseWithTasks(res http.ResponseWriter, s []entity.AddTask) {
	r, _ := json.Marshal(TaskResponse{Task: s})
	res.WriteHeader(http.StatusOK)
	res.Write(r)
}
