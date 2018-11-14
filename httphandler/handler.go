package httphandler

import (
	"fmt"
	"github.com/renwuxun/datafront/front"
	"github.com/renwuxun/datafront/helper"
	"github.com/valyala/fasthttp"
	"strconv"
)

const keySeperator = '-'

var groupKeyVersion = map[string]uint64{}

func keyGen(groupName string, incompleteKey []byte) []byte {
	incompleteKey = append(incompleteKey, keySeperator)
	if _, ok := groupKeyVersion[groupName]; !ok {
		groupKeyVersion[groupName] = 0
	}
	return strconv.AppendUint(incompleteKey, groupKeyVersion[groupName], 10)
}

func FrontPurgeGroup(ctx *fasthttp.RequestCtx) {
	groupName := helper.Bytes2str(ctx.QueryArgs().Peek("group"))
	if groupName == "" {
		ctx.Response.SetBodyString("invalid group")
		return
	}

	groupKeyVersion[groupName]++

	fmt.Fprintf(ctx, "group:%s, key version: %d", groupName, groupKeyVersion[groupName])
}

func FrontGet(ctx *fasthttp.RequestCtx) {
	groupName := helper.Bytes2str(ctx.QueryArgs().Peek("group"))
	key := ctx.QueryArgs().Peek("key")

	key = keyGen(groupName, key)

	data, err := front.Get(groupName, helper.Bytes2str(key))
	if err != nil {
		fmt.Print(err.Error())
	}

	ctx.Response.SetBody(data)
}
