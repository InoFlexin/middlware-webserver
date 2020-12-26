package base

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"sync"
)

type Boot struct {
	Protocol    string //server protocol type (tcp/ip, udp... etc)
	Port        string
	ServerName  string
	Callback    SocketEvent
	ReceiveSize uint64
	Complex     bool
}

type Message struct {
	Json   string
	Action int
}

type SocketEvent interface {
	OnMessageReceive(message *Message, client net.Conn)
	OnConnect(message *Message, client net.Conn)
	OnClose(err error)
}

func Write(message *Message, client net.Conn) {
	//Encoding Message to Json String
	e, err := json.Marshal(message)

	if err != nil {
		log.Println(err)
		return
	}

	client.Write(e)
}

func PacketUnmarshal(data []byte) *Message {
	message := Message{}
	json.Unmarshal([]byte(data), &message)

	return &message
}

/*
	Please running goroutine
*/
func Broadcast(message *Message) {
	keys, values := GetSessions()
	length := len(keys)

	for i := 0; i < length; i++ {
		conn := values[i]

		Write(message, conn)
	}
}

func _receiveAndHandle(buf []byte, connection net.Conn, boot *Boot, count int) {
	var data = buf[:count]
	message := PacketUnmarshal(data)

	switch message.Action {
	case ON_MSG_RECEIVE:
		go boot.Callback.OnMessageReceive(message, connection)
		break
	case ON_CONNECT:
		go boot.Callback.OnConnect(message, connection)
		break
	}
}

func receive(connection net.Conn, boot Boot) {
	buf := make([]byte, boot.ReceiveSize) //1kb

	for {
		count, error := connection.Read(buf)

		if nil != error {
			if io.EOF == error {
				log.Printf("connection is closed from client; %v", connection.RemoteAddr().String())
				RemoveSession(connection.RemoteAddr().String()) // if client connection refused, remove session.
				go boot.Callback.OnClose(nil)
			}

			log.Printf("fail to receive data; err: %v", error)
			return
		}

		if count > 0 {
			_receiveAndHandle(buf, connection, &boot, count)
		}
	}
}

func SetupComplexServer(boot Boot, wg *sync.WaitGroup) {
	if boot.Complex {
		wg.Done()
	}
}

func ServerStart(boot Boot, wg *sync.WaitGroup) {
	listener, error := net.Listen(boot.Protocol, boot.Port)
	log.Println(boot)
	log.Println(boot.ServerName + " get started port: " + boot.Port)

	if error != nil {
		log.Fatalf("Failed to bind address to "+boot.Port+" err: %v", error)
	}

	defer listener.Close()
	defer wg.Done()
	SetupComplexServer(boot, wg)

	for {
		conn, error := listener.Accept()

		if nil != error {
			log.Printf("Failed to accept; err: %v", error)
			continue
		}

		id, sessionError := AddSession(conn.LocalAddr().String(), conn)

		if sessionError == nil {
			go receive(conn, boot)
			log.Printf("Successfully session added id: " + id)
		} else {
			log.Fatal(sessionError)
		}
	}
}
