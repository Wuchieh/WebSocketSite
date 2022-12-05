package Server

import (
	"fmt"
	"time"
)

func tokenRecycle() {
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			timeNow := time.Now().Unix()
			var deleteTokenList []string
			for i, v := range userTokens {
				if v.ExpiredTime.Unix() < timeNow {
					deleteTokenList = append(deleteTokenList, i)
				}
			}
			for _, v := range deleteTokenList {
				delete(userTokens, v)
			}
			fmt.Println("已經刪除", len(deleteTokenList), "個過期Token")
		}
	}()
}
