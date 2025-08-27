package main

import (
	"task4/routers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	routers.UserRouterInit(r)
	routers.CommentTouterInit(r)
	routers.PostRouterInit(r)
	r.Run()
}
