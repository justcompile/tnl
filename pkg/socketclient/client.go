package socketclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/justcompile/tnl/pkg/socketserver"
	"github.com/justcompile/tnl/pkg/types"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	conn          *websocket.Conn
	serverAddress string
}

func (c *Client) Close() error {
	err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("write close:", err)
		return err
	}

	return c.conn.Close()
}

func (c *Client) Connect() {
	u := url.URL{Scheme: "ws", Host: c.serverAddress, Path: socketserver.WebSocketPath}
	log.Printf("connecting to %s", u.String())

	var err error

	headers := http.Header{}
	headers.Add(socketserver.DomainHeader, "foobar.com")

	c.conn, _, err = websocket.DefaultDialer.Dial(u.String(), headers)
	if err != nil {
		log.Fatal("dial:", err)
	}

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			log.Println(".")

			// log.Printf("recv: %s", message)

			resp, err := makeRequest(message)
			if err != nil {
				log.Println("req:", err)
				return
			}

			if err := c.conn.WriteMessage(websocket.BinaryMessage, resp); err != nil {
				log.Println("send:", err)
			}
		}
	}()
	//
	// ticker := time.NewTicker(time.Second)
	// defer ticker.Stop()
	//
	// for {
	// 	select {
	// 	case <-done:
	// 		return
	// 	case t := <-ticker.C:
	// 		err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
	// 		if err != nil {
	// 			log.Println("write:", err)
	// 			return
	// 		}
	// 	case <-interrupt:
	// 		log.Println("interrupt")
	//
	// 		// Cleanly close the connection by sending a close message and then
	// 		// waiting (with timeout) for the server to close the connection.
	// 		err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	// 		if err != nil {
	// 			log.Println("write close:", err)
	// 			return
	// 		}
	// 		select {
	// 		case <-done:
	// 		case <-time.After(time.Second):
	// 		}
	// 		return
	// 	}
	// }
}

func makeRequest(in []byte) ([]byte, error) {
	var r *types.Request

	if err := json.Unmarshal(in, &r); err != nil {
		return nil, err
	}

	forwardURL := fmt.Sprintf("http://localhost:3333%s", r.URL)

	var req *http.Request
	if len(r.Body) > 0 {
		req, _ = http.NewRequest(r.Method, forwardURL, bytes.NewReader(r.Body))
	} else {
		req, _ = http.NewRequest(r.Method, forwardURL, nil)
	}

	for key, val := range r.Headers {
		req.Header.Set(key, val)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return types.SerializeResponse(resp)
}

func New(addr string) *Client {
	return &Client{
		serverAddress: addr,
	}
}
