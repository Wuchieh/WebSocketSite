package Server

import (
	"math/rand"
	"time"
)

func RandomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
	//a := string(b)
	//return &a
}

func genNewToken(session string) *UserToken {
	if token, ok := userTokens[session]; ok {
		token.UpdateTime = time.Now()
		token.ExpiredTime = time.Now().Add(time.Minute * time.Duration(setting.ExpiredTime))
		token.Token = RandomString(32)
	} else {
		userToken := UserToken{
			CreateTime:  time.Now(),
			UpdateTime:  time.Now(),
			ExpiredTime: time.Now().Add(time.Minute * time.Duration(setting.ExpiredTime)),
			Token:       RandomString(32),
		}
		userTokens[session] = &userToken
	}
	return userTokens[session]
}

func updateToken(session string) bool {
	if token, ok := userTokens[session]; ok {
		if token.UpdateTime.Add(1*time.Minute).Unix() > time.Now().Unix() {
			return false
		}
		token.UpdateTime = time.Now()
		token.ExpiredTime = time.Now().Add(time.Minute * time.Duration(setting.ExpiredTime))
		token.Token = RandomString(32)
	} else {
		if genNewToken(session) == nil {
			return false
		}
	}
	return true
}

//func wsLogin(ws *websocket.Conn, login *struct {
//	Token string `json:"Token"`
//	Group string `json:"Group"`
//}) bool {
//	userToken, ok := tokenCheck(login.Token)
//	if !ok || userToken.InGroup {
//		return false
//	}
//	err := wsConnectGroupJoin(ws, userToken, login.Group)
//	if err != nil {
//		log.Println(err)
//		return false
//	}
//	return true
//}
//
//func wsConnectGroupJoin(ws *websocket.Conn, token *UserToken, g string) error {
//	if group, ok := wsConnectGroups[g]; ok {
//		group = append(group, token)
//	} else {
//		wsConnectGroups[g] = []*UserToken{}
//		wsConnectGroups[g] = append(wsConnectGroups[g], token)
//	}
//	token.Ws = ws
//	token.InGroup = true
//	token.Group = g
//	return nil
//}
//
//func tokenCheck(token string) (*UserToken, bool) {
//	for _, userToken := range userTokens {
//		if userToken.Token == token {
//			return userToken, true
//		}
//	}
//	return nil, false
//}
//
//func wsLoginAuthentication(ws *websocket.Conn) (bool, string) {
//	for i := 0; i < 3; i++ {
//		msgType, msg, err := ws.ReadMessage()
//		if err != nil || msgType == -1 {
//			log.Println(err)
//			_ = ws.WriteMessage(1, []byte("登入失敗請重新嘗試"))
//			continue
//		}
//		login := &struct {
//			Token string `json:"Token"`
//			Group string `json:"Group"`
//		}{}
//		err = json.Unmarshal(msg, &login)
//		if err != nil {
//			log.Println(err)
//			_ = ws.WriteMessage(1, []byte("登入失敗請重新嘗試"))
//			continue
//		}
//		if login.Token != "" && login.Group != "" {
//			if wsLogin(ws, login) {
//				return true, login.Group
//			} else {
//				return false, ""
//			}
//		}
//		_ = ws.WriteMessage(1, []byte("登入失敗請重新嘗試"))
//	}
//	return false, ""
//}
//
//func wsLogout(ws *websocket.Conn, group string) {
//	for i, v := range wsConnectGroups[group] {
//		if v.Ws == ws {
//			wsConnectGroups[group][i].InGroup = false
//			wsConnectGroups[group][i].Group = ""
//			wsConnectGroups[group] = append(wsConnectGroups[group][:i], wsConnectGroups[group][i+1:]...)
//			break
//		}
//	}
//	if len(wsConnectGroups[group]) < 1 {
//		delete(wsConnectGroups, group)
//	}
//}
