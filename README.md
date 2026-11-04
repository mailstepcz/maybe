# maybe

Generic optional values for Go (`Maybe[T]`), with JSON and `database/sql` support.

## Install

```sh
go get github.com/mailstepcz/maybe
```

Requires Go 1.27 or newer.

## Usage

```go
package main

import (
	"fmt"

	"github.com/mailstepcz/maybe"
)

func main() {
	a := maybe.Unit(42)
	b := maybe.Nothing[int]()

	fmt.Println(a.Valid, a.Val)               // true 42
	fmt.Println(b.GetOr(-1))                  // -1
	fmt.Println(b.GetOrZero())                // 0
	fmt.Println(b.Pointer())                  // <nil>
	fmt.Println(*maybe.New(new(7)).Pointer()) // 7
}
```

## API

`Maybe[T]` is a struct with the exported fields `Val T` and `Valid bool`.

Constructors

- `Unit(x T) Maybe[T]` — a present value.
- `Nothing[T]() Maybe[T]` — no value (same as the zero `Maybe[T]`).
- `New(x *T) Maybe[T]` — `Nothing` for a nil pointer, otherwise `Unit(*x)`.

Reading

- `GetOr(defVal T) T` — the value, or `defVal` if empty.
- `GetOrZero() T` — the value, or the zero value of `T` if empty.
- `Pointer() *T` — a pointer to a copy of the value, or `nil` if empty.

Combinators

- `Fmap(f func(T) U, x Maybe[T]) Maybe[U]` — applies `f` to a present value.
- `FallibleFmap(f func(T) (U, error), x Maybe[T]) (Maybe[U], error)` — like `Fmap`; returns the error of `f`.
- `Bind(f func(T) Maybe[U], x Maybe[T]) Maybe[U]` — applies `f` to a present value and returns its result.
- `Join(x Maybe[Maybe[T]]) Maybe[T]` — flattens a nested `Maybe`.
- `Identity(x T) T` — the identity function.

Reflection / non-generic access (`Iface`, implemented by `*Maybe[T]`)

- `Get() (any, bool)`, `Set(any)`, `GetPtr() unsafe.Pointer`, `SetPtr(unsafe.Pointer)`, `SetValid()`, `MaybeType() reflect.Type`.
- `Value` interface (`MaybeIface() Iface`), implemented by `Maybe[T]`; `IfaceType` is the `reflect.Type` of `Iface`.

## JSON

- `Nothing` marshals to `null`; a present value marshals as `T` would.
- Unmarshaling `null` leaves the `Maybe` unchanged (a zero `Maybe` stays `Nothing`); any other input sets `Valid` and decodes into `Val`.

```go
type payload struct {
	N maybe.Maybe[int] `json:"n"`
}

b, _ := json.Marshal(payload{})                   // {"n":null}
b, _ = json.Marshal(payload{N: maybe.Unit(1234)}) // {"n":1234}

var p payload
_ = json.Unmarshal([]byte(`{"n":1234}`), &p) // p.N == maybe.Unit(1234)
```

## SQL

`*Maybe[T]` implements `sql.Scanner` and `Maybe[T]` implements `driver.Valuer`, so it can be used directly as a query argument and scan target.

- `Scan`: `NULL` becomes `Nothing`; any other value is converted into `T` the same way `sql.Null[T]` does.
- `Value`: `Nothing` is `NULL`. A present value is resolved in this order:
  1. `T` implements `driver.Valuer` — its `Value()` is used.
  2. Built-in integer types become `int64`, `float32` becomes `float64`.
  3. A value that already is a `driver.Value` is returned as is.
  4. Otherwise `driver.DefaultParameterConverter` converts it — e.g. `uuid.UUID` (stdlib) becomes its string form, named string/int types become `string`/`int64`.
  5. If the converter fails, the raw value is returned.

```go
id := uuid.MustParse("0a65a93e-1ba8-48d7-a149-6350d05acdbb")

v, _ := maybe.Unit(id).Value()            // "0a65a93e-1ba8-48d7-a149-6350d05acdbb"
v, _ = maybe.Nothing[uuid.UUID]().Value() // nil (NULL)
```
