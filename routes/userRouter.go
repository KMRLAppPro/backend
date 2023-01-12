package routes

import (
	controller "github.com/KMRLAppPro/backend/controllers"
	"github.com/KMRLAppPro/backend/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.POST("/upload",controller.UploadFile())
	incomingRoutes.GET("/download",controller.DownloadFile())
}
