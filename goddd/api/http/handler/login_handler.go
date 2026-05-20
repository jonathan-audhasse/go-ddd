package handler

import (
	"goddd/internal/application/user"
	"goddd/internal/infrastructure/auth"
)

type AuthHandler struct {
	userService user.UserService
	jwt         *auth.JWTService
}

// func (c *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
// 	var req dto.LoginRequest

// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		httperror.Write(w, err)
// 		return
// 	}

// 	user, err := c.userService.GetByUsername(r.Context(), req.Username)
// 	if err != nil {
// 		httperror.Write(w, err)
// 		return
// 	}

// 	if !checkPassword(req.Password, user.Password) {
// 		httperror.Write(w, errors.New("invalid credentials"))
// 		return
// 	}

// 	token, err := c.jwt.Generate(user.ID.String())
// 	if err != nil {
// 		httperror.Write(w, err)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	_ = json.NewEncoder(w).Encode(map[string]string{
// 		"token": token,
// 	})
// }
