package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"

	"github.com/mapleque/bnm/server"
)

func main() {
	db := server.NewDB(os.Getenv("DB_DSN"))
	initDBSchema(db, "./sql/")
	initBusiness(db)
	initCustomerAddress(db)
	initItems(db)
	initOrders(db)
}

func initBusiness(db server.DBConn) {
	if err := db.Execute("INSERT INTO business_profile "+
		"(buid, name, avatar, `desc`, qrcode, wxid, status) "+
		"VALUES (?,?,?,?,?,?,?)",
		1,
		"BangNiMai官方旗舰店",
		"",
		"您正在使用的这个应用还可以购买一些扩展包，如有需要请戳。",
		"",
		"no wxid",
		1,
	).Error(); err != nil {
		panic(err)
	}
	for i := 2; i < 10; i++ {
		if err := db.Execute("INSERT INTO business_profile "+
			"(buid, name, avatar, `desc`, qrcode, wxid, status) "+
			"VALUES (?,?,?,?,?,?,?)",
			i,
			fmt.Sprintf("无名之店-%d", i),
			"",
			fmt.Sprintf("店虽无名，作用却很大%d", i),
			"",
			"no wxid",
			1,
		).Error(); err != nil {
			panic(err)
		}
	}
}

func initCustomerAddress(db server.DBConn) {
	for i := 1; i < 15; i++ {
		if err := db.Execute("INSERT INTO customer_address "+
			"(cuid, label, reciever, address, phone) "+
			"VALUES (?,?,?,?,?)",
			1,
			fmt.Sprintf("第%d个地址", i),
			fmt.Sprintf("收货人%d", i),
			fmt.Sprintf("地址%d", i),
			fmt.Sprintf("电话%d", i),
		).Error(); err != nil {
			panic(err)
		}
	}
}

func initItems(db server.DBConn) {
	for i := 1; i < 21; i++ {
		if err := db.Execute("INSERT INTO item "+
			"(bid, name, price, pic, `desc`, status) "+
			"VALUES (?,?,?,?,?,1)",
			1,
			fmt.Sprintf("神秘商品-%d", i),
			rand.Intn(1000),
			"",
			fmt.Sprintf("神秘商品必有神秘之处，如果描述的足够长，它将变得更加神秘"),
		).Error(); err != nil {
			panic(err)
		}
	}
}

func initOrders(db server.DBConn) {
	status := []int{0, 0, 0, 0, 1, 1, 1, 1}
	stage := []string{
		server.ORDER_STAGE_NEW,
		server.ORDER_STAGE_PAID,
		server.ORDER_STAGE_TRANSPORT,
		server.ORDER_STAGE_C_W_CANCEL,
		server.ORDER_STAGE_FINISH,
		server.ORDER_STAGE_C_CANCEL,
		server.ORDER_STAGE_B_CANCEL,
		server.ORDER_STAGE_REPAID,
	}
	for i := 1; i < 21; i++ {
		if err := db.Execute("INSERT INTO `order` "+
			"(cuid, bid, itid, name, price, counts, reciever, address, phone, status, stage, additional) "+
			"VALUES (?,?,?,?,?,?,?,?,?,?,?,?)",
			1,
			1,
			i,
			fmt.Sprintf("神秘商品-%d", i),
			rand.Intn(1000),
			i,
			fmt.Sprintf("收货人%d", i),
			fmt.Sprintf("地址也要足够长才显得真实%d", i),
			fmt.Sprintf("1380013800%d", i),
			status[i%8],
			stage[i%8],
			fmt.Sprintf("购物备注也要跟得上时代，随便说什么是没有用的（%d）", i),
		).Error(); err != nil {
			panic(err)
		}
	}
}

func initDBSchema(db server.DBConn, schemaPath string) {
	sqlDir, err := ioutil.ReadDir(schemaPath)
	if err != nil {
		panic(err)
	}
	for _, fi := range sqlDir {
		if fi.IsDir() {
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), ".SQL") {
			sqlfile := schemaPath + string(os.PathSeparator) + fi.Name()
			sourceSql, err := ioutil.ReadFile(sqlfile)
			if err != nil {
				panic(err)
			}
			for _, sql := range strings.Split(string(sourceSql), ";") {
				if len(strings.TrimSpace(sql)) > 0 {
					if err := db.Execute(sql).Error(); err != nil {
						panic(err)
					}
				}
			}
		}
	}
}
