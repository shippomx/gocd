package keyboard

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

var inDirs = []string{}
var ChanInput chan string
var ChanClose chan int

func init() {
	ChanInput = make(chan string, 0)
	ChanClose = make(chan int, 0)

	homeDir, _ := os.UserHomeDir()
	inDirs = strings.Split(homeDir[1:], "/")
}


// 作为无缓冲channel通道的内容生产者
func Readline() {
	_readLine()
}

func _readLine() {
	curPath	:= ""
	var err error
	scanner := bufio.NewScanner(os.Stdin)
	for {
		// 作为无缓冲channel通道的内容生产者
		fmt.Printf("%s", "$ ")
		scanner.Scan()
		cmd := scanner.Text()

		if cmd == "exit" {
			ChanClose <- 1
		} else if strings.HasPrefix(cmd, "cd ") {
			args := strings.Split(cmd, " ")
			if len(args) != 2 {
				continue
			}
			curPath, err = assemblePath(args[1])
			if err != nil {
				fmt.Println(err)
				continue
			}
		} else if cmd == "pwd" {
			curPath, err = assemblePath(".")
			ChanInput <- curPath
		}
	}
}

func PrintLine() {
	for {
		select {
		case curPath := <- ChanInput:
			fmt.Println(curPath)
		case <-ChanClose:
			close(ChanInput)
			os.Exit(0)
		default:
		}
	}
}

func assemblePath(cmd string) (string, error) {
	if cmd == "" {
		cmd = "."
	}
	idx := strings.IndexAny(cmd, "~") // homeDir 只能为第一个字符
	if idx > 0 {
		return "", errors.New("cd: no such file or directory: " + cmd)
	}
	if strings.Contains(cmd, "...") { // 不考虑连续多个的...
		return "", errors.New("cd: no such file or directory: " + cmd)
	}
	if string(cmd[0]) == "/" {
		inDirs = []string{} //代表从根目录进行计算统计
	}
	dirs := strings.Split(cmd, "/")
	for _, dir := range dirs {
		switch dir {
		case "":
			continue
		case ".":
		case "..":
			if len(inDirs) < 1 {
				inDirs = []string{}
			} else {
				inDirs = inDirs[:len(inDirs)-1]
			}
		default:
			inDirs = append(inDirs, dir)
		}
	}
	curPath := ""
	for _, dir := range inDirs {
		curPath += "/" + dir
	}
	if curPath == "" {
		curPath = "/"
	}
	return curPath, nil
}