package main

import (
	"context"
	"github.com/youluo1230/okex"
	"github.com/youluo1230/okex/api"
	"github.com/youluo1230/okex/events"
	"github.com/youluo1230/okex/events/public"
	ws_public_requests "github.com/youluo1230/okex/requests/ws/public"
	"log"
)

func main() {
	apiKey := ""
	secretKey := ""
	passphrase := ""

	errChan := make(chan *events.Error)
	subChan := make(chan *events.Subscribe)
	uSubChan := make(chan *events.Unsubscribe)
	logChan := make(chan *events.Login)
	sucChan := make(chan *events.Success)
	tradesChan := make(chan *public.Trades)

	ctx := context.Background()
	client, _ := api.NewClient(ctx, apiKey, secretKey, passphrase, okex.CnServer)
	client.Ws.SetChannels(errChan, subChan, uSubChan, logChan, sucChan)
	err := client.Ws.Public.Trades(ws_public_requests.Trades{
		InstID: "BTC-USD-SWAP",
	}, tradesChan)
	if err != nil {
		println(err)
	}
	for {
		select {
		case <-logChan:
			log.Print("[Authorized]")
		case sub := <-subChan:
			channel, _ := sub.Arg.Get("channel")
			log.Printf("[Subscribed]\t%s", channel)
		case err := <-client.Ws.ErrChan:
			log.Printf("[Error]\t%+v", err)
			for _, datum := range err.Data {
				log.Printf("[Error]\t\t%+v", datum)
			}
		case i := <-tradesChan:
			log.Print(i.Trades[0])
		case b := <-client.Ws.DoneChan:
			log.Printf("[End]:\t%v", b)
			return
		}
	}
}
