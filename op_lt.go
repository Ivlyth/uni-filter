package uni_filter

import (
	"errors"
	"fmt"
	"strconv"
)

type OPLT struct {
	s          string
	allowEqual bool
	f          float64
}

func (op *OPLT) Name() string {
	if op.allowEqual {
		return "lte"
	}
	return "lt"
}

func (op *OPLT) parse() error {
	if op.s == "" {
		return errors.New(fmt.Sprintf("op param for %s can not be empty", op.Name()))
	}
	i, err := strconv.ParseFloat(op.s, 64)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("value for op %s should be valid float64", op.Name()))
	}
	op.f = i
	return nil
}

func (op *OPLT) check(v any, exists bool) bool {
	if !exists {
		return false
	}

	var vf float64

	switch nv := v.(type) {
	case uint:
		vf = float64(nv)
	case uint8:
		vf = float64(nv)
	case uint16:
		vf = float64(nv)
	case uint32:
		vf = float64(nv)
	case uint64:
		vf = float64(nv)
	case int:
		vf = float64(nv)
	case int8:
		vf = float64(nv)
	case int16:
		vf = float64(nv)
	case int32:
		vf = float64(nv)
	case int64:
		vf = float64(nv)
	case float32:
		vf = float64(nv)
	case float64:
		vf = nv
	default:
		return false
	}

	if op.allowEqual {
		return vf <= op.f
	} else {
		return vf < op.f
	}
}

func NewOPLT(s string) (OP, error) {
	op := &OPLT{s: s}
	return op, op.parse()
}

func NewOPLTE(s string) (OP, error) {
	op := &OPLT{s: s, allowEqual: true}
	return op, op.parse()
}

func init() {
	register("lt", NewOPLT)
	register("lte", NewOPLTE)
}
