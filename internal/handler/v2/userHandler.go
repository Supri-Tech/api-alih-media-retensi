package handler

import (
	"encoding/json"
	"net/http"

	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/middleware"
	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/models/v2"
	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/services/v2"
	"github.com/cukiprit/api-sistem-alih-media-retensi/pkg"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (hdl *UserHandler) UserRoutes(router chi.Router) {
	router.Post("/login", hdl.Login)
	router.Post("/register", hdl.Register)

	router.Group(func(r chi.Router) {
		r.Use(middleware.VerifyToken)

		router.Put("/activate", hdl.Activate)
	})
}

func (hdl *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	type Login struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req Login
	// email, password, ok := r.BasicAuth()
	// if !ok {
	// 	pkg.Error(w, http.StatusUnauthorized, "Authorization header missing or invalid")
	// 	return
	// }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	token, err := hdl.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		pkg.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	pkg.Success(w, "Login success", token)
}

func (hdl *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	type Register struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req Register
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     "user",
		Status:   "tidak aktif",
	}

	newUser, err := hdl.service.Register(r.Context(), user)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "Email already used" {
			status = http.StatusConflict
		}
		pkg.Error(w, status, err.Error())
		return
	}

	pkg.Success(w, "User registered", newUser)
	// w.WriteHeader(http.StatusCreated)
	// json.NewEncoder(w).Encode(newUser)
}

func (hdl *UserHandler) Activate(w http.ResponseWriter, r *http.Request) {
	type Activate struct {
		Email string `json:"email"`
	}
	var req Activate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	activatedUser, err := hdl.service.Activation(r.Context(), req.Email)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "user not found" {
			status = http.StatusNotFound
		}
		pkg.Error(w, status, err.Error())
		return
	}

	pkg.Success(w, "Data activated", activatedUser)
	// json.NewEncoder(w).Encode(activatedUser)
}
