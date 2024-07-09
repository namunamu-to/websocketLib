package main

import (
	"crypto/tls"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var accessLogFilepath = "./data/accessLog.txt"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// CheckOrigin: func(r *http.Request) bool {
	// 	return true
	// },
}

type player struct {
	uid, name, roomKey string
	conn               *websocket.Conn
}

type room struct {
	players []player
}

var rooms = make(map[string]*room) //roomKeyでルーム指定

var mux *http.ServeMux

func MakeUuid() string {
	id := uuid.New()
	return id.String()
}

func getPlayerIdx(roomKey string, uid string) int {
	players := rooms[roomKey].players
	for i := 0; i < len(players); i++ {
		if players[i].uid == uid {
			return i
		}
	}

	return -1 //見つからなかった時
}

func enterRoom(roomKey string, player player) {
	idx := getPlayerIdx(roomKey, player.uid)
	if idx == -1 { //まだ自分が部屋に入ってなかったら追加
		rooms[roomKey].players = append(rooms[roomKey].players, player)
	}
}

func exitRoom(roomKey string, plData player) {
	pIdx := getPlayerIdx(roomKey, plData.uid)
	if pIdx != -1 { //自分が部屋にいたら部屋から抜ける
		a := rooms[roomKey].players
		a[pIdx] = a[len(a)-1]
		a = a[:len(a)-1]
		rooms[roomKey].players = a
	}

	if len(rooms[roomKey].players) == 1 { //他プレイヤーへ退出したことを通知
		sendMsg(rooms[roomKey].players[0].conn, "disConnect")
	}
}

func sendMsg(conn *websocket.Conn, msg string) {
	conn.WriteMessage(websocket.TextMessage, []byte(msg))
}

func readCmd(str string) ([]string, string, int) {
	cmd := strings.Split(string(str), " ")
	cmdType := cmd[0]
	cmdLen := len(cmd)
	return cmd, cmdType, cmdLen
}

func bloadcastMsg(roomKey string, msg string) {
	for i := 0; i < len(rooms[roomKey].players); i++ {
		sendMsg(rooms[roomKey].players[i].conn, msg)
	}
}

func isRoom(roomKey string) bool {
	_, ok := rooms[roomKey]
	return ok
}

func makeRoom(roomKey string) {
	rooms[roomKey] = &room{players: []player{}}
}

func sendMsgToAnother(roomKey string, exceptPl player, msg string) {
	//自分以外にコマンド送信
	exceptIdx := getPlayerIdx(roomKey, exceptPl.uid)
	toIdx := 1 - exceptIdx
	sendMsg(rooms[roomKey].players[toIdx].conn, msg)
}

func addHandleFunc(url string, fn func()) {
	handller := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "https://galleon.yachiyo.tech")

		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		//リモートアドレスからのアクセスを許可する
		conn, _ := upgrader.Upgrade(w, r, nil)

		//ログファイルにアクセスを記録
		log := time.Now().Format("2006/1/2 15:04:05") + " | " + r.RemoteAddr + url
		WriteFileAppend(accessLogFilepath, log)

		// 無限ループさせることでクライアントからのメッセージを受け付けられる状態にする
		plData := player{uid: MakeUuid(), name: "", roomKey: "matching", conn: conn}
		roomKey := plData.roomKey

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil { //通信終了時の処理
				if roomKey == "matching" {
					break
				}

				exitRoom(roomKey, plData) //部屋から抜ける
				break
			}

			//msgのコマンド読み取り
			cmd, cmdType, cmdLen := readCmd(string(msg))
			playerNum := 0
			println("ルーム内プレイヤー数", playerNum)

			//コマンドに応じた処理をする
			if cmdType == "roomMatch" && cmdLen == 2 { //マッチングコマンド。想定コマンド = "roomMatch ルームキー"
				roomKey = cmd[1]

				if !isRoom(roomKey) { //部屋が無いなら作る
					sendMsg(conn, "部屋作った")
					makeRoom(roomKey)
				}

				playerNum = len(rooms[roomKey].players)
				enterRoom(roomKey, plData)
				bloadcastMsg(roomKey, "playerNum "+strconv.Itoa(playerNum))
				sendMsgToAnother(roomKey, plData, "なんか来たw")
			}

			fn()
		}

	}

	mux.HandleFunc(url, handller)
}

func main() {
	// ハンドラの設定
	mux = http.NewServeMux() //ミューテックス。すでに起動してるか確認。
	addHandleFunc("/shogi/websocketLib", func() { println("ハンドル実行") })

	//tls設定
	cfg := &tls.Config{
		ClientAuth: tls.RequestClientCert,
	}

	//サーバー設定
	srv := http.Server{
		Addr:      ":8443",
		Handler:   mux,
		TLSConfig: cfg,
	}

	println("サーバー起動")
	err := srv.ListenAndServeTLS("./fullchain.pem", "./privkey.pem")
	if err != nil {
		println("サーバー起動に失敗")
		println(err.Error())
	}
}
