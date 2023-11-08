package v

import (
	"github.com/gin-gonic/gin"
	"mall/pkg/log"
	"mall/pkg/util"
	"mall/service"
	"net/http"
)

func ImportSkillGoods(c *gin.Context) {
	var skillGoodsImport service.SeckillGoodsImport
	file, _, _ := c.Request.FormFile("file")
	if err := c.ShouldBind(&skillGoodsImport); err == nil {
		res := skillGoodsImport.Import(c.Request.Context(), file)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		log.LogrusObj.Infoln(err, "ImportSkillGoods")
	}
}
func InitSkillGoods(c *gin.Context) {
	var skillGoods service.SkillGoodsService
	if err := c.ShouldBind(&skillGoods); err == nil {
		res := skillGoods.InitSkillGoods(c.Request.Context())
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		log.LogrusObj.Infoln(err, "InitSkillGoods")
	}
}
func SkillGoods(c *gin.Context) {
	var skillGoods service.SkillGoodsService
	claim, _ := util.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&skillGoods); err == nil {
		res := skillGoods.SkillGoods(c.Request.Context(), claim.ID)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
		log.LogrusObj.Infoln(err, "SkillGoods")
	}
}
