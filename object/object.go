package object

import "fmt"

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
