# colar

a http router for go language. 

static router, dynamic router and regular expressions are supported.

# usage

```
package main

import (
    "fmt"
    "github.com/abnersun2016/colar"
    "net/http"
    "log"
)

func Index(context *context.Context) {
    context.ResponseWriter.Write(([]byte)("hello word"))
}

func getComments(context *context.Context) {
    context.ResponseWriter.Write(([]byte)("get comments"))
}

func updateComments(context *context.Context) {
    context.ResponseWriter.Write(([]byte)("update comments"))
}

func main() {
    router := colar.New()
    router.Get("/", Index)
    router.Get("/comments/:Id", getComments)
    router.Put("/comments/:Id([0-9]*)", updateComments)
    log.Fatal(http.ListenAndServe(":8080", router))
}
```
