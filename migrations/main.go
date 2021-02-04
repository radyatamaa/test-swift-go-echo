package main

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	model "github.com/models"
)

func main() {
	//local
	//db, err := gorm.Open("mysql", "root:@(localhost:3306)/swift_logistic?charset=utf8&parseTime=True&loc=Local")
	//if err != nil {
	//	fmt.Println(err)
	//}

	db, err := gorm.Open("mysql", "adminbkni@bkni-ri:Standar123.@(bkni-ri.mysql.database.azure.com)/test_db?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
	}

	migration := model.MigrationHistory{}
	errmigration := db.AutoMigrate(&migration)
	if errmigration != nil {
		migration := model.MigrationHistory{
			DescMigration: "Add Table Migration",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}

	user := model.User{}
	erruser := db.AutoMigrate(&user)
	if erruser != nil {
		migration := model.MigrationHistory{
			DescMigration: "Add Table User",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}

	Product := model.Product{}
	errProduct := db.AutoMigrate(&Product)
	if errProduct != nil {
		migration := model.MigrationHistory{
			DescMigration: "Add Table Product",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}


	Inbound := model.Inbound{}
	errInbound := db.AutoMigrate(&Inbound).AddForeignKey("product_id", "products(id)", "RESTRICT", "RESTRICT")
	if errInbound != nil {
		migration := model.MigrationHistory{
			DescMigration: "Add Table Inbound",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}

	Outbound := model.Outbound{}
	errOutbound := db.AutoMigrate(&Outbound).AddForeignKey("product_id", "products(id)", "RESTRICT", "RESTRICT")
	if errOutbound != nil {
		migration := model.MigrationHistory{
			DescMigration: "Add Table Outbound",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}

	Stock := model.Stock{}
	errStock := db.AutoMigrate(&Stock).AddForeignKey("product_id", "products(id)", "RESTRICT", "RESTRICT")
	if errStock != nil {
		migration := model.MigrationHistory{
			DescMigration: "Add Table Stock",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}

	errStock2 := db.Model(&Stock).AddForeignKey("outbound_id", "outbounds(id)", "RESTRICT", "RESTRICT")
	if errStock2 != nil {
		migration := model.MigrationHistory{
			DescMigration: "Add foregn key outbound_id Table Stock",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}

	errStock3 := db.Model(&Stock).AddForeignKey("inbound_id", "inbounds(id)", "RESTRICT", "RESTRICT")
	if errStock3 != nil {
		migration := model.MigrationHistory{
			DescMigration: "Add foregn key outbound_id Table Stock",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}

	order := model.Order{}
	errOrder := db.AutoMigrate(&order)
	if errOrder != nil {
		migration := model.MigrationHistory{
			DescMigration: "Add Table Order",
			Date:          time.Now(),
		}

		db.Create(&migration)
	}

}
