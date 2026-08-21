package api

import (
	"net/http"
	"os"
)

var todoPassword string

func Init() {
	todoPassword = os.Getenv("TODO_PASSWORD")

	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/signin", signinHandler)
	http.HandleFunc("/api/task", auth(taskHandler))
	http.HandleFunc("/api/task/done", auth(doneTaskHandler))
	http.HandleFunc("/api/tasks", auth(tasksHandler))
}
