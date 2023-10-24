package service

import (
	"context"
	logging "github.com/sirupsen/logrus"
	"mall/dao"
	"mall/model"
	"mall/pkg/e"
	"mall/serializer"
)

type FavoritesService struct {
	ProductId  uint `form:"product_id" json:"product_id"`
	BossId     uint `form:"boss_id" json:"boss_id"`
	FavoriteId uint `form:"favorite_id" json:"favorite_id"`
	PageNum    int  `form:"pageNum"`
	PageSize   int  `form:"pageSize"`
}

func (service *FavoritesService) Create(ctx context.Context, uId uint) serializer.Response {
	code := e.Success
	favoriteDao := dao.NewFavoritesDao(ctx)
	//判断是否收藏夹是否已经存在此商品
	exist, _ := favoriteDao.FavoriteExistOrNot(service.FavoriteId, uId)
	if exist {
		code = e.ErrorExistFavorite
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	//创建userDao
	userDao := dao.NewUserDao(ctx)
	user, err := userDao.GetUserById(uId)
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	//创建BossDao
	//bossDao需要使用与userDao相同的数据库连接。
	// 由于NewUserDaoByDB接受一个*gorm.DB参数，可以通过将userDao.DB传递给NewUserDaoByDB来创建一个新的bossDao对象。
	//这样做是为了确保bossDao使用的是与userDao相同的数据库连接，以便在执行查询操作时保持一致性。
	bossDao := dao.NewUserDaoByDB(userDao.DB)
	boss, err := bossDao.GetUserById(service.BossId)
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	//创建ProductDao
	productDao := dao.NewProductDao(ctx)
	product, err := productDao.GetProductById(service.ProductId)
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	//创建FavoriteDao
	favorite := &model.Favorite{
		UserId:    uId,
		User:      *user,
		Product:   *product,
		ProductId: service.ProductId,
		Boss:      *boss,
		BossId:    service.BossId,
	}
	//进行数据插入
	favoriteDao = dao.NewFavoritesByDB(favoriteDao.DB)
	err = favoriteDao.CreateFavorite(favorite)
	if err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}
func (service *FavoritesService) Show(ctx context.Context, uId uint) serializer.Response {
	code := e.Success
	favoritesDao := dao.NewFavoritesDao(ctx)
	if service.PageSize == 0 {
		service.PageSize = 15
	}
	favorites, total, err := favoritesDao.ListFavoriteByUserId(uId, service.PageSize, service.PageNum)
	if err != nil {
		logging.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.BuildListResponse(serializer.BuildFavorites(ctx, favorites), uint(total))
}

// Delete 删除收藏夹
func (service *FavoritesService) Delete(ctx context.Context) serializer.Response {
	code := e.Success

	favoriteDao := dao.NewFavoritesDao(ctx)
	err := favoriteDao.DeleteFavoriteById(service.FavoriteId)
	if err != nil {
		logging.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Data:   e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	return serializer.Response{
		Status: code,
		Data:   e.GetMsg(code),
	}
}
