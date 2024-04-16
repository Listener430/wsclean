package websocket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
)

func (nc *NotificationClient) BroadcastMessage(msg interface{}) error {
	message, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return nc.Client.WriteMessage(1, []byte(message))
}

// readPump pumps messages from the websocket connection to the hub.

func (cl *OnlineClient) ReadPump() {
	fmt.Print("ReadPump started\n")
	var onlineEventPayload OnlineEventStructure
	defer unRegisterAndCloseOnlineConnection(cl)
	for {
		fmt.Printf("\n ReadPump listening to read message\n")
		msgtype, payload, err := cl.WebSocketConnection.ReadMessage()
		fmt.Printf("\n ReadPump msgtype=%d \n, payload %s", msgtype, payload)

		if err != nil {
			fmt.Printf("\n Error in ReadPump %s", err.Error())
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				break
			}
			break
		}

		err = json.Unmarshal(payload, &onlineEventPayload)
		if err != nil {
			fmt.Printf("\n Error in ReadPump unmarshaling %s", err.Error())
			break
		}

	}
	fmt.Print("ReadPump finished\n")
}

// writePump pumps messages from the hub to the websocket connection.

func (cl *OnlineClient) WritePump() {
	fmt.Print("WritePump started\n")
	defer func() {
		err := cl.WebSocketConnection.Close()
		if err != nil {
			return
		}
	}()

	fmt.Print("WritePump - before inf loop \n")

	// An infinite loop to keep the connection alive
	for {
		select {
		case payload, ok := <-cl.Send:
			reqBodyBytes := new(bytes.Buffer)
			fmt.Print("WritePump - case payload \n")
			if err := json.NewEncoder(reqBodyBytes).Encode(payload); err != nil {
				//log.Error("write_pump", flog.NewField("json_encode", err))
				continue
			}

			finalPayload := reqBodyBytes.Bytes()
			if !ok {
				if err := cl.WebSocketConnection.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					fmt.Print("WritePump - issue with final payload\n")
				}
				return
			}
			w, err := cl.WebSocketConnection.NextWriter(websocket.TextMessage)
			fmt.Print("WritePump - getting ready to write next message \n")
			if err != nil {
				//log.Error("write_pump", flog.NewField("web_socket_connection", "next_writer"), flog.NewField("err", err))
				return
			}
			_, err = w.Write(finalPayload)
			if err != nil {
				return
			}
			fmt.Print("WritePump - written next message \n")

			n := len(cl.Send)
			for i := 0; i < n; i++ {
				fmt.Print("WritePump -sending more messages \n")

				if err := json.NewEncoder(reqBodyBytes).Encode(<-cl.Send); err != nil {
					//	log.Error("write_pump", flog.NewField("err", err))
					continue
				}
				if _, err := w.Write(reqBodyBytes.Bytes()); err != nil {
					//	log.Error("write_pump", flog.NewField("err", err))
					continue
				}
			}
			if err := w.Close(); err != nil {
				return
			}
			//default:
			// Default action here
			//fmt.Print(".")
		}
	}
}

func unRegisterAndCloseOnlineConnection(c *OnlineClient) {
	fmt.Printf("UnregisteringAndCloseOnlineConnection started\n")
	c.OnlineHub.Unregister <- c
	err := c.WebSocketConnection.Close()
	if err != nil {
		return
	}
}
