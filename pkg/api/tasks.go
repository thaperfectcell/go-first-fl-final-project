package api

import (
	"net/http"

	"github.com/thaperfectcell/go-first-fl-final-project/pkg/db"
)

type tasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, map[string]string{"error": "method not allowed"})
		return
	}

	tasks, err := db.Tasks(50, r.FormValue("search"))
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	if tasks == nil {
		tasks = make([]*db.Task, 0)
	}

	writeJSON(w, tasksResp{Tasks: tasks})
}
