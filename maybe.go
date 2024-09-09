// Package maybe provides a sum type inspired by Haskell's Maybe monad.
// Instances can by associated with a value or be empty.
// The [Unit] function is used to create an instance with a value.
// The [Nothing] function is used to create an instance without an associated value.
package maybe

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"reflect"
	"slices"
	"unsafe"
)

var (
	null = []byte("null")
)

// Identity is the identity function.
func Identity[T any](x T) T {
	return x
}

// Iface is a non-generic interface for [Maybe].
type Iface interface {
	Set(interface{})
	SetPtr(unsafe.Pointer)
	Get() (interface{}, bool)
	GetPtr() unsafe.Pointer
	SetValid()
	MaybeType() reflect.Type
}

// Value is a value interface for a maybe type.
type Value interface {
	MaybeIface() Iface
}

// MaybeIface returns an instance of [Iface].
func (m Maybe[T]) MaybeIface() Iface {
	return &m
}

// Maybe is a maybe type.
type Maybe[T any] struct {
	Val   T
	Valid bool
}

var (
	// IfaceType is the reflected type of [Iface].
	IfaceType       = reflect.TypeFor[Iface]()
	_         Iface = (*Maybe[int])(nil)
)

// New creates a new Maybe instance from a pointer.
func New[T any](x *T) Maybe[T] {
	if x == nil {
		return Nothing[T]()
	}
	return Unit(*x)
}

// MaybeType returns the underlying type.
func (m *Maybe[T]) MaybeType() reflect.Type {
	return reflect.TypeFor[T]()
}

// GetOr returns the underlying value if valid and `defVal` otherwise.
func (m Maybe[T]) GetOr(defVal T) T {
	if m.Valid {
		return m.Val
	}
	return defVal
}

// GetOrZero returns the underlying value if valid and the zero value otherwise.
func (m Maybe[T]) GetOrZero() T {
	if m.Valid {
		return m.Val
	}
	var x T
	return x
}

// Get gets the underlying value if valid.
func (m *Maybe[T]) Get() (interface{}, bool) {
	if m.Valid {
		return m.Val, true
	}
	return nil, false
}

// Pointer gets the pointer to the underlying value or nil in case there's none.
func (m Maybe[T]) Pointer() *T {
	if m.Valid {
		return &m.Val
	}
	return nil
}

// GetPtr gets the unsafe pointer to the underlying value or nil in case there's none.
func (m *Maybe[T]) GetPtr() unsafe.Pointer {
	return unsafe.Pointer(m.Pointer())
}

// Set sets the underlying value.
func (m *Maybe[T]) Set(x interface{}) {
	m.Val = x.(T)
	m.Valid = true
}

// SetPtr sets the underlying value from an unsafe pointer.
func (m *Maybe[T]) SetPtr(ptr unsafe.Pointer) {
	if ptr == nil {
		m.Valid = false
	} else {
		m.Valid = true
		m.Val = *(*T)(ptr)
	}
}

// SetValid sets the underlying value.
func (m *Maybe[T]) SetValid() {
	m.Valid = true
}

// Unit returns a maybe instance with an underlying value.
func Unit[T any](x T) Maybe[T] {
	return Maybe[T]{Val: x, Valid: true}
}

// Nothing returns a maybe instance representing nothing.
func Nothing[T any]() Maybe[T] {
	return Maybe[T]{}
}

// Fmap is the functorial map for Maybe.
func Fmap[T, U any](f func(T) U, x Maybe[T]) Maybe[U] {
	if !x.Valid {
		return Maybe[U]{}
	}
	return Maybe[U]{Valid: true, Val: f(x.Val)}
}

// FallibleFmap is the functorial map for a possibly erring function.
func FallibleFmap[T, U any](f func(T) (U, error), x Maybe[T]) (Maybe[U], error) {
	if !x.Valid {
		return Maybe[U]{}, nil
	}
	y, err := f(x.Val)
	if err != nil {
		return Maybe[U]{}, err
	}
	return Maybe[U]{Valid: true, Val: y}, nil
}

// Bind is the monadic bind operation.
func Bind[T, U any](f func(T) Maybe[U], x Maybe[T]) Maybe[U] {
	if !x.Valid {
		return Maybe[U]{}
	}
	return f(x.Val)
}

// Join is the monadic join operation.
func Join[T any](x Maybe[Maybe[T]]) Maybe[T] {
	return Bind(Identity[Maybe[T]], x)
}

func (m Maybe[T]) MarshalJSON() ([]byte, error) {
	if !m.Valid {
		return null, nil
	}
	return json.Marshal(m.Val)
}

func (m *Maybe[T]) UnmarshalJSON(val []byte) error {
	if slices.Equal(val, null) {
		return nil
	}
	m.Valid = true
	return json.Unmarshal(val, &m.Val)
}

func (m *Maybe[T]) Scan(val any) error {
	var v sql.Null[T]
	if err := v.Scan(val); err != nil {
		return err
	}
	m.Valid = v.Valid
	m.Val = v.V
	return nil
}

func (m Maybe[T]) Value() (driver.Value, error) {
	if !m.Valid {
		return nil, nil
	}

	iface := interface{}(m.Val)

	switch v := iface.(type) {
	case driver.Valuer:
		return v.Value()

	// for numbers only int64 and float64 is supported https://pkg.go.dev/database/sql/driver@go1.22.0#Value

	case int:
		return int64(v), nil
	case uint:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint64:
		return int64(v), nil
	case float32:
		return float64(v), nil
	}

	return m.Val, nil
}

var (
	_ json.Marshaler   = Unit(0)
	_ json.Unmarshaler = (*Maybe[int])(nil)
	_ driver.Valuer    = Unit(0)
	_ sql.Scanner      = (*Maybe[int])(nil)
)
