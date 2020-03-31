package logic

import (
	"awesomeProject/application/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/micro/go-micro/v2/errors"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	DefaultAddress = ":7272"
)

type (
	ImServer struct {
		rabbitMqBroker *RabbitMqBroker
		clients        map[string]*websocket.Conn
		Address        string
		lock           sync.Mutex
		upgraer        *websocket.Upgrader
	}
	SendMsgRequest struct {
		FromToken     string `json:"fromToken"`
		ToToken       string `json:"toToken"`
		Body          string `json:"body"`
		TimeStamp     int64  `json:"timeStamp"`
		RemoteAddress string `json:"remoteAddress"`
	}

	LoginRequest struct {
		Token string `json:"token"`
	}
	SendMsgResponse struct {
		FromToken     string `json:"fromToken"`
		Body          string `json:"body"`
		RemoteAddress string `json:"remoteAddress"`
	}
	ImServerOptionsFunc func(im *ImServer)
)

func NewImServer(rabbitMqBroker *RabbitMqBroker, opts ...ImServerOptionsFunc) (*ImServer, error) {
	// 初始化
	imServer := &ImServer{
		rabbitMqBroker: rabbitMqBroker,
		clients:        make(map[string]*websocket.Conn, 0),
		upgraer: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}

	for _, opt := range opts {
		opt(imServer)
	}

	if imServer.Address == "" {
		imServer.Address = DefaultAddress
	}
	return imServer, nil
}

func (l *ImServer) Run() {
	http.HandleFunc("/ws", l.login)
	//go func(c map[string]*websocket.Conn){
	//	time.Sleep(time.Second*10)
	//	for {
	//		con := c["err"]
	//		_, message, _ := con.ReadMessage()
	//		fmt.Println(string(message))
	//		_ = con.WriteMessage(1, []byte("我特么收到消息了"))
	//	}
	//
	//}(l.clients)
	log.Fatal(http.ListenAndServe(l.Address, nil))
}

func (l *ImServer) SendMsg(r *SendMsgRequest) (*SendMsgResponse, error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	log.Printf("send SendMsgRequest  %+v", r)
	conn := l.clients[r.ToToken]
	if conn == nil {
		return nil, &errors.Error{
			Code:   utils.ErrConnTokenFailed,
			Detail: "用户连接不存在",
			Status: http.StatusText(500),
		}
	}
	r.TimeStamp = time.Now().Unix()
	r.RemoteAddress = conn.RemoteAddr().String()
	bodyMsg, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	if err := conn.WriteMessage(websocket.TextMessage, bodyMsg); err != nil {
		log.Printf("send message err %v", err)
		l.clients[r.ToToken] = nil
		log.Println(conn.Close())
		return nil, err
	}
	log.Printf("send message succes  %v", r.Body)
	return &SendMsgResponse{}, nil
}

func (l *ImServer) Subscribe() {
	err := l.rabbitMqBroker.Subscribe(func(topic string, body []byte) {
		senRes := &SendMsgRequest{
			//FromToken:     "",
			//ToToken:       "",
			//Body:          "",
			//TimeStamp:     1,
			//RemoteAddress: "",
		}
		err := json.Unmarshal(body,senRes)
		if err != nil {
			fmt.Println(err)
		}
		_, err = l.SendMsg(senRes)
		if err != nil {
			fmt.Println(err)
		}
	})
	fmt.Println(err)
}

func (l *ImServer) login(w http.ResponseWriter, r *http.Request) {
	//登录的时候会访问一次 但是这里面有个接受消息 所以就会阻塞在这
	//每次登录就会保存这个token和对应的conn
	conn, err := l.upgraer.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	msgType, message, err := conn.ReadMessage()
	if err != nil {
		log.Printf("read login message err %+v", err)
		return
	}
	if msgType != websocket.TextMessage {
		log.Printf("read login msgType err %+v", err)
		return
	}
	fmt.Println(string(message))
	loginMsgRequest := new(LoginRequest)

	if err := json.Unmarshal(message, loginMsgRequest); err != nil {
		log.Printf("json.Unmarshal msg err %+v", err)
		return
	}

	l.clients[loginMsgRequest.Token] = conn
	return
}

func ImServerAddress(adr string) ImServerOptionsFunc {
	return func(im *ImServer) {
		im.Address = adr
	}
}
