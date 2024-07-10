package main

import "strconv"

//オウム返し
// func handllerEx1(plData player, msg string) {
// 	sendMsg(plData.conn, msg)
// }

// 特定のコマンドなら、オウム返しする
func handllerEx2(plData player, msg string) {
	cmd, cmdType, cmdLen := readCmd(msg)
	if cmdType == "test1" && cmdLen == 2 {
		sendMsg(plData.conn, cmd[1])
	}
}

// 挨拶を返す
func handllerEx3(plData player, msg string) {
	if msg != "hello" {
		return
	}

	sendMsg(plData.conn, "hello!")
}

// プレイヤーデータを用いてルーム情報にアクセスし、クライアントに返す
func handllerEx4(plData player, msg string) {
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

// チャット
func handllerEx5(plData player, msg string) {
	_, cmdType, _ := readCmd(msg)
	if cmdType != "chat" {
		return
	}

	if plData.roomKey == "" {
		sendMsg(plData.conn, "まだルームに入っていません")
		return
	}

	if cmdType == "chat" {
		bloadcastMsg(plData.roomKey, msg[5:])
	}
}

func main() {
	addHandller(handllerEx2)
	addHandller(handllerEx3)
	addHandller(handllerEx4)
	addHandller(handllerEx5)
	addHandller(mashGame)
	startServer("/test", "8444", "./fullchain.pem", "./privkey.pem")
}
