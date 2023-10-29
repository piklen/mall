package service

import (
	"context"
	"errors"
	"fmt"
	logging "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"mall/dao"
	"mall/model"
	"mall/pkg/e"
	"mall/pkg/util"
	"mall/serializer"
	"strconv"
)

type OrderPay struct {
	OrderId   uint    `form:"order_id" json:"order_id"`
	Money     float64 `form:"money" json:"money"`
	OrderNo   string  `form:"orderNo" json:"orderNo"`
	ProductID int     `form:"product_id" json:"product_id"`
	PayTime   string  `form:"payTime" json:"payTime" `
	Sign      string  `form:"sign" json:"sign" `
	BossID    int     `form:"boss_id" json:"boss_id"`
	BossName  string  `form:"boss_name" json:"boss_name"`
	Num       int     `form:"num" json:"num"`
	Key       string  `form:"key" json:"key"`
}

func (o *OrderPay) PayDown(ctx context.Context, uId uint) serializer.Response {
	code := e.Success
	//创建一个订单数据访问对象，并使用Transaction方法开始一个数据库事务。ctx是传入的上下文对象。
	err := dao.NewOrderDao(ctx).Transaction(func(tx *gorm.DB) error {
		//设置加密密钥
		util.Encrypt.SetKey(o.Key)
		orderDao := dao.NewOrderByDB(tx)
		//获取订单总金额，
		order, err := orderDao.GetOrderById(o.OrderId)
		if err != nil {
			logging.Info(err)
			return err
		}
		money := order.Money
		num := order.Num
		money = money * float64(num)

		userDao := dao.NewUserDaoByDB(tx)
		user, err := userDao.GetUserById(uId)
		if err != nil {
			logging.Info(err)
			code = e.ErrorDatabase
			return err
		}

		// 对钱进行解密。减去订单。再进行加密。
		moneyStr := util.Encrypt.AesDecoding(user.Money)
		moneyFloat, _ := strconv.ParseFloat(moneyStr, 64)
		if moneyFloat-money < 0.0 { // 金额不足进行回滚
			logging.Info(err)
			code = e.ErrorDatabase
			return errors.New("金币不足")
		}

		finMoney := fmt.Sprintf("%f", moneyFloat-money)
		user.Money = util.Encrypt.AesEncoding(finMoney)
		//更新用户金额
		err = userDao.UpdateUserById(uId, user)
		if err != nil { // 更新用户金额失败，回滚
			logging.Info(err)
			code = e.ErrorDatabase
			return err
		}

		//更新boss金币数量
		boss := new(model.User)
		boss, err = userDao.GetUserById(uint(o.BossID))
		moneyStr = util.Encrypt.AesDecoding(boss.Money)
		moneyFloat, _ = strconv.ParseFloat(moneyStr, 64)
		finMoney = fmt.Sprintf("%f", moneyFloat+money)
		boss.Money = util.Encrypt.AesEncoding(finMoney)

		err = userDao.UpdateUserById(uint(o.BossID), boss)
		if err != nil { // 更新boss金额失败，回滚
			logging.Info(err)
			code = e.ErrorDatabase
			return err
		}
		//更新商品数量
		product := new(model.Product)
		productDao := dao.NewProductDaoByDB(tx)
		product, err = productDao.GetProductById(uint(o.ProductID))
		if err != nil {
			return err
		}
		product.Num -= num
		err = productDao.UpdateProduct(uint(o.ProductID), product)
		if err != nil { // 更新商品数量减少失败，回滚
			logging.Info(err)
			code = e.ErrorDatabase
			return err
		}

		// 更新订单状态
		order.Type = 2
		err = orderDao.UpdateOrderById(o.OrderId, order)
		if err != nil { // 更新订单失败，回滚
			logging.Info(err)
			code = e.ErrorDatabase
			return err
		}

		productUser := model.Product{
			Name:          product.Name,
			CategoryID:    product.CategoryID,
			Title:         product.Title,
			Info:          product.Info,
			ImgPath:       product.ImgPath,
			Price:         product.Price,
			DiscountPrice: product.DiscountPrice,
			Num:           num,
			OnSale:        false,
			BossID:        uId,
			BossName:      user.UserName,
			BossAvatar:    user.Avatar,
		}

		err = productDao.CreateProduct(&productUser)
		if err != nil { // 买完商品后创建成了自己的商品失败。订单失败，回滚
			logging.Info(err)
			code = e.ErrorDatabase
			return err
		}

		return nil

	})

	if err != nil {
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}
