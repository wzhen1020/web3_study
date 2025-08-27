<<<<<<< HEAD
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
=======
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
>>>>>>> 0a25e3ccee0f1071177d78e3d62c3db86e3c70b4
