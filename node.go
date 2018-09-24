/*
store the path info
*/
package colar

import (
	"bytes"
	pathParam "colar/context/param"
	"fmt"
	"regexp"
	"strings"
)

type nodeType uint8

const (
	root nodeType = iota
	normal
	param
)

type node struct {
	regex    *regexp.Regexp
	nType    nodeType
	nValue   string
	cPrefix  string
	children []*node
	handle   Handler
}

func (root *node) insertNode(path string, handle Handler, caseSensitive bool) {

	fullPath := path
	if !caseSensitive {
		fullPath = strings.ToLower(path)
	}
	fullPath = revampTrailSlash(fullPath)
	strs := strings.Split(fullPath, "/")
	if _, error := checkParams(strs); error != nil {
		panic(error)
		return
	}
	if len(strs) > 0 {
		curr := root
		i := 0
		for ; i < len(strs); i++ {
			if len(curr.children) == 0 {
				if strings.HasPrefix(strs[i], ":") {
					curr = curr.insert(strs[i][1:], true)
				} else {
					curr = curr.insert(strs[i], false)
				}
			} else {
				if strings.HasPrefix(strs[i], ":") {
					curr = curr.insert(strs[i][1:], true)
				} else {
					//遍历子节点，查找最长公共前缀
					buf := new(bytes.Buffer)
					for m := range curr.children {
						if curr.children[m].nType != param {
							var min string
							if len(strs[i]) < len(curr.children[m].nValue) {
								min = strs[i]
							} else {
								min = curr.children[m].nValue
							}
							for k := 0; k < len(min); k++ {
								if strings.Compare(string(curr.children[m].nValue[k]), string(strs[i][k])) == 0 {
									buf.WriteString(string(strs[i][k]))
								} else {
									break
								}
							}
							//长到公共前缀
							if buf.Len() > 0 {
								if buf.Len() < len(curr.children[m].nValue) {
									//分裂该节点
									child := curr.children[m]
									child.regex = regexp.MustCompile(child.nValue[buf.Len():])
									child.nValue = child.regex.String()
									//声明该子节点的父节点
									parent := new(node)
									parent.nType = normal
									parent.regex = regexp.MustCompile(buf.String())
									parent.nValue = buf.String()
									parent.cPrefix += child.nValue[:1]
									parent.children = append(parent.children, child)

									//删除该节点的子节点m
									curr.children = remove(curr.children, m)
									//添加新节点
									curr.children = append(curr.children, parent)

									if buf.Len() < len(strs[i]) {
										curr = parent.insert(strs[i][buf.Len():], false)
									} else {
										curr = parent
									}
								} else {
									if buf.Len() < len(strs[i]) {
										curr = curr.children[m].insert(strs[i][buf.Len():], false)
									} else {
										curr = curr.children[m]
									}
								}
								break
							}
						}
					}
					//未找到公共前缀
					if buf.Len() == 0 {
						curr = curr.insert(strs[i], false)
					}

				}

			}
			if i == len(strs)-1 {
				curr.handle = handle
			}
		}
	} else {
		if len(fullPath) == 0 || strings.Compare("/", fullPath) == 0 {
			root.handle = handle
		}
	}
}

func (curr *node) insert(nValue string, isParam bool) *node {
	if isParam {
		vIndex := -1
		for i := range curr.children {
			if curr.children[i].nType == param {
				vIndex = i
			}
		}
		var regex string
		if regexp.MustCompile(".+\\(.+\\)").MatchString(nValue) {
			regex = nValue[strings.Index(nValue, "(")+1 : strings.LastIndex(nValue, ")")]
		} else {
			regex = ".*"
		}
		if vIndex == -1 {
			n := new(node)
			n.regex = regexp.MustCompile(regex)
			n.nType = param
			if strings.Compare(regex, ".*") == 0 {
				n.nValue = nValue
			} else {
				n.nValue = nValue[:strings.Index(nValue, "(")]
			}
			curr.cPrefix += "*"
			curr.children = append(curr.children, n)
			return n
		} else {
			n := curr.children[vIndex]
			n.regex = regexp.MustCompile(regex)
			if strings.Compare(regex, ".*") == 0 {
				n.nValue = nValue
			} else {
				n.nValue = nValue[:strings.Index(nValue, "(")]
			}
			return n
		}
	} else {
		vIndex := -1
		for i := range curr.children {
			if strings.Compare(curr.children[i].nValue, nValue) == 0 && curr.children[i].nType != param {
				vIndex = i
				break
			}
		}
		if vIndex == -1 {
			n := new(node)
			n.regex = regexp.MustCompile(nValue)
			n.nValue = nValue
			n.nType = normal
			curr.cPrefix += nValue[:1]
			curr.children = append(curr.children, n)
			return n
		} else {
			n := curr.children[vIndex]
			return n
		}
	}
}

func (root *node) findNode(path string, caseSensitive bool) (*node, *pathParam.PathParams) {
	fullPath := path
	if !caseSensitive {
		fullPath = strings.ToLower(path)
	}
	fullPath = revampTrailSlash(fullPath)
	strs := strings.Split(fullPath, "/")

	curr := root
	params := make(map[string][]string)
	index := 0
	nValue := strs[index]
loop:
	for {
		//查找前缀子节点
		var matchIndex = -1
		var catchAllIndex = -1
		for i, prefix := range curr.cPrefix {
			if strings.HasPrefix(nValue, string(prefix)) {
				matchIndex = i
				for k, ch := range curr.children[i].nValue {
					if strings.Compare(string(ch), nValue[k:k+1]) != 0 {
						matchIndex = -1
					}
				}
			}
			if strings.Compare("*", string(prefix)) == 0 && curr.children[i].regex.MatchString(nValue) {
				catchAllIndex = i
			}
		}
		//确定优先级
		if matchIndex >= 0 && catchAllIndex >= 0 {
			if strings.Compare(curr.children[matchIndex].nValue, nValue) == 0 {
				catchAllIndex = -1
			} else {
				matchIndex = -1
			}
		}
		if matchIndex >= 0 {
			suffix := nValue[len(curr.children[matchIndex].nValue):]
			curr = curr.children[matchIndex]
			if len(suffix) > 0 {
				nValue = suffix
			} else {
				index += 1
				if index < len(strs) {
					nValue = strs[index]
				} else {
					body := &pathParam.PathParams{params}
					return curr, body
				}
			}
			continue loop
		} else if catchAllIndex >= 0 {
			curr = curr.children[catchAllIndex]
			if params[curr.nValue] == nil {
				params[curr.nValue] = make([]string, 0)
			}
			params[curr.nValue] = append(params[curr.nValue], strs[index])
			if curr.nValue == "filepath" {
				filepath := new(bytes.Buffer)
				for i := index + 1; i < len(strs); i++ {
					filepath.WriteString("/")
					filepath.WriteString(strs[i])
				}
				params[curr.nValue][0] += filepath.String()
				body := &pathParam.PathParams{params}
				return curr, body
			} else {
				index += 1
				if index < len(strs) {
					nValue = strs[index]
				} else {
					body := &pathParam.PathParams{params}
					return curr, body
				}
			}
		} else {
			break
		}
	}
	return nil, nil
}

func remove(nodes []*node, i int) []*node {
	newSlice := make([]*node, 0)
	for k := range nodes {
		if k != i {
			newSlice = append(newSlice, nodes[k])
		}
	}
	return newSlice
}

func checkParams(paths []string) (uint8, error) {
	var count uint8
	var error error
	for _, path := range paths {
		if strings.HasPrefix(path, ":") {
			if len(path) > 1 {
				isRegex := regexp.MustCompile(".*\\(.*\\)").MatchString(path)
				if isRegex {
					path = path[:strings.Index(path, "(")]
				}
				if len(path) > 1 {
					count += 1
					continue
				}
			}
			error = fmt.Errorf("the param length of following behind ':' must be greater than zero,paths=%s", paths)
			return count, error
		}
	}
	return count, error
}

func revampTrailSlash(path string) string {
	buf := new(bytes.Buffer)
	splits := strings.Split(path, "/")
	for _, path := range splits {
		if strings.Compare(path, "") != 0 {
			buf.WriteString(path)
			buf.WriteString("/")
		}
	}
	if strings.HasSuffix(buf.String(), "/") {
		return buf.String()[:len(buf.String())-1]
	}
	return buf.String()
}
