package tcp

import (
	"fmt"
	"io"
	"mschon/dbproxy/errors"
	"net"
	"strconv"
)

type ProxyRecDirection int

const (
	ProxyRecDirection_toServer ProxyRecDirection = iota
	ProxyRecDirection_fromServer
)

func (direction ProxyRecDirection) String() string {
	switch direction {
	case ProxyRecDirection_fromServer:
		return "S>C"
	case ProxyRecDirection_toServer:
		return "C>S"
	default:
		return strconv.Itoa(int(direction))
	}
}

type ProxyRecMsg struct {
	Direction  ProxyRecDirection
	Connection ProxyConnection

	Data []byte
}

func (msg ProxyRecMsg) String() string {
	return fmt.Sprintf("[ProxyRecMsg %s %d [%d]byte]", msg.Direction.String(), msg.Connection, len(msg.Data))
}

type ProxyConnection struct {
	Client  net.Conn
	Forward net.Conn

	quitChan chan bool
}

type Proxy struct {
	Connections []ProxyConnection

	ServerRecChan chan ProxyRecMsg
	ClientRecChan chan ProxyRecMsg
}

func NewProxy() Proxy {
	return Proxy{
		Connections:   make([]ProxyConnection, 0),
		ServerRecChan: make(chan ProxyRecMsg),
		ClientRecChan: make(chan ProxyRecMsg),
	}
}

func (proxy *Proxy) Listen(listenAddress string, forwardAddress string) {

	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		errors.WriteError(err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept error: ", err)
		}

		fwd, err := net.Dial("tcp", forwardAddress)

		if err != nil {
			fmt.Println("Dial error: ", err)
		}

		pair := ProxyConnection{
			Client:   conn,
			Forward:  fwd,
			quitChan: make(chan bool),
		}

		proxy.Connections = append(proxy.Connections, pair)

		go proxy.proxyHandler(&pair, ProxyRecDirection_toServer)
		go proxy.proxyHandler(&pair, ProxyRecDirection_fromServer)
	}
}

func (proxy *Proxy) proxyHandler(pair *ProxyConnection, direction ProxyRecDirection) {
	buffer := make([]byte, 1024)

	var local net.Conn = nil
	var forward net.Conn = nil
	var msgChan chan ProxyRecMsg = nil
	if direction == ProxyRecDirection_toServer {
		local = pair.Client
		forward = pair.Forward
		msgChan = proxy.ClientRecChan
	} else {
		local = pair.Forward
		forward = pair.Client
		msgChan = proxy.ServerRecChan
	}

	for {
		select {
		case <-pair.quitChan:
			return
		default:
			n, err := local.Read(buffer)
			if err != nil {
				if err != io.EOF {
					fmt.Println(" error: ", err)
				}
				return
			}

			forward.Write(buffer[:n])

			data := make([]byte, n)
			copy(data, buffer[:n])

			msg := ProxyRecMsg{
				Direction:  ProxyRecDirection_fromServer,
				Connection: *pair,
				Data:       data,
			}

			msgChan <- msg
		}
	}
}
