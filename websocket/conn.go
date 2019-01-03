package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/dgrr/fastws"
)

// Conn represents WebSocket connection to a Topic.
// Use Update() func to get server messages.
type Conn struct {
	cn *sync.Cond
	b  bool

	sym string
	ps  instanceServer

	up      chan interface{}
	noClose bool

	tc         Topic
	sm         string
	c          *fastws.Conn
	lastUpdate time.Time
}

func (c *Conn) lock() {
	c.cn.L.Lock()
	for c.b {
		c.cn.Wait()
	}
	c.b = true
	c.cn.L.Unlock()
}

func (c *Conn) unlock() {
	c.b = false
	c.cn.Signal()
}

// SetUpdates sets parsed channel to be used when a update fire.
func (c *Conn) SetUpdates(ch chan interface{}) {
	c.close()

	c.lock()
	c.up = ch
	c.noClose = true
	c.unlock()
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
	c.lastUpdate = time.Now()
}

func (c *Conn) close() {
	if !c.noClose && c.up != nil {
		c.lock()
		close(c.up)
		c.up = nil
		c.unlock()
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
	if c.c != nil {
		err = c.c.Close()
		if err == nil {
			c.c = nil
		}
	}
	c.close()
	return err
}

var nid uint64 = 7

func nextID() uint64 {
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
		Id:    nextID(),
		Type:  string(tp),
		Topic: url,
		Req:   0, // TODO
	}
	var data []byte
	data, err = json.Marshal(req)
	if err != nil {
		c.c.Close()
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

func (c *Conn) checkUpdates(stop chan struct{}) {
	interval := time.Duration(c.PingInterval())
	for {
		select {
		case <-stop:
			return
		case <-time.After(time.Millisecond * interval):
			if time.Now().Sub(c.lastUpdate) > interval {
				c.lastUpdate = time.Now()
				resp := pingReq{
					Id:   nextID(),
					Type: "ping",
				}
				data, err := json.Marshal(resp)
				if err != nil {
					c.sendUpdate(err)
				} else {
					c.c.Write(data)
				}
			}
		}
	}
}

func (c *Conn) handle() {
	if c.c == nil {
		c.up <- errors.New("nil connection")
		return
	}
	stop := make(chan struct{}, 1)
	go c.checkUpdates(stop)

	var fr *fastws.Frame
	var err error
	for {
		fr, err = c.c.NextFrame()
		if err != nil {
			if err == fastws.EOF {
				break
			}
			c.sendUpdate(err)
			continue
		}

		err = c.handlePingClose(fr)
		if err != nil {
			if err != fastws.EOF {
				c.sendUpdate(err)
			}
			break
		}
		res := c.doDecode(c.tc, fr.Payload())
		if c.sendUpdate(res) {
			break
		}

		fastws.ReleaseFrame(fr)
	}
	stop <- struct{}{}
}

func (c *Conn) sendUpdate(res interface{}) (brk bool) {
	c.lock()
	defer c.unlock()
	defer func() {
		if recover() != nil {
			brk = true
		}
	}()
	c.up <- res
	return
}

func (c *Conn) doDecode(tc Topic, b []byte) interface{} {
	var res wsResp
	var dst interface{}

	err := json.Unmarshal(b, &res)
	if err != nil {
		return err
	}

	switch res.Type {
	case "ack":
		return nil
	case "pong":
		return nil
	}

	switch res.Code.String() {
	case "404":
		err = fmt.Errorf("%s: %s", res.Code, res.Data)
	}
	if err != nil {
		return err
	}

	switch tc {
	case TOrderBook:
		dst = &OrderBook{
			Symbol: c.Symbol(),
		}
	case THistory:
		dst = &History{
			Symbol: c.Symbol(),
		}
	case TMarket, Tick:
		dst = new(Market)
	default:
		return errors.New("topic not valid")
	}

	err = json.Unmarshal(res.Data, dst)
	if err != nil {
		dst = err
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
		err = c.c.ReplyClose(fr)
		if err == nil {
			err = fastws.EOF
		}
	}
	return
}
