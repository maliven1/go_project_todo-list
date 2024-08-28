package handler

import (
	"net/http"
	"strconv"

	database "github.com/maliven1/go_final_project/db"
)

func DeleteTaskHandler(db database.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
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

		if _, err := db.DeleteID(n); err != nil {
			responseWhithError(w, "Ошибка удаления")
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	}

}
