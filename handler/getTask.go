package handler

import (
	"net/http"

	database "github.com/maliven1/go_final_project/db"
)

func GetTasks(d database.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		search := r.URL.Query().Get("search")
		if search == "" {
			tasks, err := database.GetTask(d)
			if err != nil {
				responseWhithError(w, "Ошибка при получении данных")
				return
			}
			responseWithTasks(w, tasks)
		} else {
			tasks, err := database.GetTaskSearch(d, search)
			if err != nil {
				responseWhithError(w, "Ошибка при получении данных")
				return
			}
			responseWithTasks(w, tasks)
		}

	}
}
