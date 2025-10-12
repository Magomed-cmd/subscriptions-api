package router

import (
	"net/http"
	"subscriptions-api/internal/handler"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRoutes(r *gin.Engine, subHandler *handler.SubscriptionHandler) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	v1 := r.Group("/api/v1")
	{
		sub := v1.Group("/subscriptions")
		{
			sub.POST("", subHandler.Create)
			sub.GET("/:id", subHandler.GetByID)
			sub.PUT("/:id", subHandler.Update)
			sub.DELETE("/:id", subHandler.Delete)
			sub.GET("", subHandler.List)
			sub.GET("/total", subHandler.TotalCostForPeriod)
		}
	}
}
