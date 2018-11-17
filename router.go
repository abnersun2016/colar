/*
Matching request connection
*/
package colar

import (
	"bytes"
	"colar/context"
	"net/http"
	"strings"
	"sync"
)

type Handler func(context *context.Context)

type Router struct {
	trees map[string]*node

	HandlerNotFound Handler

	MethodNotAllowed Handler

	RecoverHandler func(*context.Context, interface{})

	CaseSensitive bool

	pool sync.Pool
}

var (
	getNode     *node
	postNode    *node
	putNode     *node
	patchNode   *node
	deleteNode  *node
	connectNode *node
	options     *node
	trace       *node
)

//surpported http method
var (
	HTTPMETHOD = map[string]bool{
		"GET":     true,
		"POST":    true,
		"PUT":     true,
		"DELETE":  true,
		"PATCH":   true,
		"OPTIONS": true,
		"HEAD":    true,
		"TRACE":   true,
		"CONNECT": true,
	}
)

func New() *Router {
	rooter := &Router{trees: make(map[string]*node), CaseSensitive: true}
	rooter.pool.New = func() interface{} {
		return &context.Context{}
	}
	return rooter
}

func (r *Router) Get(path string, handler Handler) {
	r.AddMethod(http.MethodGet, path, handler)
}

func (r *Router) Post(path string, handler Handler) {
	r.AddMethod(http.MethodPost, path, handler)
}

func (r *Router) Put(path string, handler Handler) {
	r.AddMethod(http.MethodPut, path, handler)
}

func (r *Router) Delete(path string, handler Handler) {
	r.AddMethod(http.MethodDelete, path, handler)
}

func (r *Router) Head(path string, handler Handler) {
	r.AddMethod(http.MethodDelete, path, handler)
}

func (r *Router) Options(path string, handler Handler) {
	r.AddMethod(http.MethodOptions, path, handler)
}

func (r *Router) Trace(path string, handler Handler) {
	r.AddMethod(http.MethodTrace, path, handler)
}

func (r *Router) Connect(path string, handler Handler) {
	r.AddMethod(http.MethodConnect, path, handler)
}

func (r *Router) Patch(path string, handler Handler) {
	r.AddMethod(http.MethodPatch, path, handler)
}

func (r *Router) Any(path string, handler Handler) {
	r.AddMethod("*", path, handler)
}

func (r *Router) AddMethod(method, path string, handler Handler) {
	method = strings.ToUpper(method)
	if !HTTPMETHOD[method] && !strings.EqualFold("*", method) {
		panic("http method: " + method + "isn't be supported")
	}
	if strings.EqualFold("*", method) {
		for k := range HTTPMETHOD {
			r.initRootTree(k)
			r.trees[k].insertNode(path, handler, r.CaseSensitive)
		}
	} else {
		r.initRootTree(method)
		r.trees[method].insertNode(path, handler, r.CaseSensitive)
	}
}

func (r *Router) initRootTree(method string) {
	if r.trees[method] == nil {
		rootNode := new(node)
		rootNode.nType = root
		r.trees[method] = rootNode
	}
}

//recover from panic
func (r *Router) recov(context *context.Context) {
	if recov := recover(); r != nil {
		r.RecoverHandler(context, recov)
	}
}

func (r *Router) checkAllowedMethod(path string) string {
	method := new(bytes.Buffer)
	for k := range HTTPMETHOD {
		if node, _ := r.trees[k].findNode(path, r.CaseSensitive); node != nil {
			method.WriteString(k)
			method.WriteString("_")
		}
	}
	return method.String()
}

func (r *Router) ServeFiles(path string, root http.FileSystem) {
	fullPath := path
	if !r.CaseSensitive {
		fullPath = strings.ToLower(path)
	}
	fullPath = revampTrailSlash(fullPath)
	strs := strings.Split(fullPath, "/")
	if strs[len(strs)-1] != ":filepath" {
		panic("path must end with /:filepath or /:filepath/ in path '" + path + "'")
	}
	fileServer := http.FileServer(root)
	r.Get(path, func(context *context.Context) {
		context.Request.URL.Path = context.PathParams.GetByName("filepath").(string)
		fileServer.ServeHTTP(context.ResponseWriter, context.Request)
	})
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	context := r.pool.Get().(*context.Context)
	defer r.pool.Put(context)
	context.Refresh(w, req)
	if r.RecoverHandler != nil {
		defer r.recov(context)
	}
	if root := r.trees[req.Method]; root != nil {
		if node, pathParams := root.findNode(path, r.CaseSensitive); node != nil {
			context.PathParams = pathParams
			node.handle(context)
			return
		}
		if r.HandlerNotFound != nil {
			r.HandlerNotFound(context)
		} else {
			http.Error(w,
				http.StatusText(http.StatusNotFound),
				http.StatusNotFound,
			)
		}
		return
	} else {
		allowedMethod := r.checkAllowedMethod(path)
		if len(allowedMethod) > 0 {
			if r.MethodNotAllowed != nil {
				r.MethodNotAllowed(context)
			} else {
				http.Error(w,
					http.StatusText(http.StatusMethodNotAllowed),
					http.StatusMethodNotAllowed,
				)
			}
		} else {
			if r.HandlerNotFound != nil {
				r.HandlerNotFound(context)
			} else {
				http.Error(w,
					http.StatusText(http.StatusNotFound),
					http.StatusNotFound,
				)
			}
		}
		return
	}

}
