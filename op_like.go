package uni_filter

import (
	"errors"
	"fmt"
	"strings"
)

type OPLike struct {
	s          string
	is         string
	ignoreCase bool
}

func (op *OPLike) Name() string {
	if op.ignoreCase {
		return "ilike"
	}
	return "like"
}

func (op *OPLike) parse() error {
	if op.s == "" {
		return errors.New(fmt.Sprintf("op param for %s can not be empty", op.Name()))
	}
	if op.ignoreCase {
		op.is = strings.ToLower(op.s)
	}
	return nil
}

func (op *OPLike) check(v any, exists bool) bool {
	if !exists {
		return false
	}
	s := convert2string(v)
	if op.ignoreCase {
		s = strings.ToLower(s)
		return strings.Contains(s, op.is)
	} else {
		return strings.Contains(s, op.s)
	}
}

func NewOPLike(s string) (OP, error) {
	op := &OPLike{s: s}
	return op, op.parse()
}

func NewOPILike(s string) (OP, error) {
	op := &OPLike{s: s, ignoreCase: true}
	return op, op.parse()
}

func init() {
	register("like", NewOPLike)
	register("ilike", NewOPILike)
}
