package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	database "github.com/maliven1/go_final_project/db"
	"github.com/maliven1/go_final_project/handler"
	"github.com/maliven1/go_final_project/middlewares"
	_ "modernc.org/sqlite"
)

const webDir = "./web"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db, err := database.New()
	if err != nil {
		log.Fatalf("Ошибка создания bd")
	}
	defer db.Close()
	r := chi.NewRouter()

	fs := http.FileServer(http.Dir(webDir))

	r.Handle("/*", http.StripPrefix("/", fs))
	r.Group(func(r chi.Router) {
		r.Use(middlewares.NewAuthMeddlewares())
		r.Post("/api/task", handler.NewTaskHandler(db))
		r.Get("/api/tasks", handler.GetTasks(db))
		r.Get("/api/task", handler.GetTasksParam(db))
		r.Put("/api/task", handler.UpdateTaskHandler(db))
		r.Post("/api/task/done", handler.ConfirmTaskHandler(db))
		r.Delete("/api/task", handler.DeleteTaskHandler(db))
	})
	r.Get("/api/nextdate", handler.NextDateHandler)
	r.Post("/api/signin", handler.AuthorizationGenerateToken)
	go func() {
		log.Printf("Starting server on :%s\n", os.Getenv("TODO_PORT"))
		if err := http.ListenAndServe(":"+os.Getenv("TODO_PORT"), r); err != nil {
			panic(err)
		}
	}()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	db.Close()

}
