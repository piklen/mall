package v

import (
	"github.com/gin-gonic/gin"
	"mall/cache"
	"mall/service"
	"net/http"
	"strconv"
)

type Product struct {
	ID       int    `form:"id" json:"id"`
	UserName string `form:"user_name" json:"user_name"`
	Password string `form:"password" json:"password"`
	Key      string `form:"key" json:"key"` // 前端进行验证
}

func ListCarousels(c *gin.Context) {
	listCarouselsService := service.ListCarouselsService{}
	if err := c.ShouldBind(&listCarouselsService); err == nil {
		res := listCarouselsService.List()
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
	}
}

// View 获取点击数
func (product *Product) View() uint64 {
	countStr, _ := cache.RedisClient.Client.Get(cache.ProductViewKey(uint(product.ID))).Result()
	count, _ := strconv.ParseUint(countStr, 10, 64)
	return count
}

// AddView 商品游览
func (product *Product) AddView() {
	// 增加视频点击数
	cache.RedisClient.Client.Incr(cache.ProductViewKey(uint(product.ID)))
	// 增加排行点击数
	cache.RedisClient.Client.ZIncrBy(cache.RankKey, 1, strconv.Itoa(int(product.ID)))
}
func GetProductView(c *gin.Context) {
	var product = Product{}
	if err := c.ShouldBind(&product); err == nil {
		id, err := strconv.Atoi(strconv.Itoa(product.ID))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
			return
		}
		product := &Product{ID: int(id)}
		count := product.View()
		c.JSON(http.StatusOK, gin.H{"product_id": product.ID, "view_count": count})
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
	}

}
func AddProductView(c *gin.Context) {
	var product = Product{}
	if err := c.ShouldBind(&product); err == nil {
		id, err := strconv.Atoi(strconv.Itoa(product.ID))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
			return
		}
		product := &Product{ID: int(id)}
		product.AddView()
		c.JSON(http.StatusOK, gin.H{"message": "浏览量添加成功！！！"})
	} else {
		c.JSON(http.StatusBadRequest, ErrorResponse(err))
	}
}
