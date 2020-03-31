package users_model

import (
	"crypto/md5"
	"fmt"
	"github.com/satori/go.uuid"
)

func UserLogin(username,password string)(ok bool,err error){
	password = fmt.Sprintf("%x", md5.Sum([]byte(password)))
	user := &Users{}
	ok, err = engine.Where("username = ? and password = ?", username, password).Get(user)
	if err != nil {
		return
	}
	return
}

//验证是否用户已经存在
func UserExist(username string)(ok bool,err error){
	user := &Users{}
	ok, err = engine.Where("username = ?", username).Get(user)
	if err != nil {
		return
	}
	return
}

func Register(username,password string)(err error){
	password = fmt.Sprintf("%x", md5.Sum([]byte(password)))
	user := &Users{
		Token:uuid.NewV4().String(),
		UserName :username ,
		PassWord : password,
	}
	_, err = engine.Insert(user)
	if err != nil {
		return
	}
	return
}
