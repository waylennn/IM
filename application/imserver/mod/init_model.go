package mod

import (
	"github.com/xormplus/xorm"
	_ "github.com/go-sql-driver/mysql"
)
type (
	OfflineMsg struct {
		Id int
		FromToken string `xorm:"varchar(60) notnull 'sender'"`
		ToToken string `xorm:"varchar(60) notnull 'receiver'"`
		Body  string `xorm:"varchar(1000) 'msg'"`
		TimeStamp int64 `xorm:"BigInt 'time'"`
	}

	SendMsgRequest struct {
		FromToken     string `json:"fromToken"`
		ToToken       string `json:"toToken"`
		Body          string `json:"body"`
		TimeStamp     int64  `json:"timeStamp"`
		RemoteAddress string `json:"remoteAddress"`
	}
)

var engine *xorm.Engine

func InitModEngine(dbName,dbAddress string) (err error){
	engine, err = xorm.NewEngine(dbName, dbAddress)
	if err != nil {
		return
	}

	return nil
}
