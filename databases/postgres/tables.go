package database

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTables(c *gin.Context, request ConnectionRequest) Tables {
	row, err := PostgreConnect.Query("SELECT table_name FROM information_schema.tables WHERE table_schema=$1", request.Schema)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tables: " + err.Error()})
		return Tables{}
	}

	var tables Tables
	tables.TablesMap = make(map[string][]ColumnInfo)

	for row.Next() {
		var tableName string
		err := row.Scan(&tableName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan table name: " + err.Error()})
		}

		query := `SELECT
					attname AS column_name,
					format_type(atttypid, atttypmod) AS data_type,
					EXISTS (
						SELECT 1 FROM pg_constraint
						WHERE conrelid = attrelid AND attnum = ANY(conkey) AND contype = 'p'
					) AS is_pk
				FROM pg_attribute
				WHERE attrelid = $1 ::regclass
				AND attnum > 0
				AND NOT attisdropped
				ORDER BY attnum;`

		row, err := PostgreConnect.Query(query, fmt.Sprintf("%v.%v", request.Schema, tableName))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get column information: " + err.Error()})
			return tables
		}

		if row.Err() != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading column information: " + row.Err().Error()})
		}

		for row.Next() {
			var schema ColumnInfo
			err := row.Scan(&schema.Name, &schema.Type, &schema.PK)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan column data: " + err.Error()})
				return tables
			}

			tables.TablesMap[tableName] = append(tables.TablesMap[tableName], schema)
		}
	}

	if row.Err() != nil {
		log.Println(err)
	}

	return tables
}

func GetPK(c *gin.Context) {
	query := `SELECT
				c.conname AS fk_name,
				src.relname AS from_table,
				src_att.attname AS from_column,
				dst.relname AS to_table,
				dst_att.attname AS to_column
			FROM pg_constraint c
			JOIN pg_class src ON src.oid = c.conrelid
			JOIN pg_namespace src_ns ON src_ns.oid = src.relnamespace
			JOIN pg_class dst ON dst.oid = c.confrelid
			CROSS JOIN LATERAL generate_subscripts(c.conkey, 1) s(i)
			JOIN pg_attribute src_att ON src_att.attrelid = c.conrelid AND src_att.attnum = c.conkey[s.i]
			JOIN pg_attribute dst_att ON dst_att.attrelid = c.confrelid AND dst_att.attnum = c.confkey[s.i]
			WHERE c.contype = 'f'
			AND src_ns.nspname = $1
			ORDER BY src.relname, c.conname;`

	rows, err := PostgreConnect.Query(query, ConnInfo.Schema)
	if err != nil {
		log.Println("failed to query keys:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var relations []Relation

	for rows.Next() {
		var k Key

		err := rows.Scan(&k.FkName, &k.FromTable, &k.FromColumn, &k.ToTable, &k.ToColumn)
		if err != nil {
			log.Println("scan error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		relations = append(relations, Relation{
			ID:   k.FkName,
			From: fmt.Sprintf("%s.%s", k.FromTable, k.FromColumn),
			To:   fmt.Sprintf("%s.%s", k.ToTable, k.ToColumn),
		})
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"connected": true,
		"Key":       relations,
	})
}
