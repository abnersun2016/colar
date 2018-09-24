package colar

import (
	"log"
	"net/http"
	"testing"
)

var route *Router

func init() {
	route = New()
}
func TestRouter_AddMethod(t *testing.T) {
	route.AddMethod("*", "gitchat/posts", handle1)
	route.AddMethod("post", "gitchat/posts/:Id", handle2)
	route.AddMethod("*", "gitchat/posts/Gomments/ID", handle3)
	route.ServeFiles("/static/:filepath", http.Dir("/Users/Jean/Downloads/upload"))
	log.Fatal(http.ListenAndServe(":8080", route))
}
