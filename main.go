package main

import (
	"plpl/routers/api"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	r.POST("api/run", api.PostPLCode)

	r.Run(":8080")
}
