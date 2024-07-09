package main

import "strconv"

//オウム返し
// func test1(plData player, msg string) {
// 	sendMsg(plData.conn, msg)
// }

// 特定のコマンドなら、オウム返しする
func test2(plData player, msg string) {
	cmd, cmdType, cmdLen := readCmd(msg)
	if cmdType == "test1" && cmdLen == 2 {
		sendMsg(plData.conn, cmd[1])
	}
}

// 挨拶を返す
func test3(plData player, msg string) {
	if msg != "hello" {
		return
	}

	sendMsg(plData.conn, "hello!")
}

// プレイヤーデータを用いてルーム情報にアクセスし、クライアントに返す
func test4(plData player, msg string) {
	_, cmdType, _ := readCmd(msg)
	if cmdType != "roomInfo" {
		return
	}

	if plData.roomKey == "" {
		sendMsg(plData.conn, "まだルームに入っていません")
		return
	}

	response := len(rooms[plData.roomKey].players)
	sendMsg(plData.conn, "ルーム人数 : "+strconv.Itoa(response))
}

func SetHandllers() {
	// addedHandllers = append(addedHandllers, test1)
	addedHandllers = append(addedHandllers, test2)
	addedHandllers = append(addedHandllers, test3)
	addedHandllers = append(addedHandllers, test4)
}
