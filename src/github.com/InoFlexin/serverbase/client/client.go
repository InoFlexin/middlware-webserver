package client

import (
	"errors"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/InoFlexin/serverbase/base"
)

type ClientBoot struct {
	Protocol   string
	HostAddr   string
	HostPort   string
	Callback   base.SocketEvent
	BufferSize uint64
}

var server net.Conn = nil

func CreateError(errorMessage string) error {
	return errors.New(errorMessage)
}

func SendPing(duration time.Duration) error {
	var serverError error = nil

	if server != nil {
		for {
			base.Write(&base.Message{Json: "ping", Action: base.ON_MSG_RECEIVE}, server)
			time.Sleep(duration)
		}
	} else {
		serverError = CreateError("Server not connected error")
	}

	return serverError
}

func Handle(boot *ClientBoot, conn net.Conn, wg *sync.WaitGroup) {
	buf := make([]byte, boot.BufferSize)
	defer wg.Done() //핸들 처리용 wait group은 밑에 logic이 끝나면 Done 시킨다.

	for {
		count, error := conn.Read(buf)

		if nil != error {
			if io.EOF == error {
				log.Printf("connection is closed from server; %v", conn.RemoteAddr().String())
			}

			log.Printf("fail to receive data; err: %v", error)
			return
		}

		if count > 0 {
			data := buf[:count]
			boot.Callback.OnMessageReceive(base.PacketUnmarshal(data), conn)
		}
	}
}

func ConnectServer(boot *ClientBoot, wg *sync.WaitGroup) {
	log.Println("Run")
	conn, err := net.Dial(boot.Protocol, boot.HostAddr+boot.HostPort)
	server = conn
	boot.Callback.OnConnect(&base.Message{Json: "-connect", Action: base.ON_CONNECT}, conn)

	if err != nil {
		log.Fatalf("faild to connec to server err %v", err)
	}

	wg.Done()                    //main에서 처리한 waitGroup은 done() 시킨다.
	handleWg := sync.WaitGroup{} //고루틴 핸들러 처리용 wait group
	handleWg.Add(1)
	go Handle(boot, server, &handleWg)
	handleWg.Wait()

	defer conn.Close()
}
