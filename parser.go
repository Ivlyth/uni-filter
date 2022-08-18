package uni_filter

import (
	"fmt"
	"strings"
)

/*
- 期望用户完成的是什么呢？


- 期望提供给用户的又是什么呢？
  - 用户简单到只需要提供一个输入框从其终端用户处获取表达式
  - 将获取到的表达式提供给我们进行解析
  - 将解析好的表达式对象返回给用户，如果有错误，返回告知为何
  - 用户拿解析好的表达式对象验证任意数据


Note:
- 用户提供的数据必须是 map / array / slice / struct, 不能是简单类型(比如 int, int8 之类的)
- number 一律转换为 float64 进行比较 (但是理论上就没有丢失数据的可能吗？从 uint64 转换过来的时候)

(Pid=4837) or (Connections[].ConnectionInfo.RemoteIP=10.0.81.95 or Connections[].ConnectionInfo.RemoteIP=10.0.81.96) or (Connections[].ConnectionInfo.RemotePort=9091 or Connections[].ConnectionInfo.RemotePort=9200)
*/

func ParseFilterString(s string) (*Filter, error) {

	parts := strings.SplitN(s, FieldParamSep, 2)
	key := parts[0]
	value := ""
	if len(parts) == 2 {
		value = parts[1]
	}

	f := Filter{
		Key:   key,
		Value: value,
	}

	if f.Key == "" {
		return nil, fmt.Errorf("key can't be empty")
	}

	if key[0] == '!' {
		key = key[1:]
		f.opposite = true
	}

	parts = strings.SplitN(key, FieldOpSep, 2)
	opName := "eq"

	if len(parts) == 2 {
		key = parts[0]
		if key == "" {
			return nil, fmt.Errorf("key can't be empty")
		}

		opName = parts[1]
		if opName == "" {
			return nil, fmt.Errorf("op can't be empty")
		}
	}

	if newOpFunc, ok := ops[opName]; !ok {
		return nil, fmt.Errorf("unknown op name: %s", opName)
	} else {
		op, err := newOpFunc(f.Value)
		if err != nil {
			return nil, err
		}
		f.op = op
	}

	var keys []*Key
	for _, k := range strings.Split(key, FieldSep) {
		key, err := parseKey(k)
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	f.keys = keys

	return &f, nil
}
