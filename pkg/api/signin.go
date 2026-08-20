package api

import (
	"encoding/json"
	"net/http"
	"os"
)

func signinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, map[string]string{"error": "method not allowed"})
		return
	}

	var req signinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	pass := os.Getenv("TODO_PASSWORD")
	if pass == "" {
		writeJSON(w, map[string]string{"error": "authentication is disabled"})
		return
	}

	if req.Password != pass {
		writeJSON(w, map[string]string{"error": "invalid password"})
		return
	}

	token, err := createToken(pass)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, map[string]string{"token": token})
}
