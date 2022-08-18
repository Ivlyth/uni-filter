package uni_filter

import "time"

type OPTimeRange struct {
	s        string
	absolute bool
	ms       bool

	f1 float64
	f2 float64
}

func (op *OPTimeRange) Name() string {
	return "timerange"
}

func (op *OPTimeRange) parse() error {
	var f1, f2 float64

	now := time.Now()

	if op.absolute {
		f1 = float64(now.Unix())
		f2 = float64(now.Unix() + 600)
	} else {
		now = now.Add(time.Second * 600) // 10min
		begin := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute()-now.Minute()%10, 0, 0, now.Location())

		f1 = float64(begin.Unix())
		f2 = float64(begin.Unix() + 600)
	}

	if op.ms {
		f1 = f1 * 1000.0
		f2 = f2 * 1000.0
	}

	op.f1 = f1
	op.f2 = f2

	return nil
}

func (op *OPTimeRange) check(v any, exists bool) bool {
	if !exists {
		return false
	}
	vf, ok := v.(float64)
	if ok {
		return vf >= op.f1 && vf < op.f2
	}
	return false
}

func NewOPTimeRangeA(s string) (OP, error) {
	op := &OPTimeRange{
		s:        s,
		absolute: true,
	}
	return op, op.parse()
}

func NewOPTimeRangeAMS(s string) (OP, error) {
	op := &OPTimeRange{
		s:        s,
		absolute: true,
		ms:       true,
	}
	return op, op.parse()
}

func NewOPTimeRangeR(s string) (OP, error) {
	op := &OPTimeRange{
		s: s,
	}
	return op, op.parse()
}

func NewOPTimeRangeRMS(s string) (OP, error) {
	op := &OPTimeRange{
		s:  s,
		ms: true,
	}
	return op, op.parse()
}

func init() {
	register("a10min", NewOPTimeRangeA)
	register("a10minms", NewOPTimeRangeAMS)
	register("r10min", NewOPTimeRangeR)
	register("r10minms", NewOPTimeRangeRMS)
}
