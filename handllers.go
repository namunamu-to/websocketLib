package main

// func test1(plData player, msg string) {
// 	sendMsg(plData.conn, msg)
// }

func test2(plData player, msg string) {
	cmd, cmdType, cmdLen := readCmd(msg)
	if cmdType == "test1" && cmdLen == 2 {
		sendMsg(plData.conn, cmd[1])
	}
}

func test3(plData player, msg string) {
	if msg != "hello" {
		return
	}

	sendMsg(plData.conn, "hello!")
}

func SetHandllers() {
	// addedHandllers = append(addedHandllers, test1)
	addedHandllers = append(addedHandllers, test2)
	addedHandllers = append(addedHandllers, test3)
}
