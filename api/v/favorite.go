package v

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mall/pkg/log"
	"mall/pkg/util"
	"mall/service"
	"net/http"
)

func CreateFavorite(c *gin.Context) {
	service := service.FavoritesService{}
	claim, _ := util.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&service); err == nil {
		res := service.Create(c.Request.Context(), claim.ID)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		log.LogrusObj.Infoln(err)
	}
}
func ShowFavorites(c *gin.Context) {
	service := service.FavoritesService{}
	fmt.Println("---------标记--------")
	claim, _ := util.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&service); err == nil {
		res := service.Show(c.Request.Context(), claim.ID)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		log.LogrusObj.Infoln(err)
	}
}
func DeleteFavorite(c *gin.Context) {
	service := service.FavoritesService{}
	if err := c.ShouldBind(&service); err == nil {
		res := service.Delete(c.Request.Context())
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		log.LogrusObj.Infoln(err)
	}
}
