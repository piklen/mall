package service

import (
	"context"
	"mall/dao"
	"mall/pkg/e"
	"mall/pkg/util"
	"mall/serializer"
)

type ListCarouselsService struct {
}

func (service *ListCarouselsService) List() serializer.Response {
	code := e.Success
	carouselsCtx := dao.NewCarouselDao(context.Background())
	carousels, err := carouselsCtx.ListAddress()
	if err != nil {
		util.LogrusObj.Infoln("err", err)
		code = e.Error
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.BuildListResponse(serializer.BuildCarousels(carousels), uint(len(carousels)))
}
