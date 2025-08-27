package routers

import (
	"task4/controllers"
	"task4/middlewares"

	"github.com/gin-gonic/gin"
)

func CommentTouterInit(r *gin.Engine) {

	commentRouter := r.Group("/comment,", middlewares.JWTAuthMiddleware())
	{
		commentRouter.POST("/add", controllers.CommentController{}.Add)

		commentRouter.POST("/queryList", controllers.CommentController{}.QueryList)

	}

}
