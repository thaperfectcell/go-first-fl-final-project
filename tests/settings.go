package tests

import (
	"os"
	"github.com/thaperfectcell/go-first-fl-final-project/pkg/api"
)

var Port = 7540
var DBFile = "../scheduler.db"
var FullNextDate = true
var Search = true
var Token string

func init() {
    password := os.Getenv("TODO_PASSWORD")
    if password == "" {
        Token = ""
        return
    }

    token, _ := api.CreateToken(password)
    Token = token
}
