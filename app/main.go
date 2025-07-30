package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddHandler(c *gin.Context) {
	if values, numsExist := c.GetQueryArray("num"); numsExist {
		total := 0
		for _, val := range values {
			num, err := strconv.Atoi(val)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid number: " + val})
				return
			}
			total += num
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "sum": total})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "missing query param 'num'"})
	}
}

func main() {
	r := gin.Default()
	r.GET("/add", AddHandler)
	r.Run(":3000")
}
