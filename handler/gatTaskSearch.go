package handler

import (
	"net/http"

	"github.com/maliven1/go_final_project/db"
)

func GetTaskss(d db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		tasks, err := db.GetTask(d)
		if err != nil {
			responseWhithError(w, "Ошибка при получении данных")
			return
		}
		responseWithTasks(w, tasks)
	}
}

// func GetTasksSearch(d db.DB) func(w http.ResponseWriter, r *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 	// 	search := r.URL.Query().Get("search")
// 	// 	tasks, err := db.GetTaskSearch(d, search)
// 	// 	if err != nil {
// 	// 		responseWhithError(w, "Ошибка при получении данных")
// 	// 		return
// 	// 	}
// 	// 	responseWithTasks(w, tasks)
// 	// }

// }
