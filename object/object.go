package object

import "fmt"

var (
	NULL  = &Null{}
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (_ *Integer) Type() ObjectType {
	return OBJ_INTEGER
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

func (_ *Boolean) Type() ObjectType {
	return OBJ_BOOLEAN
}

type Null struct{}

func (_ *Null) Inspect() string {
	return "null"
}

func (_ *Null) Type() ObjectType {
	return OBJ_NULL
}

func AsInt(val int64) *Integer {
	return &Integer{Value: val}
}

func AsBool(val bool) *Boolean {
	if val {
		return TRUE
	}
	return FALSE
}

func IsTruthy(obj Object) bool {
	switch obj {
	case FALSE:
		fallthrough
	case NULL:
		return false
	default:
		return true
	}
}

type ReturnValue struct {
	Value Object
}

func (_ *ReturnValue) Type() ObjectType {
	return OBJ_RETURN_VALUE
}

func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

type Error struct {
	Msg string
}

func (_ *Error) Type() ObjectType {
	return OBJ_ERROR
}

func (e *Error) Inspect() string {
	return "ERROR: " + e.Msg
}

func FormatError(format string, args ...any) *Error {
	return &Error{
		Msg: fmt.Sprintf(format, args...),
	}
}

func IsError(obj Object) bool {
	return obj != nil && obj.Type() == OBJ_ERROR
}
