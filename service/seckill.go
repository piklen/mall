package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	logging "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"log"
	"mall/cache"
	"mall/dao"
	"mall/model"
	"mall/mq"
	"mall/pkg/e"
	"mall/serializer"
	"math/rand"
	"mime/multipart"
	"strconv"
	"time"
)

type SeckillGoodsImport struct {
}

// SkillGoodsService 限购一个
type SkillGoodsService struct {
	SkillGoodsId uint   `json:"skill_goods_id" form:"skill_goods_id"`
	ProductId    uint   `json:"product_id" form:"product_id"`
	BossId       uint   `json:"boss_id" form:"boss_id"`
	AddressId    uint   `json:"address_id" form:"address_id"`
	Key          string `json:"key" form:"key"`
}
type SeckillGoodsWithMySQL struct {
	ProductId uint   `json:"product_id" form:"product_id"`
	AddressId uint   `json:"address_id" form:"address_id"`
	Num       int    `json:"num" form:"num"`
	BossId    uint   `json:"boss_id" form:"boss_id"`
	Key       string `json:"key" form:"key"`
}
type SeckillGoods struct {
	ProductId uint   `json:"product_id" form:"product_id"`
	AddressId uint   `json:"address_id" form:"address_id"`
	Num       int    `json:"num" form:"num"`
	BossId    uint   `json:"boss_id" form:"boss_id"`
	Key       string `json:"key" form:"key"`
}

// 导入秒杀商品文件
func (service *SeckillGoodsImport) Import(ctx context.Context, file multipart.File) serializer.Response {
	xlFile, err := excelize.OpenReader(file)
	if err != nil {
		logging.Info(err)
	}
	code := e.Success
	rows := xlFile.GetRows("Sheet1")
	length := len(rows[1:])
	skillGoods := make([]*model.SeckillGoods, length, length)
	for index, colCell := range rows {
		if index == 0 {
			continue
		}
		pId, _ := strconv.Atoi(colCell[0])
		bId, _ := strconv.Atoi(colCell[1])
		num, _ := strconv.Atoi(colCell[3])
		money, _ := strconv.ParseFloat(colCell[4], 64)
		skillGood := &model.SeckillGoods{
			ProductId: uint(pId),
			BossId:    uint(bId),
			Title:     colCell[2],
			Money:     money,
			Num:       num,
		}
		skillGoods[index-1] = skillGood
	}
	err = dao.NewSeckillGoodsDao(ctx).CreateByList(skillGoods)
	if err != nil {
		code = e.ErrorUploadFile
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data:   "上传失败",
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

// 直接放到这里，初始化秒杀商品信息，将mysql的信息存入redis中
func (service *SkillGoodsService) InitSkillGoods(ctx context.Context) error {
	skillGoods, _ := dao.NewSeckillGoodsDao(ctx).ListSkillGoods()
	r := cache.RedisClient
	// 加载到redis
	for i := range skillGoods {
		fmt.Println(*skillGoods[i])
		r.Client.HSet("SK"+strconv.Itoa(int(skillGoods[i].Id)), "num", skillGoods[i].Num)
		r.Client.HSet("SK"+strconv.Itoa(int(skillGoods[i].Id)), "money", skillGoods[i].Money)
	}
	return nil
}
func (service *SkillGoodsService) SkillGoods(ctx context.Context, uId uint) serializer.Response {
	mo, _ := cache.RedisClient.Client.HGet("SK"+strconv.Itoa(int(service.SkillGoodsId)), "money").Float64()
	sk := &model.SeckillGood2MQ{
		ProductId:   service.ProductId,
		BossId:      service.BossId,
		UserId:      uId,
		AddressId:   service.AddressId,
		Key:         service.Key,
		Money:       mo,
		SkillGoodId: service.SkillGoodsId,
	}
	err := RedissonSecKillGoods(sk)
	if err != nil {
		return serializer.Response{}
	}
	return serializer.Response{}
}
func (service *SeckillGoods) SkillGoodsWithRedis(ctx context.Context, uId uint) serializer.Response {
	code := e.Success
	//先构建redis中的key值，再去redis进减少库存的操作。
	//先加锁
	cache.RedisClient.Mu.Lock()
	defer cache.RedisClient.Mu.Unlock()
	//如果该用户已经进行了秒杀活动，那么就不能再进行秒杀
	res, err := cache.RedisClient.Client.SIsMember("SK"+strconv.Itoa(int(service.ProductId))+"names", uId).Result()
	if res == true {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    "该用户已进行秒杀！！！！",
		}
	}
	productNum, err := cache.RedisClient.Client.HGet("SK"+strconv.Itoa(int(service.ProductId)), "num").Int()
	//先看是否能预扣库存
	if productNum < service.Num {
		logging.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    "商品库存数量不足,秒杀商品数量预扣失败！！！",
		}
	}
	//对商品数量进行变化
	productNum = -service.Num
	cache.RedisClient.Client.HIncrBy("SK"+strconv.Itoa(int(service.ProductId)), "num", int64(productNum))
	cache.RedisClient.Client.SAdd("SK"+strconv.Itoa(int(service.ProductId))+"names", uId)
	return serializer.Response{
		Status: code,
		Msg:    "秒杀商品数量预扣成功！！！",
	}
}
func (service *SeckillGoodsWithMySQL) SkillGoodsWithMySQL(ctx context.Context, uId uint) serializer.Response {
	//先看是否能预扣库存
	code := e.Success
	skillDao := dao.NewSeckillGoodsDao(ctx)
	err := skillDao.CanPreReduceStocks(service.ProductId, service.Num)
	if err != nil {
		logging.Info(err)
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    "秒杀商品数量预扣失败！！！",
		}
	}
	//生成订单信息

	//通过商品id查询商品信息
	//先不写
	return serializer.Response{
		Status: code,
		Msg:    "秒杀商品数量预扣成功！！！",
	}
}

func RedissonSecKillGoods(sk *model.SeckillGood2MQ) error {
	p := strconv.Itoa(int(sk.ProductId))
	uuid := getUuid(p)
	_, err := cache.RedisClient.Client.Del(p).Result()
	lockSuccess, err := cache.RedisClient.Client.SetNX(p, uuid, time.Second*3).Result()
	if err != nil || !lockSuccess {
		fmt.Println("get lock fail", err)
		return errors.New("get lock fail")
	} else {
		fmt.Println("get lock success")
	}
	//将商品信息传到消息队列中
	_ = SendSecKillGoodsToMQ(sk)
	value, _ := cache.RedisClient.Client.Get(p).Result()
	if value == uuid { // compare value,if equal then del
		_, err := cache.RedisClient.Client.Del(p).Result()
		if err != nil {
			fmt.Println("unlock fail")
			return nil
		} else {
			fmt.Println("unlock success")
		}
	}
	return nil
}
func SendSecKillGoodsToMQ(sk *model.SeckillGood2MQ) error {
	ch, err := mq.RabbitMQ.Channel()
	if err != nil {
		err = errors.New("rabbitMQ err:" + err.Error())
		return err
	}
	q, err := ch.QueueDeclare("skill_goods", true, false, false, false, nil)
	if err != nil {
		err = errors.New("rabbitMQ err:" + err.Error())
		return err
	}
	body, _ := json.Marshal(sk)
	err = ch.Publish("", q.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         body,
	})
	if err != nil {
		err = errors.New("rabbitMQ err:" + err.Error())
		return err
	}
	log.Printf("Sent %s", body)
	return nil
}
func getUuid(gid string) string {
	codeLen := 8
	// 1. 定义原始字符串
	rawStr := "jkwangagDGFHGSERKILMJHSNOPQR546413890_"
	// 2. 定义一个buf，并且将buf交给bytes往buf中写数据
	buf := make([]byte, 0, codeLen)
	b := bytes.NewBuffer(buf)
	// 随机从中获取
	rand.Seed(time.Now().UnixNano())
	for rawStrLen := len(rawStr); codeLen > 0; codeLen-- {
		randNum := rand.Intn(rawStrLen)
		b.WriteByte(rawStr[randNum])
	}
	return b.String() + gid
}
