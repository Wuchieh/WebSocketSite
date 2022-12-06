package Server

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

func tokenRecycle() {
	if setting.CheckExpiredTime > 0 {
		log.Println("已啟用定時回收Token")
	} else {
		return
	}
	go func() {
		for {
			time.Sleep(time.Duration(setting.CheckExpiredTime) * time.Second)
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
			if len(deleteTokenList) > 0 {
				log.Println("已經刪除", len(deleteTokenList), "個過期Token")
			}
		}
	}()
}

func tokenSave() {
	if setting.SaveTime > 0 {
		log.Println("已啟用定時儲存Token")
	} else {
		return
	}
	go func() {
		for {
			time.Sleep(time.Duration(setting.SaveTime) * time.Second)
			if len(userTokens) < 1 {
				err := os.WriteFile("tokens.json", []byte("{}"), 0666)
				if err != nil {
					log.Println()
					return
				}
				break
			}
			bytes, err := json.MarshalIndent(userTokens, "", "  ")
			if err != nil {
				log.Println()
				return
			}
			err = os.WriteFile("tokens.json", bytes, 0666)
			if err != nil {
				log.Println()
				return
			}
			log.Println("定時存檔完成")
		}
	}()
}
