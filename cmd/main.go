package main

import (
	"net/http"
	"note-manager-api/internals/auth"
	"note-manager-api/internals/db"
	"note-manager-api/internals/note"
	"note-manager-api/internals/redis"

	"github.com/gin-gonic/gin"
)

func main() {
	db.DBConnect()
	redis.ConnectRedis()
	r := gin.Default()
	r.POST("/register", auth.Register)
	r.POST("/login", auth.Login)
	r.LoadHTMLGlob("../templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	protected := r.Group("/api")
	protected.Use(auth.JWTMiddleware())
	{
		protected.GET("/notes", note.GetAll)
		protected.POST("/notes", note.CreateNote)
		protected.DELETE("/notes/:id", note.DeleteNote)
		protected.PUT("/notes/:id", note.UpdateNote)
		protected.GET("/favorites", note.GetAllFav)
		protected.POST("/favorites/:note_id", note.CreateFav)
		protected.DELETE("/favorites/:note_id", note.DeleteFavNote)
	}
	r.Run(":8484")
}
