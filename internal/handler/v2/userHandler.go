package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

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

		r.Put("/activate", hdl.ActivateUser)
		r.Get("/profile", hdl.GetProfile)
		r.Put("/profile", hdl.UpdateProfile)
		r.Put("/change-password", hdl.ChangePassword)
	})
}

func (hdl *UserHandler) UserAdminRoutes(router chi.Router) {
	router.Group(func(r chi.Router) {
		r.Use(middleware.VerifyToken)
		r.Use(middleware.VerifyAdmin)
		r.Get("/users", hdl.GetAllUsers)
		r.Put("/users/{id}", hdl.UpdateUser)
		r.Patch("/users/{id}/status", hdl.ToggleStatus)
		r.Put("/users/{id}/activate", hdl.ActivateUser)
	})
}

func (hdl *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

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
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

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
}

// func (hdl *UserHandler) Activate(w http.ResponseWriter, r *http.Request) {
// 	var req struct {
// 		Email string `json:"email"`
// 	}

// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		pkg.Error(w, http.StatusBadRequest, "Invalid request body")
// 		return
// 	}

// 	activatedUser, err := hdl.service.Activation(r.Context(), req.Email)
// 	if err != nil {
// 		status := http.StatusBadRequest
// 		if err.Error() == "User not found" {
// 			status = http.StatusNotFound
// 		}
// 		pkg.Error(w, status, err.Error())
// 		return
// 	}

// 	pkg.Success(w, "User activated", activatedUser)
// }

func (hdl *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := pkg.GetUserIDFromCtx(r.Context())

	user, err := hdl.service.GetProfile(r.Context(), userID)
	if err != nil {
		pkg.Error(w, http.StatusNotFound, err.Error())
		return
	}

	resp := models.UserResponse{
		ID:     user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Role:   user.Role,
		Status: user.Status,
	}

	pkg.Success(w, "Get profile success", resp)
}

func (hdl *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	userID := pkg.GetUserIDFromCtx(r.Context())
	user := models.User{
		ID:    userID,
		Name:  req.Name,
		Email: req.Email,
	}

	updated, err := hdl.service.UpdateProfile(r.Context(), user)
	if err != nil {
		pkg.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	pkg.Success(w, "Profile updated", updated)
}

func (hdl *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	userID := pkg.GetUserIDFromCtx(r.Context())
	err := hdl.service.ChangePassword(r.Context(), userID, req.OldPassword, req.NewPassword)
	if err != nil {
		pkg.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	pkg.Success(w, "Password changed successfully", nil)
}

func (hdl *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	page := 1
	perPage := 10
	users, err := hdl.service.GetAll(r.Context(), page, perPage)
	if err != nil {
		pkg.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	pkg.Success(w, "List users", users)
}

func (hdl *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		pkg.Error(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req struct {
		Name   string `json:"name"`
		Email  string `json:"email"`
		Role   string `json:"role"`
		Status string `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		pkg.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	updated, err := hdl.service.Update(r.Context(), models.User{
		ID:     id,
		Name:   req.Name,
		Email:  req.Email,
		Role:   req.Role,
		Status: req.Status,
	})
	if err != nil {
		pkg.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	pkg.Success(w, "User updated", updated)
}
func (hdl *UserHandler) ToggleStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		pkg.Error(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := hdl.service.GetProfile(r.Context(), id)
	if err != nil {
		pkg.Error(w, http.StatusNotFound, err.Error())
		return
	}

	var updated *models.User
	if user.Status == "aktif" {
		updated, err = hdl.service.UpdateStatus(r.Context(), id, "tidak aktif")
	} else {
		updated, err = hdl.service.Activation(r.Context(), id)
	}
	if err != nil {
		pkg.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	pkg.Success(w, "User status updated", updated)
}

func (hdl *UserHandler) ActivateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		pkg.Error(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	activated, err := hdl.service.Activation(r.Context(), id)
	if err != nil {
		pkg.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	pkg.Success(w, "User activated", activated)
}

// func (hdl *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
// 	idStr := chi.URLParam(r, "id")
// 	// id, err := strconv.Atoi(idStr)
// 	_, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		pkg.Error(w, http.StatusBadRequest, "Invalid user ID")
// 		return
// 	}

// 	var req struct {
// 		Name   string `json:"name"`
// 		Email  string `json:"email"`
// 		Role   string `json:"role"`
// 		Status string `json:"status"`
// 	}

// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		pkg.Error(w, http.StatusBadRequest, "Invalid request body")
// 		return
// 	}

// 	// Pakai Update dan UpdateStatus sesuai kebutuhan
// 	updated, err := hdl.service.Update(r.Context(), req.Name, req.Email, "", req.Role, req.Status)
// 	if err != nil {
// 		pkg.Error(w, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	pkg.Success(w, "User updated", updated)
// }

// func (hdl *UserHandler) ToggleStatus(w http.ResponseWriter, r *http.Request) {
// 	idStr := chi.URLParam(r, "id")
// 	id, _ := strconv.Atoi(idStr)

// 	user, err := hdl.service.GetProfile(r.Context(), id)
// 	if err != nil {
// 		pkg.Error(w, http.StatusNotFound, err.Error())
// 		return
// 	}

// 	var updated *models.User
// 	if user.Status == "aktif" {
// 		updated, err = hdl.service.UpdateStatus(r.Context(), models.User{Email: user.Email, Status: "tidak aktif"})
// 	} else {
// 		updated, err = hdl.service.Activation(r.Context(), user.Email)
// 	}
// 	if err != nil {
// 		pkg.Error(w, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	pkg.Success(w, "User status updated", updated)
// }

// func (hdl *UserHandler) ActivateUser(w http.ResponseWriter, r *http.Request) {
// 	idStr := chi.URLParam(r, "id")
// 	id, _ := strconv.Atoi(idStr)

// 	user, err := hdl.service.GetProfile(r.Context(), id)
// 	if err != nil {
// 		pkg.Error(w, http.StatusNotFound, err.Error())
// 		return
// 	}

// 	activated, err := hdl.service.Activation(r.Context(), user.Email)
// 	if err != nil {
// 		pkg.Error(w, http.StatusBadRequest, err.Error())
// 		return
// 	}
// 	pkg.Success(w, "User activated", activated)
// }

// package handler

// import (
// 	"encoding/json"
// 	"net/http"

// 	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/middleware"
// 	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/models/v2"
// 	"github.com/cukiprit/api-sistem-alih-media-retensi/internal/services/v2"
// 	"github.com/cukiprit/api-sistem-alih-media-retensi/pkg"
// 	"github.com/go-chi/chi/v5"
// )

// type UserHandler struct {
// 	service services.UserService
// }

// func NewUserHandler(service services.UserService) *UserHandler {
// 	return &UserHandler{service: service}
// }

// func (hdl *UserHandler) UserRoutes(router chi.Router) {
// 	router.Post("/login", hdl.Login)
// 	router.Post("/register", hdl.Register)

// 	router.Group(func(r chi.Router) {
// 		r.Use(middleware.VerifyToken)

// 		r.Put("/activate", hdl.Activate)
// 		r.Get("/profile", hdl.GetProfile)
// 		r.Put("/profile", hdl.UpdateProfile)
// 		r.Put("/change-password", hdl.ChangePassword)
// 	})
// }

// func (hdl *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
// 	type Login struct {
// 		Email    string `json:"email"`
// 		Password string `json:"password"`
// 	}

// 	var req Login
// 	// email, password, ok := r.BasicAuth()
// 	// if !ok {
// 	// 	pkg.Error(w, http.StatusUnauthorized, "Authorization header missing or invalid")
// 	// 	return
// 	// }
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		pkg.Error(w, http.StatusBadRequest, "Invalid request body")
// 		return
// 	}

// 	token, err := hdl.service.Login(r.Context(), req.Email, req.Password)
// 	if err != nil {
// 		pkg.Error(w, http.StatusUnauthorized, err.Error())
// 		return
// 	}

// 	http.SetCookie(w, &http.Cookie{
// 		Name:     "token",
// 		Value:    token,
// 		HttpOnly: true,
// 		Secure:   true,
// 		SameSite: http.SameSiteStrictMode,
// 		Path:     "/",
// 	})

// 	pkg.Success(w, "Login success", token)
// }

// func (hdl *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
// 	type Register struct {
// 		Name     string `json:"name"`
// 		Email    string `json:"email"`
// 		Password string `json:"password"`
// 	}

// 	var req Register
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		pkg.Error(w, http.StatusBadRequest, "Invalid request body")
// 		return
// 	}

// 	user := models.User{
// 		Name:     req.Name,
// 		Email:    req.Email,
// 		Password: req.Password,
// 		Role:     "user",
// 		Status:   "tidak aktif",
// 	}

// 	newUser, err := hdl.service.Register(r.Context(), user)
// 	if err != nil {
// 		status := http.StatusBadRequest
// 		if err.Error() == "Email already used" {
// 			status = http.StatusConflict
// 		}
// 		pkg.Error(w, status, err.Error())
// 		return
// 	}

// 	pkg.Success(w, "User registered", newUser)
// 	// w.WriteHeader(http.StatusCreated)
// 	// json.NewEncoder(w).Encode(newUser)
// }

// func (hdl *UserHandler) Activate(w http.ResponseWriter, r *http.Request) {
// 	type Activate struct {
// 		Email string `json:"email"`
// 	}
// 	var req Activate
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		pkg.Error(w, http.StatusBadRequest, "Invalid request body")
// 		return
// 	}

// 	activatedUser, err := hdl.service.Activation(r.Context(), req.Email)
// 	if err != nil {
// 		status := http.StatusBadRequest
// 		if err.Error() == "user not found" {
// 			status = http.StatusNotFound
// 		}
// 		pkg.Error(w, status, err.Error())
// 		return
// 	}

// 	pkg.Success(w, "Data activated", activatedUser)
// 	// json.NewEncoder(w).Encode(activatedUser)
// }

// func (hdl *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
// 	userID := pkg.GetUserIDFromCtx(r.Context())

// 	user, err := hdl.service.GetProfile(r.Context(), userID)
// 	if err != nil {
// 		pkg.Error(w, http.StatusNotFound, err.Error())
// 		return
// 	}

// 	resp := models.UserResponse{
// 		ID:     user.ID,
// 		Name:   user.Name,
// 		Email:  user.Email,
// 		Role:   user.Role,
// 		Status: user.Status,
// 	}

// 	pkg.Success(w, "Get profile success", resp)
// }

// func (hdl *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
// 	type Req struct {
// 		Name  string `json:"name"`
// 		Email string `json:"email"`
// 	}
// 	var req Req
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		pkg.Error(w, http.StatusBadRequest, "Invalid request body")
// 		return
// 	}

// 	userID := pkg.GetUserIDFromCtx(r.Context())
// 	user := models.User{
// 		ID:    userID,
// 		Name:  req.Name,
// 		Email: req.Email,
// 	}

// 	updated, err := hdl.service.UpdateProfile(r.Context(), user)
// 	if err != nil {
// 		pkg.Error(w, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	pkg.Success(w, "Profile updated", updated)
// }

// func (hdl *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
// 	type Req struct {
// 		OldPassword string `json:"old_password"`
// 		NewPassword string `json:"new_password"`
// 	}
// 	var req Req
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		pkg.Error(w, http.StatusBadRequest, "Invalid request body")
// 		return
// 	}

// 	userID := pkg.GetUserIDFromCtx(r.Context())
// 	err := hdl.service.ChangePassword(r.Context(), userID, req.OldPassword, req.NewPassword)
// 	if err != nil {
// 		pkg.Error(w, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	pkg.Success(w, "Password changed successfully", nil)
// }
