package colar

import (
	"github.com/abnersun2016/colar"
	"log"
	"net/http"
	"testing"
)

func TestRouter_AddMethod(t *testing.T) {
	route := colar.New()
	route.AddMethod("*", "gitchat/posts", handle1)
	route.AddMethod("post", "gitchat/posts/:Id", handle2)
	route.AddMethod("*", "gitchat/posts/Gomments/ID", handle3)
	route.ServeFiles("/static/:filepath", http.Dir("/Users/Jean/Downloads/upload"))
	log.Fatal(http.ListenAndServe(":8080", route))
}
