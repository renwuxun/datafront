package front

import (
	"errors"
	"github.com/golang/groupcache"
	"github.com/golang/groupcache/groupcachepb"
	"github.com/golang/protobuf/proto"
	"github.com/renwuxun/datafront/helper"
	"github.com/valyala/fasthttp"
	"strings"
)

var peers *groupcache.HTTPPool

func RegisterPeers(me string, others []string) {
	if peers != nil {
		panic("u can only call RegisterPeers() once")
	}

	peers = groupcache.NewHTTPPoolOpts("http://"+me, &groupcache.HTTPPoolOptions{BasePath: GroupcacheBasePath})

	all := append(others, me)
	for i, v := range all {
		all[i] = "http://" + v
	}
	peers.Set(all...)
}

func RegisterGetter(groupName string, cacheSize int64, groupGetter groupcache.GetterFunc) {
	groupcache.NewGroup(groupName, cacheSize, groupGetter)
}

func serveGroupcache(ctx *fasthttp.RequestCtx) {
	path := helper.Bytes2str(ctx.Path())
	if !strings.HasPrefix(path, GroupcacheBasePath) {
		panic("HTTPPool serving unexpected path: " + path)
	}

	parts := strings.SplitN(path[len(GroupcacheBasePath):], "/", 2)

	if len(parts) != 2 {
		ctx.Response.SetStatusCode(400)
		ctx.Response.SetBodyString("bad request")
		return
	}

	groupName := parts[0]
	key := parts[1]

	group := groupcache.GetGroup(groupName)
	if group == nil {
		ctx.Response.SetStatusCode(404)
		ctx.Response.SetBodyString("no such group: " + groupName)
		return
	}

	group.Stats.ServerRequests.Add(1)

	var value []byte
	err := group.Get(ctx, key, groupcache.AllocatingByteSliceSink(&value))
	if err != nil {
		ctx.Response.SetStatusCode(500)
		return
	}

	body, err := proto.Marshal(&groupcachepb.GetResponse{Value: value})
	if err != nil {
		ctx.Response.SetStatusCode(500)
		return
	}

	ctx.Response.Header.Set("Content-Type", "application/x-protobuf")
	ctx.Response.SetBody(body)
}

func makeHandlerCanServeGroupcache(handler fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		sPath := helper.Bytes2str(ctx.Path())
		if strings.HasPrefix(sPath, GroupcacheBasePath) { // 如果是groupcache peer之间的请求
			serveGroupcache(ctx)
		} else {
			handler(ctx)
		}
	}
}

func Get(groupName, key string) ([]byte, error) {
	var err error

	group := groupcache.GetGroup(groupName)
	if group == nil {
		return nil, errors.New("no such group: " + groupName)
	}

	// PickPeer returns the peer that owns the specific key
	// and true to indicate that a remote peer was nominated.
	// It returns nil, false if the key owner is the current peer.
	peer, ok := peers.PickPeer(key)
	if !ok { // 说明这个key归自己管
		var value []byte
		err = group.Get(nil, key, groupcache.AllocatingByteSliceSink(&value))
		if err != nil {
			return nil, errors.New("error on pick peer: " + err.Error())
		}
		return value, nil
	}

	req := &groupcachepb.GetRequest{Group: &groupName, Key: &key}
	res := &groupcachepb.GetResponse{}
	err = peer.Get(nil, req, res)
	if err != nil {
		return nil, err
	}

	return res.Value, nil
}

func MakeCanServeGroupcache(fun func(addr string, h fasthttp.RequestHandler)) func(addr string, h fasthttp.RequestHandler) {
	return func(addr string, h fasthttp.RequestHandler) {
		h = makeHandlerCanServeGroupcache(h)
		fun(addr, h)
	}
}

// UpdatePeers 使得外层应用在运行中可以实时增减节点
func UpdatePeers(ps []string) {
	if peers == nil {
		panic("peers is nil")
	}
	if len(ps) < 2 {
		return
	}
	peers.Set(ps...)
}
