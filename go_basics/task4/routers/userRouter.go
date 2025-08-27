package routers

import (
	"task4/handlers"

	"github.com/gin-gonic/gin"
)

func UserRouterInit(r *gin.Engine) {
	// r.Group("/user",middlewares.)

	UserRouter := r.Group("/user")
	{
		UserRouter.POST("/login", handlers.Login)
	}

}
