package helper

import "fmt"

func PrintErr(err error) error{
	if err != nil {
		fmt.Println(err)
		return err
	}
}
