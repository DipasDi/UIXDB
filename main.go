package main

import (
	database "UIXDB/databases/postgres"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Static("/assets", "./assets")
	r.LoadHTMLGlob("templates/*.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.POST("/api/connect", database.InitDB)
	r.GET("/api/getKey", database.GetPK)
	r.POST("/api/query", database.SendConsole)
	r.POST("/api/table", database.GetTable)

	if err := r.Run("0.0.0.0:8082"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
