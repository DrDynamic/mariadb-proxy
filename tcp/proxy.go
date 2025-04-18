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
	return fmt.Sprintf("[ProxyRecMsg %s %v [%d]byte]", msg.Direction.String(), msg.Connection, len(msg.Data))
}

type ConnectionStatus string

const (
	ConnectionStatus_Unknown = "unknown"
	ConnectionStatus_Active  = "active"
	ConnectionStatus_Closed  = "closed"
)

type ProxyConnection struct {
	Client  net.Conn
	Forward net.Conn

	Status ConnectionStatus

	ServerRecChan    chan ProxyRecMsg
	ClientRecChan    chan ProxyRecMsg
	StatusChangeChan chan ConnectionStatus
	quitChan         chan bool
}

type Proxy struct {
	Connections []*ProxyConnection
	ConnectChan chan *ProxyConnection
}

func NewProxy() Proxy {
	return Proxy{
		Connections: make([]*ProxyConnection, 0),
		ConnectChan: make(chan *ProxyConnection),
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
			Client:           conn,
			Forward:          fwd,
			ServerRecChan:    make(chan ProxyRecMsg),
			ClientRecChan:    make(chan ProxyRecMsg),
			StatusChangeChan: make(chan ConnectionStatus),
			Status:           ConnectionStatus_Active,
			quitChan:         make(chan bool),
		}

		proxy.Connections = append(proxy.Connections, &pair)

		proxy.ConnectChan <- &pair

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
		msgChan = pair.ClientRecChan
	} else {
		local = pair.Forward
		forward = pair.Client
		msgChan = pair.ServerRecChan
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
				pair.Status = ConnectionStatus_Closed
				pair.StatusChangeChan <- pair.Status
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
