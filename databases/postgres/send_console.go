package database

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func SendConsole(c *gin.Context) {
	if PostgreConnect == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Database connection is not initialized",
		})
		return
	}

	var req struct {
		Statements []string `json:"statements"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload: " + err.Error(),
		})
		return
	}

	if len(req.Statements) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No SQL statements provided",
		})
		return
	}

	tx, err := PostgreConnect.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to begin transaction: " + err.Error(),
		})
		return
	}
	defer tx.Rollback()

	var allRows [][]string

	for _, stmt := range req.Statements {
		trimmed := strings.TrimSpace(stmt)

		if strings.HasPrefix(strings.ToUpper(trimmed), "SELECT") {
			rows, err := tx.Query(stmt)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Query execution failed: " + err.Error(),
				})
				return
			}

			cols, err := rows.Columns()
			if err != nil {
				rows.Close()
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to retrieve columns: " + err.Error(),
				})
				return
			}

			values := make([]any, len(cols))
			scanArgs := make([]any, len(cols))
			for idx := range values {
				scanArgs[idx] = &values[idx]
			}

			allRows = append(allRows, cols)

			for rows.Next() {
				if err := rows.Scan(scanArgs...); err != nil {
					rows.Close()
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": "Failed to scan row: " + err.Error(),
					})
					return
				}

				rowStrings := make([]string, len(cols))
				for idx, val := range values {
					switch v := val.(type) {
					case nil:
						rowStrings[idx] = "NULL"
					case []byte:
						rowStrings[idx] = string(v)
					default:
						rowStrings[idx] = fmt.Sprint(v)
					}
				}

				allRows = append(allRows, rowStrings)
			}

			if err := rows.Err(); err != nil {
				rows.Close()
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Row iteration error: " + err.Error(),
				})
				return
			}
			rows.Close()

		} else {
			if _, err = tx.Exec(stmt); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Execution failed: " + err.Error(),
				})
				return
			}
		}
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to commit transaction: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"rows":    allRows,
	})
}
