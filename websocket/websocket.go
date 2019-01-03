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
	token    string
	userType string
	ps       []instanceServer
	hs       []historyServer
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
	ws.ps = append(ws.ps[:0], data.InstanceServers...)
	ws.hs = append(ws.hs[:0], data.HistoryServers...)

	return nil
}

func (ws *WebSocket) selectServers() (ps *instanceServer, hs *historyServer) {
	for _, s := range ws.ps {
		if ws.userType == "" || s.UserType == ws.userType {
			ps = &s
			break
		}
	}
	for _, s := range ws.hs {
		if ws.userType == "" || s.UserType == ws.userType {
			hs = &s
			break
		}
	}
	return
}

// Subscribe subscribes client to a topic.
func (ws *WebSocket) Subscribe(tc Topic, sym string) (c *Conn, err error) {
	var conn *fastws.Conn
	var ps instanceServer
	conn, ps, err = ws.dial(Subscribe, tc, sym)
	if err == nil {
		c = &Conn{
			cn:  sync.NewCond(&sync.Mutex{}),
			sym: sym,
			ps:  ps,
			sm:  sym,
			tc:  tc,
			c:   conn,
		}
		_, err = c.Send(Subscribe, tc, sym)
		if err == nil {
			c.init()
			go c.handle()
		} else {
			conn.Close()
			c = nil
		}
	}
	return
}

func (ws *WebSocket) dial(t Type, tc Topic, sym string) (c *fastws.Conn, ps instanceServer, err error) {
	pps, _ := ws.selectServers()
	if pps == nil {
		err = fmt.Errorf("error selecting server. Server for %s not found", ws.userType)
		return
	}
	ps = *pps

	uri := fasthttp.AcquireURI()
	defer fasthttp.ReleaseURI(uri)

	uri.Update(ps.Endpoint)
	args := uri.QueryArgs()
	args.Add("bulletToken", ws.token)
	args.Add("format", "json")
	args.Add("resource", "api")

	c, err = fastws.Dial(uri.String())
	return
}
