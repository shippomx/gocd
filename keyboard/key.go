package keyboard

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

var inDirs = []string{"Users", "moumooun"}
var ChanInput chan string
var ChanClose chan int

func init() {
	ChanInput = make(chan string, 0)
	ChanClose = make(chan int, 0)
}


// 作为无缓冲channel通道的内容生产者
func Readline() {
	_readLine()
}

func _readLine() {
	//reader := bufio.NewReader(os.Stdin)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		// 作为无缓冲channel通道的内容生产者
		scanner.Scan()
		cmd := scanner.Text()
		if cmd == "exit" {
			ChanClose <- 1
		} else if strings.HasPrefix(cmd, "cd ") {
			args := strings.Split(cmd, " ")
			if len(args) != 2 {
				continue
			}
			ChanInput <- args[1]
		}
	}
}

func PrintLine() {
	for {
		select {
		case cmd := <- ChanInput:
			curPath, err := assemblePath(cmd)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(curPath)
			}
		case <-ChanClose:
			close(ChanInput)
			os.Exit(0)
		default:
		}
	}
}

func assemblePath(cmd string) (string, error) {
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
	return curPath, nil
}