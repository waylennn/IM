package users_model

import (
	"time"
)

type (
	Users struct {
		Id  int64  `json:"-" xorm:"int(11) pk 'id'"`
		Token string `json:"token" xorm:"varchar(11) notnull 'token'"`
		UserName string `json:"token" xorm:"varchar(30) notnull 'username'"`
		PassWord string `json:"token" xorm:"varchar(40) notnull 'password'"`
		CreateTime time.Time `json:"token" xorm:"DateTime notnull created 'create_time'"`
		UpdateTime time.Time `json:"token" xorm:"DateTime notnull updated 'update_time'"`
	}
)

