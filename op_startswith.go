package uni_filter

import (
	"errors"
	"fmt"
	"strings"
)

type OPStarts struct {
	s          string
	is         string
	ignoreCase bool
}

func (op *OPStarts) Name() string {
	if op.ignoreCase {
		return "istarts"
	}
	return "starts"
}

func (op *OPStarts) parse() error {
	if op.s == "" {
		return errors.New(fmt.Sprintf("op param for %s can not be empty", op.Name()))
	}
	if op.ignoreCase {
		op.is = strings.ToLower(op.s)
	}
	return nil
}

func (op *OPStarts) check(v any, exists bool) bool {
	if !exists {
		return false
	}
	s := convert2string(v)
	if op.ignoreCase {
		s = strings.ToLower(s)
		return strings.HasPrefix(s, op.is)
	} else {
		return strings.HasPrefix(s, op.s)
	}
}

func NewOPStarts(s string) (OP, error) {
	op := &OPStarts{s: s}
	return op, op.parse()
}

func NewOPIStarts(s string) (OP, error) {
	op := &OPStarts{s: s, ignoreCase: true}
	return op, op.parse()
}

func init() {
	register("starts", NewOPStarts)
	register("istarts", NewOPIStarts)
}
