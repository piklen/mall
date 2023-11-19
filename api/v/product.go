package v

import (
	"github.com/gin-gonic/gin"
	"mall/pkg/log"
	"mall/pkg/util"
	"mall/service"
	"net/http"
)

// 创建商品
func CreateProduct(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["file"]
	claim, _ := util.ParseToken(c.GetHeader("Authorization"))
	createProductService := service.ProductService{}
	if err := c.ShouldBind(&createProductService); err == nil {
		res := createProductService.Create(c.Request.Context(), claim.ID, files)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		log.LogrusObj.Infoln(err)
	}
}

// ListProducts 商品列表
func ListProducts(c *gin.Context) {
	listProductsService := service.ProductService{}
	if err := c.ShouldBind(&listProductsService); err == nil {
		res := listProductsService.List(c.Request.Context())
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		log.LogrusObj.Infoln(err)
	}
}

// SearchProducts 搜索商品
func SearchProducts(c *gin.Context) {
	searchProductsService := service.ProductService{}
	if err := c.ShouldBind(&searchProductsService); err == nil {
		res := searchProductsService.Search(c.Request.Context())
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		log.LogrusObj.Infoln(err)
	}
}

func ShowProduct(c *gin.Context) {
	showProductService := service.ProductService{}
	if err := c.ShouldBind(&showProductService); err == nil {
		res := showProductService.Show(c.Request.Context(), c.Query("id"))
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		log.LogrusObj.Infoln(err)
	}
}
