package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"sync"
	"time"
)

var c = 5
var key string
var wg sync.WaitGroup

func main() { //user
	ADDRESS, PORT, USERNAME := config()
	IP := ADDRESS + ":" + PORT
	conn, err := net.Dial("tcp", IP)
	if err != nil {
		fmt.Println("connection failed\ntrying again...\n----", c, "/5----")
		c--
		if c >= 0 {
			time.Sleep(10 * 100 * time.Millisecond)
			main()
		}
		fmt.Println("----exit----")
		os.Exit(1)
	}
	fmt.Println("connecting...")
	_, err = conn.Write([]byte("ping from user!\n"))
	if err != nil {
		fmt.Println("ping wrong\nexiting...")
		os.Exit(2)
	}
	fmt.Println("success!")
	var buf [20]byte
	r, err := conn.Read(buf[:])
	if err != nil {
		fmt.Println("人数信息读取错误")
		os.Exit(3)
	}
	if string(buf[:r]) == "1" {
		fmt.Println("人数已满，拒绝连接！(1/1)")
		os.Exit(4)
	} else if string(buf[:r]) == "ok" {
		_, err := conn.Write([]byte("ok"))
		if err != nil {
			fmt.Println("确认信息发送错误")
			os.Exit(3)
		}
		fmt.Println("-----------自检结束，开始远程命令行协助-----------")
		fmt.Println("输入24位密匙(如果聊天消息乱码说明有一方输入错误)")
		_, err = fmt.Scan(&key)
		if err != nil {
			os.Exit(4)
		}
		fmt.Println("直接输入并回车发送消息")
		wg.Add(2)
		go lis(conn)
		go say(conn, USERNAME)
		wg.Wait()
	}
}
func lis(conn net.Conn) {

	var buf [200]byte
	for {
		r, err := conn.Read(buf[:])
		go li(buf, r, conn)
		if err != nil {
			fmt.Println("接收错误")
		}
	}
	wg.Done()
}

var sh = 0

func li(buf [200]byte, r int, conn net.Conn) {
	var nr string
	nr = AesDecrypt(string(buf[:r]), key)
	if nr == "-s" {
		sh = 0
	} else if nr == "-ss" {
		sh = 1
	} else if sh == 0 {
		fmt.Println(nr)
	} else if sh == 1 {
		fmt.Println("执行命令...")
		go ml(nr, conn)
	}
}

func say(conn net.Conn, USERNAME string) {
	var ssay string
	for {
		_, err := fmt.Scan(&ssay)
		if err != nil {
			fmt.Println("输入错误")
		}
		go sa(conn, ssay, USERNAME)
	}
	wg.Done()
}

func sa(conn net.Conn, ssay string, USERNAME string) {
	if ssay == "exit" {
		_, err := conn.Write([]byte(AesEncrypt(ssay, key)))
		if err != nil {
			fmt.Println("发送错误")
		}
	} else {
		_, err := conn.Write([]byte(AesEncrypt("{"+USERNAME+" : "+ssay+"}", key)))
		if err != nil {
			fmt.Println("发送错误")
		} else {
			fmt.Println(USERNAME + " : " + ssay)
			_, _ = conn.Write([]byte("1"))
		}
	}
}

func ml(nr string, conn net.Conn) {
	shell := exec.Command(nr)
	out, err := shell.CombinedOutput()
	if err != nil {
		fmt.Println("执行错误")
	}
	_, err = conn.Write([]byte(AesEncrypt(string(out), key)))
	if err != nil {
		fmt.Println("发送错误")
	} else {
		_, _ = conn.Write([]byte("1"))
	}
}
