package main

import (
	"flag"
	"fmt"
	"github.com/golang/groupcache"
	"github.com/renwuxun/datafront/front"
	"github.com/renwuxun/datafront/httphandler"
	"github.com/renwuxun/fasthttproute"
	"log"
	"os"
	"strings"
)

var (
	me     = flag.String("me", "0.0.0.0:8080", "Current peer")
	others = flag.String("others", "0.0.0.0:8081", `Other peers [separate by ","]`)
	help   = flag.Bool("h", false, "Show usage")
)

func usage() {
	_, _ = fmt.Fprintf(os.Stderr,
		`Usage: %s [-me] [-others]
Options:
`, os.Args[0])
	flag.PrintDefaults()
}

func init() {
	flag.Usage = usage
	flag.Parse()
}

func main() {
	if *help {
		flag.Usage()
		os.Exit(0)
	}

	otherPeers := strings.Split(*others, ",")
	front.RegisterPeers(*me, otherPeers)
	front.RegisterGetter("dummygroup", 64<<20, groupcache.GetterFunc(
		func(ctx groupcache.Context, key string, dest groupcache.Sink) error {
			_ = dest.SetString("key: [" + key + "] from peer: " + *me)
			log.Print("getter: " + key + " from peer: " + *me)
			return nil
		}))

	fasthttproute.Handle("/front", httphandler.FrontGet)
	fasthttproute.Handle("/front/purge", httphandler.FrontPurgeGroup)

	ServeFasthttpAndGroupcache := front.MakeCanServeGroupcache(fasthttproute.ServeFasthttp)
	ServeFasthttpAndGroupcache(*me, fasthttproute.DefaultHandler)
}
