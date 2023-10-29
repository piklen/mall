package v

import (
	"github.com/gin-gonic/gin"
	"mall/pkg/log"
	"mall/service"
	"net/http"
)

func ListProductImg(c *gin.Context) {
	var listProductImgService service.ListProductImgService
	if err := c.ShouldBind(&listProductImgService); err == nil {
		res := listProductImgService.List(c.Request.Context(), c.Query("id"))
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		log.LogrusObj.Infoln(err)
	}
}
