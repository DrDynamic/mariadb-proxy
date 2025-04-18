package mariadb

import (
	"mschon/dbproxy/tcp"
	"mschon/dbproxy/tcp/mariadb/packets"
)

type OnConnectionChangeListener func(MariadbConnection)
type OnNewConnectionListener func(MariadbConnection)
type OnNewServerPacketListener func(MariadbConnection)
type OnNewClientPacketListener func(MariadbConnection)

type MariadbConnectionManager struct {
	Connections []*MariadbConnection

	proxy                       *tcp.Proxy
	onConnectionChangeListeners []OnConnectionChangeListener
	onNewConnectionListeners    []OnNewConnectionListener
	onNewServerPacketListeners  []OnNewServerPacketListener
	onNewClientPacketListeners  []OnNewClientPacketListener
	quitSignal                  chan bool
}

type MariadbConnection struct {
	ServerPackets []packets.BasePacket
	ClientPackets []packets.BasePacket

	connection          *tcp.ProxyConnection
	serverPacketFactory ServerPacketFactory
	clientPacketFactory ClientPacketFactory
	quitSignal          chan bool
}

func (connection MariadbConnection) GetProxyConnection() *tcp.ProxyConnection {
	return connection.connection
}

func NewConnectionManager(proxy *tcp.Proxy) MariadbConnectionManager {
	return MariadbConnectionManager{
		proxy:                       proxy,
		Connections:                 make([]*MariadbConnection, 0),
		onConnectionChangeListeners: make([]OnConnectionChangeListener, 0),
		onNewConnectionListeners:    make([]OnNewConnectionListener, 0),
		onNewServerPacketListeners:  make([]OnNewServerPacketListener, 0),
		onNewClientPacketListeners:  make([]OnNewClientPacketListener, 0),
		quitSignal:                  make(chan bool),
	}
}

func (manager *MariadbConnectionManager) Listen(listenAddress string, forwardAddress string) {
	go manager.handleProxy()
	manager.proxy.Listen(listenAddress, forwardAddress)
}

func (manager *MariadbConnectionManager) AddOnConnectionChangeListener(listener OnConnectionChangeListener) {
	manager.onConnectionChangeListeners = append(manager.onConnectionChangeListeners, listener)
}

func (manager *MariadbConnectionManager) AddOnNewConnectionListener(listener OnNewConnectionListener) {
	manager.onNewConnectionListeners = append(manager.onNewConnectionListeners, listener)
}

func (manager *MariadbConnectionManager) AddOnNewServerPacketListener(listener OnNewServerPacketListener) {
	manager.onNewServerPacketListeners = append(manager.onNewServerPacketListeners, listener)
}

func (manager *MariadbConnectionManager) AddOnNewClientPacketListener(listener OnNewClientPacketListener) {
	manager.onNewClientPacketListeners = append(manager.onNewClientPacketListeners, listener)
}

func (manager *MariadbConnectionManager) handleProxy() {
	for {
		select {
		case <-manager.quitSignal:
			return
		case connection := <-manager.proxy.ConnectChan:
			facrotyState := NewPacketFactoryState()
			mariadbConnection := MariadbConnection{
				connection:          connection,
				serverPacketFactory: NewServerPacketFactory(&facrotyState),
				clientPacketFactory: NewClientPacketFactory(&facrotyState),
				ServerPackets:       make([]packets.BasePacket, 0),
				ClientPackets:       make([]packets.BasePacket, 0),
				quitSignal:          make(chan bool),
			}
			manager.Connections = append(manager.Connections, &mariadbConnection)
			go manager.handleConnection(&mariadbConnection)

			for _, listener := range manager.onNewConnectionListeners {
				listener(mariadbConnection)
			}
		}
	}
}

func (manager *MariadbConnectionManager) handleConnection(connection *MariadbConnection) {
	for {
		select {
		case <-connection.quitSignal:
			return
		case <-connection.connection.StatusChangeChan:
			for _, listener := range manager.onConnectionChangeListeners {
				listener(*connection)
			}
		case msg := <-connection.connection.ServerRecChan:
			connection.serverPacketFactory.AddBytes(msg.Data)
			var packet packets.BasePacket = connection.serverPacketFactory.CreatePacket()

			for packet != nil {
				connection.ServerPackets = append(connection.ServerPackets, packet)

				packet = connection.serverPacketFactory.CreatePacket()
			}
			for _, listener := range manager.onNewServerPacketListeners {
				listener(*connection)
			}
		case msg := <-connection.connection.ClientRecChan:
			connection.clientPacketFactory.AddBytes(msg.Data)
			var packet packets.BasePacket = connection.clientPacketFactory.CreatePacket()

			for packet != nil {
				connection.ClientPackets = append(connection.ClientPackets, packet)

				packet = connection.clientPacketFactory.CreatePacket()
			}
			for _, listener := range manager.onNewClientPacketListeners {
				listener(*connection)
			}
		}
	}
}
