package main

import (
	"gocd/keyboard"
)

func main() {
	go keyboard.PrintLine() // 无缓冲的chan 接收要出现在发送之前
	keyboard.Readline()
}

