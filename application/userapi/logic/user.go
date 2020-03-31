package logic

import (
	users_model "awesomeProject/application/userapi/mod"
	"awesomeProject/application/userapi/utils"
	"fmt"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	username := c.DefaultPostForm("username","")
	password := c.DefaultPostForm("password","")
	if password == "" || username == "" {
		utils.ResponseError(c,1001)
		return
	}

	ok,err := users_model.UserLogin(username,password)
	if !ok {
		utils.ResponseError(c,1005)
		return
	}
	if err != nil {
		utils.ResponseError(c,1001)
		return
	}

	token,err := utils.GetToken()
	if err != nil {
		utils.ResponseError(c,1011)
		return
	}

	utils.ResponseSuccess(c,token)
	return
}


func Register(c *gin.Context){
	username := c.DefaultPostForm("username","")
	password := c.DefaultPostForm("password","")
	if password == "" || username == "" {
		utils.ResponseError(c,1001)
		return
	}
	//用户是否已经存在
	ok, err := users_model.UserExist(username)
	if ok {
		utils.ResponseError(c,1002)
		return
	}
	if err != nil {
		fmt.Println(err)
		utils.ResponseError(c,1001)
		return
	}
	//用户注册
	err = users_model.Register(username,password)
	if err != nil {
		fmt.Println(err)
		utils.ResponseError(c,1012)
		return
	}
	utils.ResponseSuccess(c,"插入成功")
	return
}