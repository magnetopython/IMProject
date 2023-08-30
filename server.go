package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	//在线用户列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	//消息广播的channel
	Message chan string
}

// 监听Message广播消息channel的goroutine,一旦有消息就发送给全部的在线User
func (this *Server) ListenMessager() {
	for {
		msg := <-this.Message
		//将msg发送给全部的在线user

		this.mapLock.Lock()

		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}

		this.mapLock.Unlock()
	}
}

// 创建server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// 广播消息
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	this.Message <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
	//当前连接的业务
	fmt.Println("连接建立成功")

	user := NewUser(conn, this)

	user.Online()

	//接受客户端发送的消息
	isLive := make(chan bool)

	//接受客户端传递的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("conn read err:", err)
				return
			}
			//提取用户的消息('\n')
			msg := string(buf[:n-1])

			//将得到的消息进行广播
			user.DoMessage(msg)

			//用户的任意消息，代表当前用户是活跃的
			isLive <- true
		}
	}()
	//当前handler阻塞
	for {
		select {
		case <-isLive:
			//当前用户是活跃的，应该重置定时器
			//为了激活select,更新下面的定时器
		case <-time.After(time.Second * 10):
			//已经超时
			//将当前的User强制关闭
			user.SendMsg("you have been kicked off")

			//销毁用户的资源
			close(user.C)
			//关闭链接
			conn.Close()

			//退出当前Handler
			return
		}
	}
}

// 启动服务器的接口
func (this *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	//accept
	if err != nil {
		fmt.Println("net Listen err:", err)
		return
	}
	defer listener.Close()

	//启动监听Message的goroutine
	go this.ListenMessager()

	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}
		//do handler
		go this.Handler(conn)
	}
	//close listen socket
}
