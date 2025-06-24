package note

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"note-manager-api/internals/db"
	"note-manager-api/internals/redis"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetAllFav(c *gin.Context) {
	userID := c.GetInt("user_id")
	cacheKey := fmt.Sprintf("note-manager:favorites:%d", userID)


	if cached, err := redis.Client.Get(redis.Ctx, cacheKey).Result(); err == nil {
		var notes []Note
		if err := json.Unmarshal([]byte(cached), &notes); err == nil {
			c.JSON(http.StatusOK, notes)
			return
		}
	}


	rows, err := db.Pool.Query(context.Background(), `
		SELECT n.id, n.title, n.content 
		FROM notes n
		JOIN favorites f ON f.note_id = n.id
		WHERE f.user_id = $1`, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "db error"})
		return
	}
	defer rows.Close()

	var favorites []Note
	for rows.Next() {
		var note Note
		if err := rows.Scan(&note.ID, &note.Title, &note.Content); err != nil {
			continue
		}
		favorites = append(favorites, note)
	}


	data, _ := json.Marshal(favorites)
	redis.Client.Set(redis.Ctx, cacheKey, data, time.Minute*5)

	c.JSON(http.StatusOK, favorites)
}

func CreateFav(c *gin.Context) {
	userID := c.GetInt("user_id")
	noteIDStr := c.Param("note_id")
	noteID, err := strconv.Atoi(noteIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid note_id"})
		return
	}


	var exists bool
	err = db.Pool.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM notes WHERE id=$1 AND user_id=$2)",
		noteID, userID).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"message": "note not found or not yours"})
		return
	}


	_, err = db.Pool.Exec(context.Background(),
		"INSERT INTO favorites(user_id, note_id) VALUES($1, $2) ON CONFLICT DO NOTHING",
		userID, noteID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "insert into favorites failed"})
		return
	}


	cacheKey := fmt.Sprintf("note-manager:favorites:%d", userID)
	redis.Client.Del(redis.Ctx, cacheKey)

	c.JSON(http.StatusOK, gin.H{"message": "added to favorites"})
}

func DeleteFavNote(c *gin.Context) {
	userID := c.GetInt("user_id")
	noteIDStr := c.Param("note_id")
	noteID, err := strconv.Atoi(noteIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid note_id"})
		return
	}

	result, err := db.Pool.Exec(context.Background(),
		"DELETE FROM favorites WHERE note_id=$1 AND user_id=$2", noteID, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "db error"})
		return
	}

	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "favorite not found or not yours"})
		return
	}

	cacheKey := fmt.Sprintf("note-manager:favorites:%d", userID)
	redis.Client.Del(redis.Ctx, cacheKey)

	c.JSON(http.StatusOK, gin.H{"message": "favorite removed"})
}
