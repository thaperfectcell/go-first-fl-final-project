package db

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type Task struct {
	ID      string `json:"id" db:"id"`
	Date    string `json:"date" db:"date"`
	Title   string `json:"title" db:"title"`
	Comment string `json:"comment" db:"comment"`
	Repeat  string `json:"repeat" db:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetTask(id string) (*Task, error) {
	row := DB.QueryRow(`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`, id)

	var (
		idVal   int64
		date    string
		title   string
		comment string
		repeat  string
	)

	if err := row.Scan(&idVal, &date, &title, &comment, &repeat); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, err
	}

	return &Task{
		ID:      fmt.Sprint(idVal),
		Date:    date,
		Title:   title,
		Comment: comment,
		Repeat:  repeat,
	}, nil
}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

func UpdateDate(next string, id string) error {
	res, err := DB.Exec(`UPDATE scheduler SET date = ? WHERE id = ?`, next, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

func DeleteTask(id string) error {
	res, err := DB.Exec(`DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

func Tasks(limit int, search string) ([]*Task, error) {
	if limit <= 0 {
		limit = 50
	}

	search = strings.TrimSpace(search)
	query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`
	args := []any{limit}

	if search != "" {
		if t, err := time.Parse("02.01.2006", search); err == nil {
			query = `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date LIMIT ?`
			args = []any{t.Format("20060102"), limit}
		} else {
			like := "%" + search + "%"
			query = `SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`
			args = []any{like, like, limit}
		}
	}

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]*Task, 0)
	for rows.Next() {
		var (
			id      int64
			date    string
			title   string
			comment string
			repeat  string
		)

		if err := rows.Scan(&id, &date, &title, &comment, &repeat); err != nil {
			return nil, err
		}

		tasks = append(tasks, &Task{
			ID:      fmt.Sprint(id),
			Date:    date,
			Title:   title,
			Comment: comment,
			Repeat:  repeat,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
