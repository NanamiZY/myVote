package model

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var GlobalConn *gorm.DB

func New() {
	my := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", "root", "886dsbzy", "127.0.0.1:3306", "myvote")
	conn, err := gorm.Open(mysql.Open(my), &gorm.Config{})
	if err != nil {
		fmt.Printf("err:%s\n", err)
	}
	GlobalConn = conn
}
func Close() {
	db, _ := GlobalConn.DB()
	_ = db.Close()
}
