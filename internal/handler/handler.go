package handler

import "net/http"

// AuthHandler is an interface for handling authentication
type AuthHandler interface {
    Login(w http.ResponseWriter, r *http.Request)
    Me(w http.ResponseWriter, r *http.Request)
}
