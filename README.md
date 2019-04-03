# colar

a http router for go language. 

static router, dynamic router and regular expressions are supported.

# usage

```
package main

import (
    "fmt"
    "github.com/aidensuen/colar"
    "net/http"
    "log"
)

func index(context *context.Context) {
    context.ResponseWriter.Write(([]byte)("hello word"))
}

func getComments(context *context.Context) {
    context.ResponseWriter.Write(([]byte)("get comments"))
    fmt.Println(context.PathParams.GetByName("Id"))
}

func updateComments(context *context.Context) {
    context.ResponseWriter.Write(([]byte)("update comments"))
    fmt.Println(context.PathParams.GetByName("Id"))
}

func main() {
    router := colar.New()
    router.Get("/", index)
    router.Get("/comments/:Id", getComments)
    router.Put("/comments/:Id([0-9]*)", updateComments)
    log.Fatal(http.ListenAndServe(":8080", router))
}
```
