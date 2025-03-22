package handlers

import (
	"net/http"
)

type Register interface {
	SignIn(w http.ResponseWriter, r *http.Request)
	SignUp(w http.ResponseWriter, r *http.Request)
}

type Chat interface {
	SendMessage(w http.ResponseWriter, r *http.Request)
}
