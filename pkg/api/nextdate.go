package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func afterNow(date, now time.Time) bool {
	d := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	n := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	return d.After(n)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	repeat = strings.TrimSpace(repeat)
	if repeat == "" {
		return "", errors.New("empty repeat")
	}

	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", err
	}

	rules := strings.Fields(repeat)
	if len(rules) == 0 {
		return "", errors.New("invalid repeat format")
	}

	switch rules[0] {
	case "d":
		if len(rules) != 2 {
			return "", errors.New("invalid repeat format")
		}
		interval, err := strconv.Atoi(rules[1])
		if err != nil || interval < 1 || interval > 400 {
			return "", errors.New("invalid repeat format")
		}
		date = date.AddDate(0, 0, interval)
		for !afterNow(date, now) {
			date = date.AddDate(0, 0, interval)
		}
		return date.Format(dateFormat), nil

	case "y":
		if len(rules) != 1 {
			return "", errors.New("invalid repeat format")
		}
		date = date.AddDate(1, 0, 0)
		for !afterNow(date, now) {
			date = date.AddDate(1, 0, 0)
		}
		return date.Format(dateFormat), nil

	case "w":
		if len(rules) != 2 {
			return "", errors.New("invalid repeat format")
		}
		dayAllowed := [8]bool{}
		for _, part := range strings.Split(rules[1], ",") {
			if part == "" {
				return "", errors.New("invalid repeat format")
			}
			d, err := strconv.Atoi(part)
			if err != nil || d < 1 || d > 7 {
				return "", errors.New("invalid repeat format")
			}
			dayAllowed[d] = true
		}

		date = date.AddDate(0, 0, 1)
		for {
			weekday := int(date.Weekday())
			if weekday == 0 {
				weekday = 7
			}
			if afterNow(date, now) && dayAllowed[weekday] {
				return date.Format(dateFormat), nil
			}
			date = date.AddDate(0, 0, 1)
		}

	case "m":
		if len(rules) != 2 && len(rules) != 3 {
			return "", errors.New("invalid repeat format")
		}

		dayAllowed := map[int]bool{}
		for _, part := range strings.Split(rules[1], ",") {
			if part == "" {
				return "", errors.New("invalid repeat format")
			}
			d, err := strconv.Atoi(part)
			if err != nil {
				return "", errors.New("invalid repeat format")
			}
			if (d < 1 || d > 31) && d != -1 && d != -2 {
				return "", errors.New("invalid repeat format")
			}
			dayAllowed[d] = true
		}

		monthAllowed := [13]bool{}
		if len(rules) == 3 {
			for _, part := range strings.Split(rules[2], ",") {
				if part == "" {
					return "", errors.New("invalid repeat format")
				}
				m, err := strconv.Atoi(part)
				if err != nil || m < 1 || m > 12 {
					return "", errors.New("invalid repeat format")
				}
				monthAllowed[m] = true
			}
		} else {
			for i := 1; i <= 12; i++ {
				monthAllowed[i] = true
			}
		}

		date = date.AddDate(0, 0, 1)
		for {
			if afterNow(date, now) && monthAllowed[int(date.Month())] {
				day := date.Day()
				lastDay := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
				if dayAllowed[day] || (dayAllowed[-1] && day == lastDay) || (dayAllowed[-2] && day == lastDay-1) {
					return date.Format(dateFormat), nil
				}
			}
			date = date.AddDate(0, 0, 1)
		}
	}

	return "", fmt.Errorf("unsupported repeat format")
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	nowParam := strings.TrimSpace(r.FormValue("now"))
	date := strings.TrimSpace(r.FormValue("date"))
	repeat := strings.TrimSpace(r.FormValue("repeat"))

	now := time.Now()
	if nowParam != "" {
		parsed, err := time.Parse(dateFormat, nowParam)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		now = parsed
	}

	next, err := NextDate(now, date, repeat)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(next))
}
