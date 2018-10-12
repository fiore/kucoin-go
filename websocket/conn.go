package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/dgrr/fastws"
)

// Conn represents WebSocket connection to a Topic.
//
// Use Update() func to get server messages.
type Conn struct {
	sym string
	ps  instanceServer
	up  chan interface{}
	tc  Topic
	sm  string
	c   *fastws.Conn
}

// Symbol returns current connected symbol.
func (c *Conn) Symbol() string {
	return c.sym
}

// PingInterval returns server ping interval.
func (c *Conn) PingInterval() int {
	return c.ps.PingInterval
}

// Encrypt returns if connection is encrypted (wss or ws).
func (c *Conn) Encrypt() bool {
	return c.ps.Encrypt
}

// PingTimeout returns server ping timeout.
func (c *Conn) PingTimeout() int {
	return c.ps.PingTimeout
}

// UserType returns connection user type.
func (c *Conn) UserType() string {
	return c.ps.UserType
}

func (c *Conn) init() {
	c.close()
	c.up = make(chan interface{}, 10)
}

func (c *Conn) close() {
	if c.up != nil {
		close(c.up)
		c.up = nil
	}
}

// Updates is the notification channel.
//
// The types which can be sended through channel are:
// error, History, OrderBook and Market
func (c *Conn) Updates() <-chan interface{} {
	return c.up
}

// IsClosed returns if connection is closed.
func (c *Conn) IsClosed() bool {
	return c.up == nil && c.c == nil
}

// Close closes websocket connection and updates channel.
func (c *Conn) Close() (err error) {
	err = c.c.Close("Bye")
	if err == nil {
		c.close()
		c.c = nil
	}
	return err
}

var nid uint64 = 7

func nextId() uint64 {
	return atomic.AddUint64(&nid, 2)
}

// Send sends actions to perform
func (c *Conn) Send(tp Type, tc Topic, sym string) (r Response, err error) {
	var url string
	switch tc {
	case TOrderBook:
		url = urlOrderBook
	case THistory:
		url = urlHistory
	case Tick:
		url = urlTick
	case TMarket:
		url = urlMarket
	default:
		err = fmt.Errorf("invalid topic: %d", tc)
		return
	}
	url = fmt.Sprintf(url, sym)

	req := wsReq{
		Id:    nextId(),
		Type:  string(tp),
		Topic: url,
		Req:   0, // TODO
	}
	var data []byte
	data, err = json.Marshal(req)
	if err != nil {
		c.c.Close("error marshaling data")
		return
	}
	_, err = c.c.Write(data)
	if err == nil {
		var fr *fastws.Frame
		fr, err = c.c.NextFrame() // must read ack
		if err == nil {
			err = json.Unmarshal(fr.Payload(), &r)
			fastws.ReleaseFrame(fr)
		}
	}
	return
}

func (c *Conn) handle() {
	if c.c == nil {
		c.up <- errors.New("nil connection")
		return
	}
	c.init()

	var fr *fastws.Frame
	var err error
	for {
		fr, err = c.c.NextFrame()
		if err != nil {
			if err != fastws.EOF {
				c.up <- err
			}
			break
		}

		err = c.handlePingClose(fr)
		if err != nil {
			if err != fastws.EOF {
				c.up <- err
			}
			break
		}
		if c.up != nil {
			c.up <- doDecode(c.tc, fr.Payload())
		}
		fastws.ReleaseFrame(fr)
	}
}

func doDecode(tc Topic, b []byte) interface{} {
	var res wsResp
	var dst interface{}

	err := json.Unmarshal(b, &res)
	if err != nil {
		return err
	}

	switch tc {
	case TOrderBook:
		dst = new(OrderBook)
	case THistory:
		dst = new(History)
	case TMarket, Tick:
		dst = new(Market)
	default:
		return errors.New("topic not valid")
	}

	err = json.Unmarshal(res.Data, dst)
	if err != nil {
		return err
	}
	return dst
}

func (c *Conn) handlePingClose(fr *fastws.Frame) (err error) {
	switch {
	case fr.IsPing():
		fr.Reset()
		fr.SetFin()
		fr.SetPong()
		_, err = c.c.WriteFrame(fr)
	case fr.IsClose():
		err = c.c.SendCode(fr.Code(), fr.Status(), nil)
		if err == nil {
			err = fastws.EOF
		}
	}
	return
}
