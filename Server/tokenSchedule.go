package Server

import (
	"encoding/json"
	"github.com/pkg/errors"
	"log"
	"os"
	"time"
)

func tokenSchedule() {
	ScheduleChannel <- 1
	var err error
	printError := func(err error) {
		if err != nil {
			log.Println(err)
		}
	}
	go func() {
		var running bool
		for {
			sc := <-ScheduleChannel
			err = checkTokenExpired()
			printError(err)
			err = saveUserTokens()
			printError(err)
			if sc == 0 {
				ScheduleChannel <- 0
				return
			}
			go func() {
				if running {
					return
				}
				running = true
				time.Sleep(time.Minute * time.Duration(setting.ScheduleTime))
				ScheduleChannel <- 1
				running = false
			}()
		}
	}()

}

func saveUserTokens() (err error) {
	// 將所有的UserToken以json的形式進行存檔
	log.Println("正在存檔")
	errorOccurred := true
	defer func() {
		log.Println("存檔完畢")
		if err != nil {
			return
		}
		if errorOccurred {
			err = errors.Errorf("發生錯誤")
		}
	}()
	var bytes []byte
	if len(userTokens) < 1 {
		bytes, err = json.Marshal([]byte("{}"))
	} else {
		bytes, err = json.MarshalIndent(userTokens, "", "  ")
	}
	if err != nil {
		return
	}
	err = os.WriteFile("tokens.json", bytes, 0666)
	errorOccurred = false
	return
}

func checkTokenExpired() (err error) {
	// 將已經過期的UserToken刪除
	errorOccurred := true
	defer func() {
		if errorOccurred {
			err = errors.Errorf("發生錯誤")
		}
	}()
	timeNowUnix := time.Now().Unix()
	var expiredUserToken []string

	for key, value := range userTokens {
		if timeNowUnix > value.ExpiredTime.Unix() {
			expiredUserToken = append(expiredUserToken, key)
			delete(userTokens, key)
		}
	}
	errorOccurred = false
	if len(expiredUserToken) > 0 {
		log.Println("已刪除", len(expiredUserToken), "個UserToken")
	}
	return
}
