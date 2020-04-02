package logic

import (
	"awesomeProject/application/gateway/main/config"
	gateway_model "awesomeProject/application/gateway/mod"
	im "awesomeProject/application/imserver/protoc"
	user "awesomeProject/application/userserver/protoc"
	"awesomeProject/application/utils"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/errors"
	"math/rand"
	"time"
)

var (
	UserClient  user.UserService
	ImClient    im.ImService
	AddressList []*config.ImRpc
)

func NewAddressList(addressList []*config.ImRpc) {
	AddressList = addressList
}

func Send(c *gin.Context) {

	fromToken := c.DefaultPostForm("fromToken", "")
	toToken := c.DefaultPostForm("toToken", "")
	body := c.DefaultPostForm("body", "")
	timestamp := c.DefaultPostForm("timeStamp", "")
	if fromToken == "" || toToken == "" || timestamp == "" {
		utils.ResponseError(c, utils.ErrCodeParamNotExist, nil)
		return
	}
	//用户rpc查找TOKEN
	_, err := UserClient.FindByToken(c, &user.FindByTokenRequest{Token: toToken})
	if err != nil {
		err2, ok := err.(*errors.Error)
		if !ok {
			utils.ResponseError(c, utils.ErrInvalidToken, err)
		} else {
			utils.ResponseError(c, utils.ErrRPC, err2.Detail)
		}
		return
	}
	//本身维护的关系表中查找token
	ok, err, userGate := gateway_model.FindByToken(toToken)
	if !ok {
		utils.ResponseError(c, utils.ErrInvalidToken, "网关关系表中不存在该TOKEN")
		return
	}
	if err != nil {
		utils.ResponseError(c, utils.ErrCodeDb, err)
		return
	}
	//发布消息
	_, err = ImClient.PublishMessage(c, &im.PublishMessageRequest{
		FromToken:  fromToken,
		ToToken:   toToken,
		Body:       body,
		ServerName: userGate.ServerName,
		Topic:      userGate.Topic,
		Address:    userGate.ImAddress,
	})

	if err != nil {
		utils.ResponseError(c, utils.ErrRPC, err)
		return
	}
	utils.ResponseSuccess(c, "发送成功")
}

func GetServerAddress(c *gin.Context) {
	token := c.DefaultPostForm("token", "")
	if token == "" {
		utils.ResponseError(c, 1009, nil)
		return
	}

	res, err := UserClient.FindByToken(c, &user.FindByTokenRequest{Token: token})
	if err != nil {
		utils.ResponseError(c, 1009, err)
		return
	}
	if res.Token == "" {
		utils.ResponseError(c, 1009, "用户服务中未查找到该token")
		return
	}

	im_rpc := loadBalanceToAddress()
	err = gateway_model.RelationInsertUpdate(res.Token, im_rpc.ServerName, im_rpc.Topic, im_rpc.Address)
	if err != nil {
		utils.ResponseError(c, 1012, err)
		return
	}

	utils.ResponseSuccess(c, im_rpc.Address)
}

func loadBalanceToAddress() (imRpc *config.ImRpc) {
	all_weight := 0
	for _, i := range AddressList {
		all_weight += i.Weight
	}

	rand.Seed(time.Now().UnixNano())
	rand_w := rand.Intn(all_weight)
	cur_w := 0
	cur_ind := 0
	for ind, val := range AddressList {
		cur_w += val.Weight
		if rand_w-cur_w <= 0 {
			cur_ind = ind
			break
		}
	}

	return AddressList[cur_ind]
}
