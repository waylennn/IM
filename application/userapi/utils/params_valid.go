package utils

import (
	"github.com/gin-gonic/gin"
	"strconv"
)


func GetQueryInt64(c *gin.Context, key string) (value int64, err error) {

	idstr, ok := c.GetQuery(key)
	if !ok {
		return
	}

	value, err = strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		return
	}

	return
}


