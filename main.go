package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/maliven1/go_final_project/db"
	"github.com/maliven1/go_final_project/handler"
	_ "modernc.org/sqlite"
)

const webDir = "./web"
const port = "7540"

func main() {
	db, err := db.New()
	if err != nil {
		log.Fatalf("Ошибка создания bd")
	}
	defer db.Close()
	r := chi.NewRouter()

	fs := http.FileServer(http.Dir(webDir))

	r.Handle("/*", http.StripPrefix("/", fs))
	r.HandleFunc("/api/nextdate", handler.NextDateHandler)
	r.Post("/api/task", handler.NewTaskHandler(db))
	r.Get("/api/tasks", handler.GetTasks(db))
	log.Printf("Starting server on :%s\n", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		panic(err)
	}
}
