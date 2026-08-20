package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/thaperfectcell/go-first-fl-final-project/pkg/api"
)

const webDir = "./web"

func Run() error {
	api.Init()
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	fmt.Println("Запуск сервера на порту", port)
	return http.ListenAndServe(":"+port, nil)
}
