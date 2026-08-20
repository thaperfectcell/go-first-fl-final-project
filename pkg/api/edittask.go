package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/thaperfectcell/go-first-fl-final-project/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.FormValue("id"))
	if id == "" {
		writeJSON(w, map[string]string{"error": "id is required"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, task)
}

func putTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	task.ID = strings.TrimSpace(task.ID)
	task.Title = strings.TrimSpace(task.Title)
	if task.ID == "" {
		writeJSON(w, map[string]string{"error": "id is required"})
		return
	}
	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "title is required"})
		return
	}

	if err := checkDate(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, map[string]string{})
}
