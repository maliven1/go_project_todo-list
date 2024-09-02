package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	database "github.com/maliven1/go_final_project/db"
	"github.com/maliven1/go_final_project/handler"
	"github.com/maliven1/go_final_project/middlewares"
	_ "modernc.org/sqlite"
)

const webDir = "./web"
const port = "7540"

func main() {
	db, err := database.New()
	if err != nil {
		log.Fatalf("Ошибка создания bd")
	}
	defer db.Close()
	r := chi.NewRouter()

	fs := http.FileServer(http.Dir(webDir))

	r.Use(middlewares.NewAuthMeddlewares())

	r.Handle("/*", http.StripPrefix("/", fs))

	r.Get("/api/nextdate", handler.NextDateHandler)
	r.Post("/api/task", handler.NewTaskHandler(db))
	r.Get("/api/tasks", handler.GetTasks(db))
	r.Get("/api/task", handler.GetTasksParam(db))
	r.Put("/api/task", handler.UpdateTaskHandler(db))
	r.Post("/api/task/done", handler.ConfirmTaskHandler(db))
	r.Delete("/api/task", handler.DeleteTaskHandler(db))
	r.Post("/api/signin", handler.AuthorizationGenerateToken)

	log.Printf("Starting server on :%s\n", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		panic(err)
	}
}
