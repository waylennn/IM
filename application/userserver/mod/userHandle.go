package users_model

func FindByToken(token string)(ok bool, err error,user *Users){
	user = &Users{}
	ok, err = engine.Where("Token = ?",token).Get(user)
	if err != nil {
		return
	}
	return
}
