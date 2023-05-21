package model

import (
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	_ "github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("滴滴滴滴"))
var sessionName = "session-name"

func SetSession(c *gin.Context, name string, id int64) error {
	session, _ := store.Get(c.Request, sessionName)
	session.Values["name"] = name
	session.Values["id"] = id
	return session.Save(c.Request, c.Writer)
}
func FlushSession(c *gin.Context) error {
	session, _ := store.Get(c.Request, sessionName)
	session.Values["name"] = ""
	session.Values["id"] = 0
	return session.Save(c.Request, c.Writer)
}
func GetSession(c *gin.Context) map[interface{}]interface{} {
	session, _ := store.Get(c.Request, sessionName)
	return session.Values
}
