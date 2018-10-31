package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/fiore/kucoin-go/websocket"
)

func main() {
	// Starting websocket
	ws, err := websocket.NewWS()
	if err != nil {
		log.Fatalln(err)
	}

	// Creating History, OrderBook and Market conns
	ch, co, cm := startConns(ws, "ETH", "BTC")

	// Handling SIGINT
	sch := make(chan os.Signal, 1)
	signal.Notify(sch, os.Interrupt)

	var update interface{}
loop:
	for {
		update = nil
		// Getting updates
		select {
		case update = <-ch.Updates():
		case update = <-co.Updates():
		case update = <-cm.Updates():
		case <-sch:
			log.Println("SIGINT")
			signal.Stop(sch)
			signal.Reset()
			break loop
		}
		if update != nil {
			switch up := update.(type) {
			case error:
				err = up
			case *websocket.History:
				log.Println("History:", up)
			case *websocket.OrderBook:
				log.Println("OrderBook:", up)
			case *websocket.Market:
				log.Println("Market:", up)
			}
		}
		if err != nil {
			log.Fatalln(err)
		}
	}
	// closing signal channel
	close(sch)

	// Closing connections
	ch.Close()
	co.Close()
	cm.Close()
}

func startConns(ws *websocket.WebSocket, c1, c2 string) (ch, co, cm *websocket.Conn) {
	var err error
	ch, err = ws.Subscribe(websocket.THistory, c1+"-"+c2)
	if err != nil {
		log.Fatalln("history:", err)
	}
	co, err = ws.Subscribe(websocket.TOrderBook, c1+"-"+c2)
	if err != nil {
		log.Fatalln("order book:", err)
	}
	cm, err = ws.Subscribe(websocket.TMarket, c2)
	if err != nil {
		log.Fatalln("market:", err)
	}
	return
}
