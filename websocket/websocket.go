package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/dgrr/fastws"
	"github.com/valyala/fasthttp"
)

// WebSocket represents websocket connection handler.
type WebSocket struct {
	token          string
	userType       string
	instanceServer []instanceServer
	historyServer  []historyServer
}

// NewWS returns initilised websocket connection.
func NewWS() (*WebSocket, error) {
	ws := &WebSocket{}
	return ws, ws.init()
}

// SetUserType sets user type. Can be vip or normal.
//
// By default userType is normal.
//
// Use this function before calling Dial.
func (ws *WebSocket) SetUserType(userType string) {
	ws.userType = userType
}

func (ws *WebSocket) init() error {
	_, body, err := fasthttp.Get(nil, urlServers)
	if err != nil {
		return err
	}
	// Unmarshal server response
	var res wsResp
	err = json.Unmarshal(body, &res)
	if err != nil {
		return err
	}
	if !res.Success {
		return errors.New(res.Msg)
	}

	var data connData
	// Unmarshal response data
	err = json.Unmarshal(res.Data, &data)
	if err != nil {
		return err
	}

	ws.token = data.BulletToken
	ws.instanceServer = append(ws.instanceServer[:0], data.InstanceServers...)
	ws.historyServer = append(ws.historyServer[:0], data.HistoryServers...)

	return nil
}

func (ws *WebSocket) selectServers() (is *instanceServer, hs *historyServer) {
	for _, s := range ws.instanceServer {
		if ws.userType == "" || s.UserType == ws.userType {
			is = &s
			break
		}
	}
	for _, s := range ws.historyServer {
		if ws.userType == "" || s.UserType == ws.userType {
			hs = &s
			break
		}
	}
	return
}

// Subscribe subscribes client to a Topic (Orderbook level2, History, Tick, Market)
func (ws *WebSocket) Subscribe(topic Topic, symbol string) (c *Conn, err error) {
	conn, is, err := ws.dial()
	if err == nil {
		c = &Conn{
			cond:           sync.NewCond(&sync.Mutex{}),
			symbol:         symbol,
			instanceServer: is,
			topic:          topic,
			conn:           conn,
		}
		_, err = c.Send(Subscribe, topic, symbol)
		if err == nil {
			c.init()
			go c.handle()
		} else {
			conn.Close(err.Error())
			c = nil
		}
	}
	return
}

func (ws *WebSocket) dial() (conn *fastws.Conn, is instanceServer, err error) {
	iss, _ := ws.selectServers()
	if iss == nil {
		err = fmt.Errorf("error selecting server. Server for %s not found", ws.userType)
		return
	}
	is = *iss

	uri := fasthttp.AcquireURI()
	defer fasthttp.ReleaseURI(uri)

	uri.Update(is.Endpoint)
	args := uri.QueryArgs()
	args.Add("bulletToken", ws.token)
	args.Add("format", "json")
	args.Add("resource", "api")

	conn, err = fastws.Dial(uri.String())
	return
}
