package models

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
	"net/http"
	"sync"
)

var rwLocker sync.RWMutex

type MessageBasic struct {
	gorm.Model
	TargetId   string // 目标
	UserId     string // 谁发送的
	CreateTime int64  // 发送时间
	Type       int    // 消息类型
	Media      int    // 消息的内容类型，1、文本 4、表情
	Content    string // 内容
	Pic        string //
	Url        string // 表情URL
	Desc       string //
	Amount     int    //
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
	userId := query.Get("userId")
	// 升级http连接升级为websocket
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isValidated
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 创建聊天节点
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe), // 保证线程
	}
	rwLocker.Lock()
	// 放入map中
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
			fmt.Println("send >>>>>>>>", string(data))
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
		//broadMsg(message)
		fmt.Println("<<<<<<<<recv", string(message))
		dispatch(message)

	}
}

var udpSendChan = make(chan []byte, 1024)

func broadMsg(message []byte) {
	udpSendChan <- message
}

//func init() {
//	go udpSendProc() // 分别将udp
//	go udpRecvProc()
//}
//
//func udpSendProc() {
//	con, err := net.DialUDP(
//		"udp",
//		nil,
//		&net.UDPAddr{
//			//IP:   net.IPv4(172, 31, 226, 255),
//			IP: net.IPv4(172, 30, 98, 255),
//
//			Port: 3000,
//		},
//	)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer func(con *net.UDPConn) {
//		err := con.Close()
//		if err != nil {
//			fmt.Println(err)
//			return
//		}
//	}(con)
//	for {
//		select {
//		// udpSendChan有消息则取出,并向websocket发送
//		case data := <-udpSendChan:
//			_, err := con.Write(data)
//			if err != nil {
//				fmt.Println(err)
//				return
//			}
//		}
//	}
//}
//
//func udpRecvProc() {
//	con, err := net.ListenUDP("udp",
//		&net.UDPAddr{
//			IP:   net.IPv4zero,
//			Port: 3000,
//		})
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer func(con *net.UDPConn) {
//		err := con.Close()
//		if err != nil {
//			fmt.Println(err)
//			return
//		}
//	}(con)
//	for {
//		var buf [512]byte
//		n, err := con.Read(buf[0:])
//		if err != nil {
//			fmt.Println(err)
//			return
//		}
//		dispatch(buf[0:n])
//	}
//}

// 将收到的数据解包json化
func dispatch(data []byte) {
	msg := MessageBasic{}
	err := json.Unmarshal(data, &msg) // 解析成MessageBasic，前端根据Media字段认消息的类型，如文本、表情等
	if err != nil {
		fmt.Println(err)
		return
	}
	// 根据数据的聊天类型，决定向谁发送
	switch msg.Type {
	case 1: // 私聊
		targetUserId := msg.TargetId
		sendDirectMessage(targetUserId, data)
	case 2:
		targetGroupId := msg.TargetId
		fromId := msg.UserId
		sendGroupMsg(fromId, targetGroupId, data)
	}
}

// 点对点私聊
func sendDirectMessage(targetId string, msg []byte) {
	rwLocker.RLock()
	node, ok := clientMap[targetId]
	rwLocker.RUnlock()
	if ok {
		node.DataQueue <- msg
	}
}

func sendGroupMsg(fromId, targetId string, msg []byte) {
	// 查找群内用户
	umbList := FindMembersByGroupId(targetId)
	// 向这些群内用户发送数据
	//fmt.Println(string(msg))
	//fmt.Println(len(ubList))
	for _, umb := range umbList {
		if fromId == umb.MemberIdentity {
			continue
		}
		rwLocker.RLock()
		node, ok := clientMap[umb.MemberIdentity]
		rwLocker.RUnlock()
		if ok {
			node.DataQueue <- msg
		}
	}
}
