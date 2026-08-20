package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/thaperfectcell/go-first-fl-final-project/pkg/db"
)

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_ = json.NewEncoder(w).Encode(data)
}

func checkDate(task *db.Task) error {
	now := time.Now()
	today := now.Format(dateFormat)

	task.Date = strings.TrimSpace(task.Date)
	task.Repeat = strings.TrimSpace(task.Repeat)

	if task.Date == "" {
		task.Date = today
	}

	t, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return err
	}

	next := ""
	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	if afterNow(now, t) {
		if task.Repeat == "" {
			task.Date = today
		} else {
			task.Date = next
		}
	}

	return nil
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	task.Title = strings.TrimSpace(task.Title)
	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "title is required"})
		return
	}

	if err := checkDate(&task); err != nil {
		var parseErr *time.ParseError
		if errors.As(err, &parseErr) {
			writeJSON(w, map[string]string{"error": "invalid date format"})
			return
		}
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, map[string]string{"id": strconv.FormatInt(id, 10)})
}
