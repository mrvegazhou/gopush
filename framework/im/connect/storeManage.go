package connect

import (
	"encoding/json"
	"fmt"
	"gopush/framework/db/redis"
)

func init() {
	redigoUtil.NewConnectionWithConf(ConfTcpServer)
}

func store(deviceId int64, ctx *ConnContext) error{
	res, err := json.Marshal(ctx)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = redigoUtil.RedigoConn.String.Set(string(deviceId), res).Error()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func load(deviceId int64) (*ConnContext, error) {
	val, err := redigoUtil.RedigoConn.String.Get(string(deviceId)).String()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var ctx ConnContext
	err = json.Unmarshal([]byte(val), &ctx)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &ctx, nil
}

func delete(deviceId int64) {
	redigoUtil.RedigoConn.String.Delete(string(deviceId))
}