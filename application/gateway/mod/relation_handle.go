package gateway_model

func RelationInsertUpdate (token,serverName,topic,imAddress string)(err error){
	gate := &GateWay{}
	ok, err := engine.Where("token = ?", token).Get(gate)
	if err != nil {
		return err
	}
	gate.Topic =topic
	gate.ServerName = serverName
	gate.ImAddress = imAddress
	if ok {
		_, err := engine.Update(gate)
		if err != nil {
			return err
		}
		return nil
	}
	gate.Token = token
	_, err = engine.Insert(gate)
	if err != nil {
		return err
	}

	return nil
}

func FindByToken(token string)(ok bool, err error,gate *GateWay){
	gate = &GateWay{}
	ok, err = engine.Where("token = ?", token).Get(gate)
	if err != nil {
		return
	}
	return
}
