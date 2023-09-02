package dao

import (
	"fmt"
	"mall/model"
)

func migration() {
	err := _db.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(
			&model.User{},
			&model.Product{},
			&model.Category{},
			&model.Address{},
			&model.Favorite{},
			&model.Notice{},
			&model.Order{},
			&model.ProductImg{},
			&model.Cart{},
			&model.Admin{},
			&model.Carousel{})
	if err != nil {
		fmt.Println("err", err)
	}
	return
}
