package colar

import (
	"bytes"
	"colar/context"
	"fmt"
	"regexp"
	"strings"
	"testing"
)

var n *node

const (
	path1 = "app/ab/admin/simple/:Id"
	path2 = "app/abd/admin/simple/hjds"
	path3 = "app/ad/admin/:apple/"
	path4 = "app/ad/admin/:hdjdjd/Id"
	path5 = "app/ad/admin/:simple/:Idd"
	path6 = "app/ad/ads/simple/:sId"
	path7 = "app/chj/admin/simple/:dsds([0-9]*)"
	path8 = "app/kj/admin/simple/:Idsd"
)

func init() {
	n = new(node)
	n.nType = root
	n.handle = nil

	n.insertNode(path1, handle1, true)
	n.insertNode(path2, handle2, true)
	n.insertNode(path3, handle3, true)
	n.insertNode(path4, handle4, true)
	n.insertNode(path5, handle5, true)
	n.insertNode(path6, handle6, true)
	n.insertNode(path7, handle7, true)
	n.insertNode(path8, handle8, true)
}

func handle1(context *context.Context) {
	context.ResponseWriter.Write(([]byte)("hello word"))
	fmt.Println(1)
}
func handle2(context *context.Context) {
	context.ResponseWriter.Write(([]byte)("hello word"))
	fmt.Println(context.PathParams.GetByName("Id"))
}
func handle3(context *context.Context) {
	context.ResponseWriter.Write(([]byte)("hello word"))
	fmt.Println(3)
}
func handle4(context *context.Context) {
	fmt.Println(4)
}
func handle5(context *context.Context) {
	fmt.Println(5)
}
func handle6(context *context.Context) {
	fmt.Println(6)
}
func handle7(context *context.Context) {
	fmt.Println(7)
}
func handle8(context *context.Context) {
	fmt.Println(8)
}

func TestFindNodeStatic(t *testing.T) {
	node, body := n.findNode("app/ad/admin/apple/Idd", true)
	fmt.Println(node.nValue)
	fmt.Println(node.regex.String())
	fmt.Println(node.cPrefix)
	fmt.Println(body.Params)

	node, body = n.findNode("app/chj/admin/simple/32635f26", true)
	fmt.Println(node.nValue)
	fmt.Println(node.regex.String())
	fmt.Println(node.cPrefix)
	fmt.Println(body.Params)

	ctx := &context.Context{}
	node, body = n.findNode("app/ad/admin/apple/323fdf23", true)
	fmt.Println(node.nValue)
	fmt.Println(node.regex.String())
	fmt.Println(node.cPrefix)
	fmt.Println(body.Params)
	fmt.Println(body.GetByName("Idd"))
	node.handle(ctx)
	fmt.Println(strings.Compare("a", "a"))
	fmt.Println(regexp.MustCompile("A").MatchString("a"))
}
