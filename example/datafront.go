package main

import (
	"flag"
	"github.com/golang/groupcache"
	"github.com/renwuxun/datafront/front"
	"github.com/renwuxun/datafront/httphandler"
	"github.com/renwuxun/fasthttproute"
	"log"
	"strings"
)

var (
	me     = flag.String("me", "0.0.0.0:8080", "Current peer's TCP address to listen to")
	others = flag.String("others", "0.0.0.0:8081", `Other peers's TCP address [separate by ","]`)
)

func main() {
	flag.Parse()

	otherPeers := strings.Split(*others, ",")
	front.RegisterPeers(*me, otherPeers)
	front.RegisterGetter("dummygroup", 64<<20, groupcache.GetterFunc(
		func(ctx groupcache.Context, key string, dest groupcache.Sink) error {
			dest.SetString("key: [" + key + "] from peer: " + *me)
			log.Print("getter: " + key + " from peer: " + *me)
			return nil
		}))

	fasthttproute.Handle("/front", httphandler.FrontGet)
	fasthttproute.Handle("/front/purge", httphandler.FrontPurgeGroup)

	ServeFasthttpAndGroupcache := front.MakeCanServeGroupcache(fasthttproute.ServeFasthttp)
	ServeFasthttpAndGroupcache(*me, fasthttproute.DefaultHandler)
}
