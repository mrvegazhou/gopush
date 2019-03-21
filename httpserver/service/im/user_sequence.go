package imService

import (
	"gopush/framework/db"
	"gopush/framework/helper"
	"gopush/httpserver/models/im"
)

type userSequenceService struct{}

var UserRequenceService = new(userSequenceService)

// GetNext 获取下一个序列
func (*userSequenceService) GetNext(session *db.Session, userId int64) (int64, error) {
	db := session.DB
	tx := db.Begin()
	err := imModel.UserSequenceDao.Increase(session, userId)
	helper.PrintErrRollback(err, tx)

	seq, err := imModel.UserSequenceDao.GetSequence(session, userId)
	helper.PrintErrRollback(err, tx)

	return seq, nil
}
