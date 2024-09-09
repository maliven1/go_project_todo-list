package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/maliven1/go_final_project/database"
	"github.com/maliven1/go_final_project/entity"
)

// доделать удаление
func ConfirmTaskHandler(db database.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var task entity.DeleteTask
		n := r.URL.Query().Get("id")
		_, err := strconv.Atoi(n)
		if err != nil {
			responseWhithError(w, "Не верный формат id")
			return
		}
		if n == "" {
			responseWhithError(w, "Не верный формат id")
			return
		}

		if _, err := db.ConfirmTask(n); err != nil {
			log.Println(err)
			responseWhithError(w, "Ошибка удаления")
			return

		}

		responseWithTasksConfirm(w, task)

	}

}
