package uni_filter

type OPEmpty struct {
	s string
}

func (op *OPEmpty) Name() string {
	return "empty"
}

func (op *OPEmpty) parse() error {
	return nil
}

func (op *OPEmpty) check(v any, exists bool) bool {
	if !exists {
		return false
	}
	s, ok := v.(string)
	if ok {
		return s == ""
	}
	return false
}

func NewOPEmpty(s string) (OP, error) {
	op := &OPEmpty{s: s}
	return op, op.parse()
}

func init() {
	register("em", NewOPEmpty)
	register("empty", NewOPEmpty)
}
