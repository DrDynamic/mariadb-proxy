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

	RecChan chan ProxyRecMsg
}

func NewProxy() Proxy {
	return Proxy{
		Connections: make([]ProxyConnection, 0),
		RecChan:     make(chan ProxyRecMsg),
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

		go proxyHandler(pair.Client, pair.Forward, pair.quitChan, func(data []byte) {
			proxy.RecChan <- ProxyRecMsg{
				Direction:  ProxyRecDirection_toServer,
				Connection: pair,
				Data:       data,
			}
		})
		go proxyHandler(pair.Forward, pair.Client, pair.quitChan, func(data []byte) {
			proxy.RecChan <- ProxyRecMsg{
				Direction:  ProxyRecDirection_fromServer,
				Connection: pair,
				Data:       data,
			}
		})
	}
}

func proxyHandler(local net.Conn, forward net.Conn, quitChan chan bool, f func([]byte)) {
	buffer := make([]byte, 1024)

	for {
		select {
		case <-quitChan:
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
			f(buffer[:n])
		}
	}
}
