package goio

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
)

type indexedWSConn struct {
	Index      int64
	SignalChan chan struct{}
	ExitChan   chan struct{}
}

func NewIndexedWSConn(index int64) *indexedWSConn {
	return &indexedWSConn{
		Index:      index,
		SignalChan: make(chan struct{}),
		ExitChan:   make(chan struct{}),
	}
}

func HandleSignalWs(key string, signalChan chan struct{}, exitChan chan struct{}) func(*websocket.Conn) {
	var wskey string = key
	var counter int64 = 0
	connections := make(map[int64]*indexedWSConn)
	addWSChan := make(chan *indexedWSConn)
	removeWSChan := make(chan *indexedWSConn)
	go func() {
		for {
			select {
			case <-signalChan: // Send a signal to all clients
				for _, conn := range connections {
					conn.SignalChan <- struct{}{}
				}
			case conn := <-addWSChan: // Pur a new client in the pool
				connections[conn.Index] = conn
			case conn := <-removeWSChan: // Take index out of client list
				close(conn.ExitChan)
				delete(connections, conn.Index)
			case <-exitChan: // Close the handler
				return
			}
		}
	}()
	return func(conn *websocket.Conn) {
		defer func() { conn.Close() }()
		counter++
		var index = counter
		log := func(message string) {
			fmt.Printf("\n[%d - %s]\t%s\n", index, wskey, message)
		}
		log("Got ws connection")
		connEntry := NewIndexedWSConn(index)
		addWSChan <- connEntry

		go func() {
			for {
				var data struct {
					Key   string
					Value string
				}
				if err := websocket.JSON.Receive(conn, &data); err != nil {
					log(fmt.Sprintf("Error receiving ws: %s", err))
					removeWSChan <- connEntry
					return
				}
				log(fmt.Sprintf("Got data: %v\n", data))
			}
		}()

		go func() {
			for {
				select {
				case <-connEntry.SignalChan:
					var data struct{}
					if err := websocket.JSON.Send(conn, &data); err != nil {
						log(fmt.Sprintf("Error sending on socket: %s", err))
						removeWSChan <- connEntry
					}
					log("Sent signal")
				case <-connEntry.ExitChan:
					log("Sending terminated")
					return
				case <-exitChan:
					log("Exiting send routine")
					return
				}
			}
		}()
		select {
		case <-connEntry.ExitChan: // The ws connection terminated
		case <-exitChan: // The server was terminated
		}
		log("Connection terminated")
	}
}
