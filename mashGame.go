package main

import (
	"net/http"
	"strconv"
	"time"
)

func mashGame(plData player, msg string) {
	// プレイ中の情報
	playing := false
	nowPushCount := 0
	timeLimit := 10
	name := "名無し"
	myRank := len(rankingData)

	//msgのコマンド読み取り
	cmd, _, cmdLen := readCmd(string(msg))

	if cmd[0] == "startGame" && cmdLen == 2 { //ゲーム開始コマンド。想定コマンド = startGame userName
		if cmd[1] != "" { //名前が空じゃなかったら、名前を更新
			name = cmd[1]
		}

		if !playing {
			timeLimit = 10
			playing = true
			nowPushCount = 0

			//ゲーム中の処理
			go func() {
				timeLimit = 10
				timer := time.NewTicker(time.Duration(1) * time.Second)
				for {
					<-timer.C
					timeLimit--
					if timeLimit == 0 { //プレイが終わったら次のプレイ準備をし、スコアの処理を行う
						playing = false
						myRank = updateRanking(name, nowPushCount)
						sendMsg(plData.conn, "rankingData "+strconv.Itoa(myRank)+" "+SliceToCsvStr(rankingData[:5]))
						return
					}

				}
			}()
		}

	} else if cmd[0] == "pushBtn" && cmdLen == 1 { //連打ボタンコマンド。想定コマンド = pushBtn
		if playing {
			nowPushCount += 1
		}
	} else if cmd[0] == "getRanking" && cmdLen == 1 { //ランキング取得コマンド。想定コマンド = getRanking 自分のスコア
		sendMsg(plData.conn, "rankingData "+strconv.Itoa(myRank)+" "+SliceToCsvStr(rankingData[:5]))
	}

}

type Files struct {
	ranking   string
	accessLog string
}

var files = Files{
	ranking:   "./data/mashGameRanking.csv",
	accessLog: "./data/mashGameAccessLog.txt",
}

var rankingData = ReadCsv(files.ranking)

// ランキング更新
func updateRanking(userName string, newScore int) int {
	addData := []string{userName, strconv.Itoa(newScore)}
	ranking := len(rankingData)
	for i, line := range rankingData {
		lineScore, _ := strconv.Atoi(line[1])

		if lineScore > newScore {
			continue
		}

		ranking = i
		break
	}

	slice1 := rankingData[:ranking]
	slice2 := [][]string{addData}
	slice3 := rankingData[ranking:]
	slice2 = append(slice2, slice3...)
	rankingData = append(slice1, slice2...)
	WriteCsv(files.ranking, rankingData)
	return ranking + 1
}

func MashGameCmd(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	//リモートアドレスからのアクセスを許可する
	conn, _ := upgrader.Upgrade(w, r, nil)

	log := time.Now().Format("2006/1/2 15:04:05") + " | " + r.RemoteAddr
	WriteFileAppend(files.accessLog, log)

	//プレイ中の情報
	playing := false
	nowPushCount := 0
	timeLimit := 10
	name := "名無し"
	myRank := len(rankingData)

	// 無限ループさせ、接続が切れないようにする
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil { //通信終了時の処理
			break
		}

		//msgのコマンド読み取り
		cmd, _, cmdLen := readCmd(string(msg))

		if cmd[0] == "startGame" && cmdLen == 2 { //ゲーム開始コマンド。想定コマンド = startGame userName
			if cmd[1] != "" { //名前が空じゃなかったら、名前を更新
				name = cmd[1]
			}

			if !playing {
				timeLimit = 10
				playing = true
				nowPushCount = 0

				//ゲーム中の処理
				go func() {
					timeLimit = 10
					timer := time.NewTicker(time.Duration(1) * time.Second)
					for {
						<-timer.C
						timeLimit--
						if timeLimit == 0 { //プレイが終わったら次のプレイ準備をし、スコアの処理を行う
							playing = false
							myRank = updateRanking(name, nowPushCount)
							sendMsg(conn, "rankingData "+strconv.Itoa(myRank)+" "+SliceToCsvStr(rankingData[:5]))
							return
						}

					}
				}()
			}

		} else if cmd[0] == "pushBtn" && cmdLen == 1 { //連打ボタンコマンド。想定コマンド = pushBtn
			if playing {
				nowPushCount += 1
			}
		} else if cmd[0] == "getRanking" && cmdLen == 1 { //ランキング取得コマンド。想定コマンド = getRanking 自分のスコア
			sendMsg(conn, "rankingData "+strconv.Itoa(myRank)+" "+SliceToCsvStr(rankingData[:5]))
		}
	}
}
