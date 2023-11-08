package dao

import (
	"context"
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
