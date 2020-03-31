package user

import (
	"awesomeProject/application/userapi/utils"
	users_model "awesomeProject/application/userserver/mod"
	user "awesomeProject/application/userserver/protoc"
	"github.com/micro/go-micro/v2/errors"
	"golang.org/x/net/context"
	"net/http"
)

type User struct{}

func (u *User) FindByToken(ctx context.Context, req *user.FindByTokenRequest, res *user.UserResponse) error {
	ok,err,us := users_model.FindByToken(req.Token)
	if err != nil {
		return err
	}
	res.Token = us.Token
	res.Username = us.UserName
	//
	if !ok {
		err := &errors.Error{
			Code:   utils.ErrInvalidToken,
			Status: http.StatusText(500),
			Detail:"token不存在",
		}
		return err
	}
	return nil
}
func (u *User) FindById(ctx context.Context, req *user.FindByIdRequest, res *user.UserResponse) error {
	return nil
}
