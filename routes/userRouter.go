package routes

import (
	controller "github.com/KMRLAppPro/backend/controllers"
	"github.com/KMRLAppPro/backend/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.POST("/upload/:image_name",controller.UploadFile())
	incomingRoutes.GET("/download/:image_name",controller.DownloadFile())
}
