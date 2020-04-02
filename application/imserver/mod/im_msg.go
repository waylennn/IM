package mod

func SaveOfflineMsg (Body,FromToken,ToToken string,TimeStamp int64) error{
	offlineMsg := &OfflineMsg{}
	offlineMsg.FromToken = Body
	offlineMsg.ToToken = FromToken
	offlineMsg.Body = ToToken
	offlineMsg.TimeStamp = TimeStamp

	_,err := engine.Insert(offlineMsg)
	if err != nil {
		return err
	}

	return nil
}

func GetOffLineMsg(receiver string)(offlineMsgList []*OfflineMsg,err error){
	offlineMsgList = make([]*OfflineMsg, 0)
	session := engine.Where("receiver = ?",receiver)
	err = session.Find(&offlineMsgList)
	if err != nil {
		return
	}

	_,err = session.Where("receiver = ?",receiver).Delete(&OfflineMsg{})
	if err != nil {
		return
	}

	return
}