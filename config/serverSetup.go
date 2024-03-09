package config

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"open-note-ne-go/database"
	"open-note-ne-go/domain/notes"
)

func SetupServer(migrationsPath string) (*gin.Engine, *sql.DB) {
	db, err := database.InitialiseDatabase(migrationsPath)
	if err != nil {
		panic(err)
	}

	app := gin.New()
	app.Use(gin.Recovery())
	app.Use(corsMiddleware())

	notesService := notes.Service{Db: db}
	notesService.RegisterRoutes(app)

	return app, db
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
