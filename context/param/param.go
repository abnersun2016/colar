/*
store the params in path
*/
package param

type PathParams struct {
	Params map[string][]string
}

//get param value by name
func (pathParams *PathParams) GetByName(key string) interface{} {
	if len(pathParams.Params[key]) > 1 {
		return pathParams.Params[key]
	} else if len(pathParams.Params[key]) == 1 {
		return pathParams.Params[key][0]
	}
	return pathParams.Params[key]
}
