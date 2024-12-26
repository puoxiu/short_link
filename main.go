package main

import (
	"flag"
	"fmt"
	"short_link_svc/core"
	urlsmodel "short_link_svc/model/urls_model"
)

type Options struct {
	DB bool
}

func main() {
	var opt Options
	// 解析命令行参数，-db 表示是否初始化数据库，默认为 false
	flag.BoolVar(&opt.DB, "db", false, "database")
	flag.Parse()

	if opt.DB {
		db := core.InitGorm("root:root@tcp(127.0.0.1:13306)/short_link_db?charset=utf8mb4&parseTime=True&loc=Local")
		err := db.AutoMigrate(&urlsmodel.UrlsModel{},
		)
		if err != nil {
			fmt.Println("数据库表生成失败，错误信息：", err)
			return
		}
		fmt.Println("数据库表生成成功")
	}
}