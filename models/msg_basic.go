package models

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
	"net"
	"net/http"
	"sync"
)

var rwLocker sync.RWMutex

type MessageBasic struct {
	gorm.Model
	FromIdentity string // 发送者
	ToIdentity   string // 接收者
	Type         string // 消息类型
	Media        int    //
	Content      string //
	Pic          string //
	Url          string //
	Desc         string //
	Amount       int    //
}

func (msg *MessageBasic) TableName() string {
	return "message_basic"
}

// 启动了的聊天节点
type Node struct {
	Conn      *websocket.Conn // websocket连接
	DataQueue chan []byte     // 消息队列
	GroupSets set.Interface   // ？
}

// 映射关系
var clientMap map[string]*Node = make(map[string]*Node, 0)

// todo: map的锁

func Chat(writer http.ResponseWriter, request *http.Request) {
	// todo: 校验token -> isValidated
	// token := query.Get("token")
	isValidated := true
	query := request.URL.Query()
	userId := query.Get("user_id")
	//msgType := query.Get("msg_type")
	//targetId := query.Get("target_id")
	//content := query.Get("content")
	// 升级连接为websocket
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isValidated
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe), // 保证线程
	}
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()
	go sendProc(node) // 将消息发送到由node表示的WebSocket连接。
	go recvProc(node) // 接收WebSocket连接上的消息并处理它们。
	sendDirectMessage(userId, []byte("欢迎进入聊天"))
}

func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func recvProc(node *Node) {
	for {
		_, message, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		broadMsg(message)
		fmt.Println(message)
	}
}

var udpSendChan = make(chan []byte, 1024)

func broadMsg(message []byte) {
	udpSendChan <- message
}

func init() {
	go udpSendProc() // 分别将udp
	go udpRecvProc()
}

func udpSendProc() {
	con, err := net.DialUDP(
		"udp",
		nil,
		&net.UDPAddr{
			IP:   net.IPv4(172, 31, 226, 255),
			Port: 3000,
		},
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(con *net.UDPConn) {
		err := con.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(con)
	for {
		select {
		// udpSendChan有消息则取出,并向websocket发送
		case data := <-udpSendChan:
			_, err := con.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func udpRecvProc() {
	con, err := net.ListenUDP("udp",
		&net.UDPAddr{
			IP:   net.IPv4zero,
			Port: 3000,
		})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(con *net.UDPConn) {
		err := con.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(con)
	for {
		var buf [512]byte
		n, err := con.Read(buf[0:])
		if err != nil {
			fmt.Println(err)
			return
		}
		dispatch(buf[0:n])
	}
}

func dispatch(data []byte) {
	msg := MessageBasic{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch msg.Type {
	case "1": // 私聊
		sendDirectMessage(msg.ToIdentity, data)
		//case 2:
		//	sendGroupMsg()

	}
}

func sendDirectMessage(targetId string, msg []byte) {
	rwLocker.RLock()
	node, ok := clientMap[targetId]
	rwLocker.RUnlock()
	if ok {
		node.DataQueue <- msg
	}
}
