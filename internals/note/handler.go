package note

import (
	"context"
	"fmt"
	"net/http"
	"note-manager-api/internals/db"
	"note-manager-api/internals/redis"
	"github.com/gin-gonic/gin"
)


type Note struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func GetAll(c *gin.Context) {
	userID := c.GetInt("user_id")
	rows, err := db.Pool.Query(context.Background(), "SELECT id, title, content FROM notes WHERE user_id=$1", userID)
	if err != nil {
		fmt.Println("Error fetching notes:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "db error"})
		return
	}
	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var note Note
		if err := rows.Scan(&note.ID, &note.Title, &note.Content); err != nil {
			fmt.Println("Scan error:", err)
			continue
		}
		notes = append(notes, note)
	}
	if err := rows.Err(); err != nil {
		fmt.Println("Rows iteration error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "rows iteration error"})
		return
	}
	c.JSON(http.StatusOK, notes)
}

func CreateNote(c *gin.Context) {
	userID := c.GetInt("user_id")
	var note Note
	if err := c.BindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid body"})
		return
	}

	_, err := db.Pool.Exec(context.Background(),
		"INSERT INTO notes(user_id, title, content) VALUES($1, $2, $3)",
		userID, note.Title, note.Content)

	if err != nil {
		fmt.Println("Error creating note:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "inserting error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "note created"})
}

func UpdateNote(c *gin.Context){
	user_id := c.GetInt("user_id")
	note_id := c.Param("id")
	var note Note
	if err := c.BindJSON(&note); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"err":"invalid body"})
		return
	}

	result, err := db.Pool.Exec(context.Background(),"UPDATE notes SET title=$1, content=$2 WHERE id=$3 and user_id=$4",note.Title,note.Content,note_id,user_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":"error while updating"})
		return
	}

	if result.RowsAffected() == 0{
		c.JSON(http.StatusNotFound, gin.H{"message":"note not found or is not urs"})
		return
	}
	c.JSON(http.StatusOK,gin.H{"message":"updated"})
}

func DeleteNote(c *gin.Context) {
	userID := c.GetInt("user_id")
	noteID := c.Param("id")

	result, err := db.Pool.Exec(context.Background(),
		"DELETE FROM notes WHERE id=$1 AND user_id=$2", noteID, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "db error"})
		return
	}

	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "note not found or not yours"})
		return
	}
	redis.Client.Del(redis.Ctx, fmt.Sprintf("note-manager:favorites:%d", userID))


	c.JSON(http.StatusOK, gin.H{"message": "note deleted"})
}