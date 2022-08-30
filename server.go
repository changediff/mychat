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

//创建server接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

//监听Message，有消息就发送给所有User
func (this *Server) ListenMessager() {
	for {
		msg := <-this.Message

		//将msg发送给全部在线User
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

// 广播消息
func (this *Server) Broadcast(user *User, msg string) {
	sendMsg := "[ 世界消息 ]" + user.Name + ": " + msg + "\n"
	this.Message <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
	//...当前连接任务
	// fmt.Println("当前链接成功")

	user := NewUser(conn, this)

	user.Online()

	//监听用户是否活跃的channel
	isLive := make(chan bool)

	//接收客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err: ", err)
				return
			}

			//提取用户的消息，去除"\n"
			msg := string(buf[:n-1])

			//用户处理msg
			user.DoMessage(msg)

			//用户发消息，表示当前用户存活
			isLive <- true
		}
	}()

	//当前handler阻塞
	for {
		select {
		case <-isLive:
			//当前用户活跃，应该重置定时器
			//不做任何事情，为了激活select，更新下面的定时器，所以isLive需要写在前面
		case <-time.After(time.Second * 600):
			//已经超时10s，将当前的User强制关闭
			user.SendMsg("超时被踢")

			//销毁资源
			close(user.C)

			//关闭链接
			conn.Close()

			//退出当前Handler
			return //rutime.GoExit

		}
	}
}

//启动服务器的接口
func (this *Server) Start() {
	// socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen error:", err)
		return
	}

	// close listen socket
	defer listener.Close()

	//启动监听Message的goroutine
	go this.ListenMessager()

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept error: ", err)
			continue
		}

		// do handler
		go this.Handler(conn)
	}

}
