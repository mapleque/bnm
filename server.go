package main

import (
	"os"

	"github.com/mapleque/bnm/server"
)

func main() {
	db := server.NewDB(os.Getenv("DB_DSN"))
	wx := server.NewWxSDK(os.Getenv("WX_APPID"), os.Getenv("WX_SECRET"))
	s := server.NewServer(
		db,
		wx,
		os.Getenv("HOST"),
		os.Getenv("QINIU_ACCESS_KEY"),
		os.Getenv("QINIU_SECRET_KEY"),
	)
	s.Run()
}
