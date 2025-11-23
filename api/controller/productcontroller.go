package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Product struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name"`
	Price       uint   `json:"price"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	Category    string `json:"category"`
	Isfeatured  bool   `json:"is_featured"`
	Stock       int    `json:"stock"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type ProductResponse struct {
	Products []Product `json:"products"`
	Total    int       `json:"total"`
	Page     int       `json:"page"`
	Limit    int       `json:"limit"`
}

// ProductController ...
type ProductController struct{}

// GetProducts ...
func (pc *ProductController) GetProducts(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "get products",
	})
}

// GetProduct ...
func (pc *ProductController) GetProduct(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "get product " + id,
	})
}

// CreateProduct ...
func (pc *ProductController) CreateProduct(c *gin.Context) {
	var json Product
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": json})
}

// UpdateProduct ...
func (pc *ProductController) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "update product " + id,
	})
}

// DeleteProduct ...
func (pc *ProductController) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "delete product " + id,
	})
}
