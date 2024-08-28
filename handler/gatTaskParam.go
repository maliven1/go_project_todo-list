package handler

import (
	"net/http"

	database "github.com/maliven1/go_final_project/db"
)

func GetTasksParam(d database.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		n := r.URL.Query().Get("id")
		tasks, err := database.GetTaskParam(d, n)
		if err != nil {
			responseWhithError(w, "Ошибка при получении данных")
			return
		}
		responseWithTasksT(w, tasks)
	}

}
