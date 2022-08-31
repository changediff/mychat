package main

import (
	"net"
	"strings"
)

type User struct {
	Name          string
	Addr          string
	Email         string
	Password      string
	LastReplyUser *User
	C             chan string
	conn          net.Conn
	server        *Server
}

//创建用户API
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:     userAddr,
		Addr:     userAddr,
		Email:    "",
		Password: "",
		C:        make(chan string),
		conn:     conn,
		server:   server,
	}

	//启动监听当前user channel消息的goroutine
	go user.ListenMessage()
	return user
}

func (this *User) Online() {
	//用户上线，将用户加入到onlineMap中
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	//广播上线消息
	this.server.Broadcast(this, "已上线")
}

func (this *User) Offline() {
	//用户下线，将用户从onlineMap中删除
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	//广播下线消息
	this.server.Broadcast(this, Colorize("已下线", FgBlack, BgWhite))
}

//给当前User的客户端发消息，不群发
func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg + "\n"))
}

func (this *User) DoMessage(msg string) {
	if msg == "/help" {
		context := "/help 帮助\n/who 查询在线用户\n/w <用户名> <消息> 私聊\n/r <消息> 快速回复私聊\n/rename <用户名> 改名"
		this.SendMsg(Colorize(context, FgBlack, BgWhite))
	} else if msg == "/who" {
		//查询当前在线用户都有哪些
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "在线..."
			this.SendMsg(onlineMsg)
		}
		this.server.mapLock.Unlock()
	} else if len(msg) >= 6 && msg[:3] == "/w " {
		//消息格式： /w 阿宁 你好呀

		//1 获取对方用户名
		remoteName := strings.Split(msg, " ")[1]
		if remoteName == "" {
			this.SendMsg("私聊消息格式不正确，请使用\"/w 姓名 消息\"格式。")
			return
		}

		//2 根据用户名 得到对方User对象
		remoteUser, ok := this.server.OnlineMap[remoteName]
		if !ok {
			this.SendMsg("该用户名不存在，请使用\"/who\"命令确认")
			return
		}

		//3 获取消息内容，通过对方的User对象将消息私发出去
		content := strings.Join(strings.Split(msg, " ")[2:], " ")
		if content == "" {
			this.SendMsg("消息内容为空，请重新输入")
			return
		}
		remoteUser.SendMsg(Colorize("[ 来自 "+this.Name+" ]"+": "+content, FgBlue, BgDefault))
		remoteUser.LastReplyUser = this

	} else if len(msg) >= 9 && msg[:8] == "/rename " {
		//改名消息格式： /rename 阿成
		newName := strings.Split(msg, " ")[1]

		//判断name是否存在
		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.SendMsg(Colorize("当前用户名已被占用，请尝试其他用户名", FgBlack, BgYellow))
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()

			this.Name = newName
			this.SendMsg(Colorize("用户名已更新为:"+this.Name, FgBlack, BgWhite))
		}
	} else if len(msg) >= 4 && msg[:3] == "/r " {
		content := strings.Join(strings.Split(msg, " ")[1:], " ")
		if this.LastReplyUser == nil {
			this.SendMsg(Colorize("最近无私聊消息，无法快速回复", FgBlack, BgWhite))
		} else {
			this.LastReplyUser.SendMsg(Colorize("[ 来自 "+this.Name+" ]"+": "+content, FgBlue, BgDefault))
			this.LastReplyUser.LastReplyUser = this //套娃🪆，怎么简化？
		}
	} else if msg == "" {
		this.SendMsg(Colorize("请输入消息", FgBlack, BgYellow))
	} else {
		this.server.Broadcast(this, msg)
	}
}

//监听当前User channel的方法，有消息就发给对端客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg))
	}
}
