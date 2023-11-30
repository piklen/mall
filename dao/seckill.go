package dao

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"mall/model"
)

type SeckillGoodsDao struct {
	*gorm.DB
}

func NewSeckillGoodsDao(ctx context.Context) *SeckillGoodsDao {
	return &SeckillGoodsDao{NewDBClient(ctx)}
}

func (dao *SeckillGoodsDao) Create(in *model.SeckillGoods) error {
	return dao.Model(&model.SeckillGoods{}).Create(&in).Error
}
func (dao *SeckillGoodsDao) CreateByList(in []*model.SeckillGoods) error {
	return dao.Model(&model.SeckillGoods{}).Create(&in).Error
}

// ListSkillGoods 将MySQL秒杀商品库中的商品数量大于1的进行返回
func (dao *SeckillGoodsDao) ListSkillGoods() (resp []*model.SeckillGoods, err error) {
	err = dao.Model(&model.SeckillGoods{}).Where("num > 0").Find(&resp).Error
	return
}

// 判断能否进行预扣库存

//	func (dao *SeckillGoodsDao) CanPreReduceStocks(productID uint, productNum int) (err error) {
//		err = dao.Model(&model.SeckillGoods{}).
//			Where("product_id = ? AND num - ? >= 0", productID, productNum).
//			Update("num", gorm.Expr("num - ?", productNum)).
//			Error
//		return err
//	}
func (dao *SeckillGoodsDao) CanPreReduceStocks(productID uint, productNum int) (err error) {
	var goods model.SeckillGoods

	// 查询商品记录
	err = dao.Model(&model.SeckillGoods{}).
		Where("product_id = ?", productID).
		First(&goods).
		Error

	if err != nil {
		return err
	}

	// 检查库存是否足够
	if goods.Num-productNum >= 0 {
		// 更新库存
		err = dao.Model(&model.SeckillGoods{}).
			Where("product_id = ? AND num - ? >= 0", productID, productNum).
			Updates(map[string]interface{}{"num": gorm.Expr("num - ?", productNum)}).
			Error
	} else {
		// 库存不足，返回错误
		err = errors.New("库存预修改失败")
	}

	return err
}
