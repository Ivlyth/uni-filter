package uni_filter

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

var zeroValue = reflect.Value{}

type Logic uint8

func (l Logic) Strings() string {
	switch l {
	case AND:
		return "AND"
	case OR:
		return "OR"
	default:
		return "UNKNOWN LOGIC"
	}
}

var AND Logic = 0
var OR Logic = 1

type IExpr interface {
	Match(v any) bool
	Strings() string
}

type Expr struct {
	token Token
}

func (e *Expr) Strings() string {
	return e.token.value
}

type LogicExpr struct {
	Expr
	logic Logic
	left  IExpr // SimpleExpr or LogicExpr or StartGroupExpr
	right IExpr // SimpleExpr or LogicExpr or StartGroupExpr
}

func (le *LogicExpr) Strings() string {
	if le.left != nil && le.right != nil {
		return fmt.Sprintf("%s %s %s", le.left.Strings(), le.logic.Strings(), le.right.Strings())
	} else {
		return "logic.left or right is nil"
	}
}

func (e *LogicExpr) Match(data any) bool {
	lret := e.left.Match(data)
	if e.logic == OR && lret {
		return true
	}

	rret := e.right.Match(data)
	if e.logic == OR {
		return rret
	}

	return lret && rret
}

type StartGroupExpr struct {
	Expr
	len      int     // length of exprs
	exprs    []IExpr // SimpleExpr or LogicExpr
	stack    *Stack[IExpr]
	complete bool // already end
}

func (sge *StartGroupExpr) Strings() string {
	var buffs []string
	for _, expr := range sge.exprs {
		buffs = append(buffs, expr.Strings())
	}
	return fmt.Sprintf("(%s)", strings.Join(buffs, " "))
}

func (e *StartGroupExpr) Match(data any) bool {
	var flags = make([]bool, len(e.exprs))
	for i, expr := range e.exprs {
		flags[i] = expr.Match(data)
	}
	return All(flags)
}

type EndGroupExpr struct {
	Expr
}

func (e *EndGroupExpr) Match(data any) bool {
	return false
}

type SimpleExpr struct {
	Expr
	expr string // eg: Pid=32, Connections[].ConnectionInfo.FD=44

	filter *Filter
}

func (se *SimpleExpr) Strings() string {
	return fmt.Sprintf("'%s'", se.token.value)
}

func (se *SimpleExpr) matchStruct(instance reflect.Value, keys []*Key) bool {
	var v reflect.Value
	var keyExists bool

	instancePointer := reflect.New(instance.Type())
	instancePointer.Elem().Set(instance)

	for i, key := range keys {
		v = instance.FieldByName(key.name)

		if v == zeroValue { // if not found as field, then try as method
			v = instancePointer.MethodByName(key.name)
		}

		if v == zeroValue { // still zero
			keyExists = false
		} else {
			keyExists = true
		}

		if i != len(keys)-1 && !keyExists {
			return false
		}

		switch v.Kind() {
		case reflect.Slice, reflect.Array:
			fallthrough
		case reflect.Func:
			fallthrough
		case reflect.Struct, reflect.Map:
			return se.match(v.Interface(), keys[i+1:], true)
		case reflect.Pointer:
			if v.IsNil() {
				return false
			}
			v = reflect.Indirect(v)
			switch v.Kind() {
			case reflect.Slice, reflect.Array:
				fallthrough
			case reflect.Func:
				fallthrough
			case reflect.Struct, reflect.Map:
				return se.match(v.Interface(), keys[i+1:], true)
			case reflect.Chan, reflect.Interface, reflect.Invalid, reflect.UnsafePointer:
				// not support yet
				return false
			}
		case reflect.Chan, reflect.Interface, reflect.Invalid, reflect.UnsafePointer:
			// not support yet
			return false
		}

		// here should always be simple type, like int, bool, float, etc
		//panic("can not go here")
	}

	if v == zeroValue {
		return false
	}

	ret := se.filter.op.check(v.Interface(), keyExists)
	if se.filter.opposite {
		ret = !ret
	}
	return ret
}

func (se *SimpleExpr) matchSlice(datas *[]any, keys []*Key) bool {
	for _, data := range *datas {
		ret := se.match(data, keys, true)
		if ret {
			return true
		}
	}
	return false
}

func (se *SimpleExpr) matchFunc(val reflect.Value, keys []*Key) bool {
	t := val.Type()

	if t.NumIn() != 0 {
		return false
	}

	if t.NumOut() != 1 {
		return false
	}

	rets := val.Call(nil)
	if len(rets) != 1 {
		return false
	}

	return se.match(rets[0].Interface(), keys, true)
}

func (se *SimpleExpr) matchMap(data *map[string]any, keys []*Key) bool {
	var v any
	var keyExists bool

	for i, key := range keys {
		v, keyExists = (*data)[key.name]

		if i == len(keys)-1 {
			break
		} else if !keyExists {
			return false
		}

		val := reflect.ValueOf(v)
		switch val.Kind() {
		case reflect.Slice, reflect.Array:
			fallthrough
		case reflect.Func:
			fallthrough
		case reflect.Struct, reflect.Map:
			return se.match(v, keys[i+1:], true)
		case reflect.Pointer:
			if val.IsNil() {
				return false
			}
			val = reflect.Indirect(val)
			switch val.Kind() {
			case reflect.Slice, reflect.Array:
				fallthrough
			case reflect.Func:
				fallthrough
			case reflect.Struct, reflect.Map:
				return se.match(val.Interface(), keys[i+1:], true)
			}
		case reflect.Chan, reflect.Interface, reflect.Invalid, reflect.UnsafePointer:
			// not support yet
			return false
		}

		// here should always be simple type, like int, bool, float, etc
		panic("can not go here")
	}

	ret := se.filter.op.check(v, keyExists)
	if se.filter.opposite {
		ret = !ret
	}
	return ret
}

func (se *SimpleExpr) Match(data any) bool {
	return se.match(data, se.filter.keys, false)
}

func (se *SimpleExpr) match(data any, keys []*Key, allowArray bool) bool {
	val := reflect.ValueOf(data)

	switch val.Kind() {
	case reflect.Struct:
		return se.matchStruct(val, keys)
	case reflect.Map:
		n := val.Len()
		iter := val.MapRange()
		nm := make(map[string]any, n)
		for iter.Next() {
			k := iter.Key()
			if k.Kind() != reflect.String {
				return false
			}
			v := iter.Value()
			nm[k.String()] = v.Interface()
		}
		return se.matchMap(&nm, keys)
	case reflect.Slice, reflect.Array:
		if !allowArray {
			return false
		}
		n := val.Len()
		if n == 0 {
			return false
		}
		var ns = make([]any, n)
		for idx := 0; idx < n; idx++ {
			ns[idx] = val.Index(idx).Interface()
		}
		return se.matchSlice(&ns, keys)
	case reflect.Func:
		if !allowArray {
			return false
		}
		return se.matchFunc(val, keys)
	case reflect.Pointer:
		if val.IsNil() {
			return false
		}
		val = reflect.Indirect(val)
		return se.match(val.Interface(), keys, allowArray)
	default: // unsupported type
		return false
	}
}

type Token struct {
	start int
	end   int
	value string
	err   error
}

func (t Token) String() string {
	return fmt.Sprintf("token '%s' from %d to %d", t.value, t.start, t.end)
}

type TokenError struct {
	start int
	end   int
	msg   string
}

func (te TokenError) Error() string {
	return fmt.Sprintf("%s (from %d to %d)", te.msg, te.start, te.end)
}

type Tokenizer struct {
	expr string
	c    chan Token
}

func (t *Tokenizer) Parse() chan Token {
	t.c = make(chan Token, 20)
	go t._parse()
	return t.c
}

func (t *Tokenizer) _parse() {
	defer close(t.c)

	lastToken := ""

	currentToken := ""

	start := 0
	current := 0

	var isToken = func(t string) bool {
		lt := strings.ToLower(t)
		switch lt {
		case "and", "or":
			return true
		}
		return false
	}

	for i, c := range t.expr {
		current = i

		if unicode.IsSpace(c) {
			if currentToken != "" {
				if isToken(currentToken) {
					if lastToken != "" {
						t.c <- Token{
							start: start,
							end:   current + 1,
							value: lastToken,
						}
						lastToken = ""
					}
					t.c <- Token{
						start: start,
						end:   current + 1,
						value: currentToken,
					}
				} else {
					lastToken += currentToken
				}
				currentToken = ""
			}
		} else if c == '(' || c == ')' {
			if (lastToken + currentToken) != "" {
				t.c <- Token{
					start: start,
					end:   current + 1,
					value: lastToken + currentToken,
				}
				lastToken = ""
				currentToken = ""
			}
			t.c <- Token{
				start: current,
				end:   current + 1,
				value: string(c),
			}
		} else {
			if currentToken == "" {
				start = i
			}
			currentToken += string(c)
		}
	}

	if (lastToken + currentToken) != "" {
		t.c <- Token{
			start: start,
			end:   current + 1,
			value: lastToken + currentToken,
		}
	}
}

func createExprFromToken(t Token) IExpr {
	if t.value == "(" {
		return &StartGroupExpr{
			Expr: Expr{
				token: t,
			},
			stack: NewStack[IExpr](),
		}
	}
	if t.value == ")" {
		return &EndGroupExpr{
			Expr: Expr{
				token: t,
			},
		}
	}
	iValue := strings.ToLower(t.value)
	if iValue == "and" {
		return &LogicExpr{
			Expr: Expr{
				token: t,
			},
			logic: AND,
		}
	}

	if iValue == "or" {
		return &LogicExpr{
			Expr: Expr{
				token: t,
			},
			logic: OR,
		}
	}

	return &SimpleExpr{
		Expr: Expr{
			token: t,
		},
		expr: t.value,
	}
}

func Parse(expr string) (IExpr, error) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return nil, errors.New("empty expr")
	}

	exprs := []IExpr{}

	var preExpr IExpr

	stack := NewStack[IExpr]()

	tokenizer := Tokenizer{expr: expr}

	for token := range tokenizer.Parse() {
		ex := createExprFromToken(token)
		//fmt.Printf(">>>>>>> %s\n", token)

		if preExpr != nil {
			switch preExpr.(type) {
			case *StartGroupExpr:
				switch ex.(type) {
				case *StartGroupExpr:
				case *SimpleExpr:
				default:
					return nil, errors.New(fmt.Sprintf("expect simple expr, but got %s", token))
				}
			case *EndGroupExpr:
				switch ex.(type) {
				case *EndGroupExpr:
				case *LogicExpr:
				default:
					return nil, errors.New(fmt.Sprintf("expect logic expr, but got %s", token))
				}
			case *LogicExpr:
				switch ex.(type) {
				case *StartGroupExpr:
				case *SimpleExpr:
				default:
					return nil, errors.New(fmt.Sprintf("expect logic or start group expr, but got %s", token))
				}
			case *SimpleExpr:
				switch ex.(type) {
				case *EndGroupExpr:
				case *LogicExpr:
				default:
					return nil, errors.New(fmt.Sprintf("expect logic or end group expr, but got %s", token))
				}
			}
		} else {
			switch ex.(type) {
			case *StartGroupExpr:
			case *SimpleExpr:
			default:
				return nil, errors.New(fmt.Sprintf("expect start group expr or simple expr, but got %s", token))
			}
		}

		switch cv := ex.(type) {
		case *StartGroupExpr:
			stack.Push(ex)
		case *EndGroupExpr:
			if stack.IsEmpty() {
				return nil, errors.New(fmt.Sprintf("single close group expr: %s", token))
			}
			group := stack.Pop()
			switch v := (group).(type) {
			case *StartGroupExpr:
				if v.stack.Size() == 0 {
					return nil, errors.New(fmt.Sprintf("empty group: %s", token))
				} else if v.stack.Size() > 1 {
					return nil, errors.New(fmt.Sprintf("too many expr in group's stack: %s", token))
				}
				v.exprs = append(v.exprs, v.stack.Pop())

				v.complete = true

				if stack.Size() > 0 {
					top := stack.Top()
					switch logic := top.(type) {
					case *LogicExpr:
						logic.right = group
					}
				} else {
					stack.Push(group)
				}
			default:
				return nil, errors.New(fmt.Sprintf("non-pairs group: %s", token))
			}
		case *LogicExpr:
			top := stack.Top()
			switch v := top.(type) {
			case *StartGroupExpr:
				if v.complete {
					cv.left = stack.Pop()
					stack.Push(ex)
				} else if v.stack.Size() >= 1 {
					cv.left = v.stack.Pop()
					v.stack.Push(ex)
				} else {
					return nil, errors.New(fmt.Sprintf("can not go here"))
				}
			case *SimpleExpr:
				cv.left = stack.Pop()
				stack.Push(ex)
			case *LogicExpr:
				cv.left = stack.Pop()
				stack.Push(ex)
			}
		case *SimpleExpr:

			f, err := ParseFilterString(token.value)
			if err != nil {
				return nil, err
			}
			cv.filter = f

			if !stack.IsEmpty() { // group or logic
				group := stack.Top()
				switch v := group.(type) {
				case *StartGroupExpr: // incomplete-group
					if v.stack.Size() == 0 {
						v.stack.Push(ex)
					} else {
						last := v.stack.Pop()
						switch logic := last.(type) {
						case *LogicExpr:
							logic.right = ex
							v.stack.Push(logic)
						}
					}
				case *LogicExpr:
					if stack.Size() < 1 {
						return nil, errors.New(fmt.Sprintf("not enough expr for logic: %s", token))
					}
					logic := stack.Pop() // group itself
					//simple := stack.Pop()
					switch v := logic.(type) {
					case *LogicExpr:
						//v.left = simple
						v.right = ex
						stack.Push(logic)
					}
				default:
					return nil, errors.New(fmt.Sprintf("can not go here: %s", token))
				}
			} else {
				stack.Push(ex)
				//exprs = append(exprs, ex)
			}
		}

		preExpr = ex
	}

	if !stack.IsEmpty() {
		e := stack.Pop()
		switch v := e.(type) {
		case *StartGroupExpr:
			exprs = append(exprs, e)
			//return nil, errors.New(fmt.Sprintf("non-pairs group: %s", v.token))
		case *LogicExpr:
			exprs = append(exprs, e)
		case *SimpleExpr:
			exprs = append(exprs, e)
		default:
			return nil, errors.New(fmt.Sprintf("can not go here: %s", v.Strings()))
		}
	}

	if len(exprs) == 0 {
		return nil, errors.New("empty exprs")
	}

	if len(exprs) != 1 {
		return nil, errors.New("invalid expr, can not go here")
	}

	last := exprs[0]
	switch last.(type) {
	case *SimpleExpr:
	case *LogicExpr:
	case *StartGroupExpr:
	default:
		return nil, errors.New(fmt.Sprintf("expect group or simple or logic expr at the end, but got: %s", last))
	}

	return last, nil
}
