package helper

import (
	_ "database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	"errors"
	"gopush/const"
)

func PrintErr(err error) error{
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func PrintErrRollback(err error, tx *gorm.DB) error {
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return err
	}
	return nil
}

func PrintErrMsg(code int, err error) error {
	if code >= 0 {
		msg := constdefine.GetMsg(constdefine.IM_ERROR_NOT_IN_GROUP)
		fmt.Println(msg)
		return errors.New(msg)
	}
	return nil
}