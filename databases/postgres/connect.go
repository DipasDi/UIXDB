package database

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func InitDB(c *gin.Context) {
	var err error
	var request ConnectionRequest
	err = c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ConnInfo = request

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		request.Host,
		request.Port,
		request.User,
		request.Password,
		request.Database,
	)

	PostgreConnect, err = sql.Open("postgres", dsn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database driver failed to open connection"})
		return
	}

	if err := PostgreConnect.PingContext(c.Request.Context()); err != nil {
		PostgreConnect.Close()
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	tables := GetTables(c, request)

	c.JSON(http.StatusOK, gin.H{
		"connected": true,
		"database":  request.Database,
		"host":      request.Host,
		"port":      request.Port,
		"table":     tables,
	})
}
