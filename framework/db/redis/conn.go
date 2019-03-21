package redigoUtil

import (
	"errors"
	"gopush/conf"
	"gopush/const"
	"gopush/framework/log"
)

type redigoUtil struct {
	String stringRds
	List   listRds
	Hash   hashRds
	Key    keyRds
	Set    setRds
	ZSet   zSetRds
	Bit    bitRds
	Db     dbRds
}

var RedigoConn = new(redigoUtil)

func NewConnectionWithConf(config *conf.MainConfig) (*redigoUtil, error) {
	if config == nil {
		err := constdefine.GetMsg(constdefine.ERROR_CONFIG)
		logger.Sugar.Error(err)
		return nil, errors.New(err)
	}

	initPool(config)
	logger.Sugar.Info("redis 配置为: ", config)
	return RedigoConn, nil
}
