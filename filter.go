package uni_filter

import (
	"errors"
	"fmt"
)

const FieldParamSep = "="
const FieldOpSep = "__"
const FieldSep = "."
const OPParamsSep = "__"

type Key struct {
	name string
}

func parseKey(s string) (*Key, error) {

	if len(s) == 0 {
		return nil, errors.New("empty key")
	}

	return &Key{
		name: s,
	}, nil
}

type Value struct {
	S string // clean string
	V any
}

type Filter struct {
	// 原始未解析的 key
	Key string
	// 原始未解析的 value
	Value string

	// 解析后的 keys
	keys []*Key

	//// 从原始 key 中提取出来的 op
	op OP

	// 从原始 key 中提取到的 ! flag
	opposite bool
}

func (f *Filter) String() string {
	s := f.Key
	if f.Value != "" {
		s = fmt.Sprintf("%s=%s", s, f.Value)
	}
	// FIXME show about op
	return s
}
