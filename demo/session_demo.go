package demo

import (
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SessionTestRoute(r *gin.Engine)  {
	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.GET("/incr", func(c *gin.Context) {
		session := sessions.Default(c)
		var count int
		v := session.Get("count")
		if v == nil {
			count = 0
		} else {
			count = v.(int)
			count++
		}
		session.Set("count", count)
		session.Save()
		c.JSON(200, gin.H{"count": count})
	})
}
