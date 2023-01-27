package main

import (
	"fmt"
	"net"
	"os"
	"sync"
)

var user = 0
var key string
var USERNAME string
var wg sync.WaitGroup

func main() { //server
	ADDRESS, PORT, uSERNAME := config()
	USERNAME = uSERNAME
	IP := ADDRESS + ":" + PORT
	serv, err := net.Listen("tcp", IP)
	if err != nil {
		fmt.Println("-----------error----------\n", err, "\n--------------------------------")
		os.Exit(1)
	}
	defer serv.Close()
	fmt.Println("listen start...")
	for {
		conn, err := serv.Accept()
		if err != nil {
			fmt.Println("-----------error----------\n", err, "\n--------------------------------")
		}
		if user == 0 {
			user = 1
			go conne(conn)
			user = 0
		} else if user == 1 {
			go out(conn)
		}
	}
}
func conne(conn net.Conn) {
	var buf [15]byte
	r, err := conn.Read(buf[:])
	if err != nil {
		fmt.Println("-----------error----------\n", err, "\n--------------------------------")
	}
	fmt.Println(string(buf[:r]))
	ok(conn)
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("-----------error----------\n", err, "\n--------------------------------")
		}
	}(conn)
	fmt.Println("connect close")
	key = ""
	user = 0
}
func out(conn net.Conn) {
	_, err := conn.Write([]byte("1"))
	if err != nil {
		fmt.Println("-----------error----------\n", err, "\n--------------------------------")
	}
}
func ok(conn net.Conn) {
	_, err := conn.Write([]byte("ok"))
	if err != nil {
		fmt.Println("-----------error----------\n", err, "\n--------------------------------")
	}
	fmt.Println("-----------自检结束，开始远程命令行协助-----------")
	fmt.Println("输入24位密匙(如果聊天消息乱码说明有一方输入错误)")
	_, err = fmt.Scan(&key)
	if err != nil {
		fmt.Println("-----------error----------\n", err, "\n--------------------------------")
	}
	fmt.Println("直接输入命令或-s 后跟消息以发送")
	wg.Add(2)
	go lis(conn)
	go say(conn)
	wg.Wait()
}
func lis(conn net.Conn) {
	var buf [200]byte
	for {
		_, _ = conn.Read(buf[:])
		r, err := conn.Read(buf[:])
		go li(conn, buf, r)
		if err != nil {
			fmt.Println("-----------error----------\n", err, "\n--------------------------------")
		}
	}
	wg.Done()
}

func li(conn net.Conn, buf [200]byte, r int) {
	if AesDecrypt(string(buf[:r]), key) == "exit" {
		err := conn.Close()
		if err != nil {
			if err != nil {
				fmt.Println("-----------error----------\n", err, "\n--------------------------------")
			}
		}
		return
	} else {
		fmt.Println(AesDecrypt(string(buf[:r]), key))
	}
}

func say(conn net.Conn) {
	var ssay string
	sh := 0
	for {
		_, err := fmt.Scan(&ssay)
		if err != nil {
			fmt.Println("-----------error----------\n", err, "\n--------------------------------")
		}
		if ssay == "-s" {
			sh = 0
		} else if ssay == "-ss" {
			sh = 1
		}
		go sa(ssay, sh, conn)
	}
	wg.Done()
}

func sa(ssay string, sh int, conn net.Conn) {
	var s string
	if ssay == "-s" || ssay == "-ss" || sh == 1 {
		s = ssay
	} else {
		s = USERNAME + " : " + ssay
	}
	_, err := conn.Write([]byte(AesEncrypt(s, key)))
	if err != nil {
		fmt.Println("-----------error----------\n", err, "\n--------------------------------")
	} else {
		fmt.Println("{" + USERNAME + " : " + ssay + "}")
	}
}
