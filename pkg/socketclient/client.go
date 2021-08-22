package socketclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/justcompile/tnl/pkg/socketserver"
	"github.com/justcompile/tnl/pkg/types"
	"github.com/justcompile/tnl/pkg/ui"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	conn            *websocket.Conn
	localBinding    string
	remoteProtocol  string
	wsServerAddress string
}

func (c *Client) Close() error {
	err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("write close:", err)
		return err
	}

	return c.conn.Close()
}

func (c *Client) Connect(infoUpdates, requestUpdates chan interface{}) {
	u := url.URL{Scheme: "ws", Host: c.wsServerAddress, Path: socketserver.WebSocketPath}

	domain := c.generateSubdomain()

	var err error

	headers := http.Header{}
	headers.Add(socketserver.DomainHeader, domain)

	go func() {
		infoUpdates <- &ui.BannerText{
			Endpoint: fmt.Sprintf("%s://%s", c.remoteProtocol, domain),
			Port:     c.localBinding,
		}
	}()

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

			resp, err := c.makeRequest(requestUpdates, message)
			if err != nil {
				log.Println("req:", err)
				return
			}

			if err := c.conn.WriteMessage(websocket.BinaryMessage, resp); err != nil {
				log.Println("send:", err)
			}
		}
	}()
}

func (c *Client) generateSubdomain() string {
	sub := uuid.NewV4().String()

	baseDomain := strings.Split(c.wsServerAddress, ":")[0]

	return sub + "." + baseDomain
}

func (c *Client) makeRequest(requestChan chan interface{}, in []byte) ([]byte, error) {
	var r *types.Request

	if err := json.Unmarshal(in, &r); err != nil {
		return nil, err
	}

	start := time.Now()

	forwardURL := fmt.Sprintf("http://localhost:%s%s", c.localBinding, r.URL)

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
	requestDur := time.Since(start)
	if err != nil {
		requestChan <- []string{
			r.URL, fmt.Sprintf("[%s](fg:red)", err.Error()), requestDur.String(),
		}

		return nil, err
	}

	var colour string
	if resp.StatusCode < 400 {
		colour = "green"
	} else if resp.StatusCode < 500 {
		colour = "yellow"
	} else {
		colour = "red"
	}

	requestChan <- []string{
		r.URL, fmt.Sprintf("[%d](fg:%s)", resp.StatusCode, colour), requestDur.String(),
	}

	return types.SerializeResponse(resp)
}

func New(opts *Options) *Client {
	return &Client{
		localBinding:    opts.LocalBindAddress,
		remoteProtocol:  opts.Protocol,
		wsServerAddress: opts.WebsocketServerBindAddress,
	}
}
