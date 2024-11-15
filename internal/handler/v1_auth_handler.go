package handler

import (
    "errors"
    "log"
    "net/http"
    "strings"

    "github.com/go-playground/validator/v10"
    "github.com/goccy/go-json"
    "github.com/yemyoaung/managing-vehicle-tracking-auth-svc/internal/services"
    "github.com/yemyoaung/managing-vehicle-tracking-common"
)

var (
    ErrUnauthorized     = errors.New("unauthorized")
    ErrMethodNotAllowed = errors.New("method was not allowed")
)

type V1AuthHandler struct {
    authService services.AuthService
    validate    *validator.Validate
}

func NewV1AuthHandler(authService services.AuthService, validate *validator.Validate) AuthHandler {
    return &V1AuthHandler{authService: authService, validate: validate}
}
func (h *V1AuthHandler) methodWasNotAllowed(w http.ResponseWriter) {
    common.HandleError(http.StatusMethodNotAllowed, w, ErrMethodNotAllowed)
}

// Login handles the login request
func (h *V1AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        h.methodWasNotAllowed(w)
        return
    }

    var req services.LoginRequest

    if body, ok := r.Context().Value(common.Body).([]byte); ok {
        if err := json.Unmarshal(body, &req); err != nil {
            common.HandleError(http.StatusUnprocessableEntity, w, err)
            return
        }
    }

    if err := h.validate.Struct(&req); err != nil {
        common.HandleError(http.StatusUnprocessableEntity, w, err)
        return
    }

    user, token, err := h.authService.Login(r.Context(), &req)
    if err != nil {
        common.HandleError(http.StatusUnprocessableEntity, w, err)
        return
    }

    resp := map[string]any{
        "token": token,
        "user":  user,
    }
    err = json.NewEncoder(w).Encode(common.DefaultSuccessResponse(resp, "login success"))
    if err != nil {
        log.Printf("Failed to encode response: %v", err)
    }
}

func (h *V1AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        h.methodWasNotAllowed(w)
        return
    }
    authorization := r.Header.Get(common.Authorization)
    if authorization == "" {
        common.HandleError(http.StatusUnauthorized, w, ErrUnauthorized)
        return
    }

    split := strings.Split(authorization, "Bearer ")

    if len(split) != 2 {
        common.HandleError(http.StatusUnauthorized, w, ErrUnauthorized)
        return
    }
    user, err := h.authService.ValidateToken(r.Context(), split[1])

    if err != nil {
        log.Println("ValidateToken Err: ", err)
        common.HandleError(http.StatusUnauthorized, w, ErrUnauthorized)
        return
    }

    if err = json.NewEncoder(w).Encode(common.DefaultSuccessResponse(user, "authenticated")); err != nil {
        log.Printf("Failed to encode response: %v", err)
    }
}
