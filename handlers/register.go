package handlers

import (
	"chater/models"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRepository struct {
	Db *sqlx.DB
}

func NewRegisterRepository(db *sqlx.DB) *RegisterRepository {
	return &RegisterRepository{Db: db}
}

func (c *RegisterRepository) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid method", http.StatusBadRequest)
		return
	}
	var user models.User

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "cant read request", http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, "invalid data", http.StatusBadRequest)
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		http.Error(w, "failed to generate hashcode", http.StatusInternalServerError)
		return
	}
	user.Password = string(hash)

	query := `INSERT INTO users (email, name, password) VALUES ($1, $2, $3)`
	_, err = c.Db.Exec(query, user.Email, user.Name, user.Password)
	if err != nil {
		http.Error(w, "failed to insert user", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "create user successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (c *RegisterRepository) SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid method", http.StatusBadRequest)
		return
	}

	var user models.User

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "cant read request", http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}
	hash := user.Password

	query := `SELECT * FROM users WHERE email = $1`

	row := c.Db.QueryRow(query, user.Email)

	err = row.Scan(&user.ID, &user.Email, &user.Name, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		} else {
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(hash)); err != nil {
		http.Error(w, "password dont match", http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"message": "signIn successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

}
