package uni_filter

type OPExists struct {
	s string
}

func (op *OPExists) Name() string {
	return "exists"
}

func (op *OPExists) parse() error {
	return nil
}

func (op *OPExists) check(v any, exists bool) bool {
	return exists
}

func NewOPExists(s string) (OP, error) {
	op := &OPExists{s: s}
	return op, op.parse()
}

func init() {
	register("ex", NewOPExists)
	register("exists", NewOPExists)
}
