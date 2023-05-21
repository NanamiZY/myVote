package model

import "fmt"

func Check(id, name string) bool {
	user := &User{}
	if GlobalConn.Table("users").Where("id=?", id).First(user).Error != nil {
		return false
	}
	if user.Name != name {
		return false
	}
	return true
}
func GetUser(name, pwd string) *User {
	ret := &User{}
	sql := "select * from `users` where name=? and pwd=? limit 1"
	err := GlobalConn.Raw(sql, name, pwd).Scan(ret).Error
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
	}
	return ret
}
