package router

import (
	"go_web/api/controller"

	"github.com/gin-gonic/gin"
)

func ProductRouter(r *gin.Engine) {
	pdr := r.Group("/api/v1/products")
	pc := &controller.ProductController{}
	{
		pdr.GET("/", pc.GetProducts)
		pdr.GET("/:id", pc.GetProduct)
		pdr.POST("/", pc.CreateProduct)
		pdr.PUT("/:id", pc.UpdateProduct)
		pdr.DELETE("/:id", pc.DeleteProduct)
	}
}
