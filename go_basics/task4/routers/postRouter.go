package routers

import (
	"task4/controllers"
	"task4/middlewares"

	"github.com/gin-gonic/gin"
)

func PostRouterInit(r *gin.Engine) {

	postRouter := r.Group("/post", middlewares.JWTAuthMiddleware())
	{
		postRouter.POST("/create", controllers.PostController{}.Create)

		postRouter.POST("/edit", controllers.PostController{}.Edit)

		postRouter.POST("/delete", controllers.PostController{}.Delete)

		postRouter.POST("/queryList", controllers.PostController{}.QueryList)

		postRouter.POST("/queryPostInfo", controllers.PostController{}.QueryPostInfo)

	}

}
