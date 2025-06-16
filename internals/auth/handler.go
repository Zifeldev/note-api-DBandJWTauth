package auth

import (
	"context"
	"net/http"
	"note-manager-api/internals/db"
	"strings"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"fmt"
)

func Register(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil || user.Username == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "password hashing error"})
		return
	}

	_, err = db.Pool.Exec(context.Background(),
		"INSERT INTO users(username, password_hash) VALUES($1, $2)",
		user.Username, string(hashed))

	if err != nil {
		fmt.Println("DB insert error:", err)
		if strings.Contains(err.Error(), "duplicate") {
			c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user created"})
}

func Login(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil || user.Username == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	var id int
	var hashed string
	err := db.Pool.QueryRow(context.Background(),
		"SELECT id, password_hash FROM users WHERE username=$1", user.Username).
		Scan(&id, &hashed)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid username or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(user.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid username or password"})
		return
	}

	token, err := GenerateJWT(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
