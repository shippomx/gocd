package main

import (
	"fmt"
	"gocd/keyboard"
)

func init() {
	fmt.Println("Ready to receive keyborad input.")
}

func main() {
	go keyboard.PrintLine() // 无缓冲的chan 接收要出现在发送之前
	keyboard.Readline()
}

