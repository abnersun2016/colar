/*
store the http Request,ResponseWriter and PathParams
*/
package context

import (
	"colar/context/param"
	"net/http"
)

type Context struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	PathParams     *param.PathParams
}

//when a new connect build,refresh the connect context
func (ctx *Context) Refresh(rw http.ResponseWriter, r *http.Request) {
	ctx.Request = r
	ctx.ResponseWriter = rw
	ctx.PathParams = nil
}
