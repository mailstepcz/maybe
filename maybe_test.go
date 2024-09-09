package maybe

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/mailstepcz/pointer"
	"github.com/stretchr/testify/require"
)

type s struct {
	N Maybe[int] `json:"n"`
}

func TestUnit(t *testing.T) {
	req := require.New(t)

	m := Unit(1234)
	req.Equal(1234, m.Val)
	req.Equal(true, m.Valid)
}

func TestNothing(t *testing.T) {
	req := require.New(t)

	m := Nothing[int]()
	req.Equal(false, m.Valid)
}

func TestMarshal(t *testing.T) {
	req := require.New(t)

	b, err := json.Marshal(s{})
	req.NoError(err)
	req.Equal([]byte(`{"n":null}`), b)

	b, err = json.Marshal(s{N: Unit(1234)})
	req.NoError(err)
	req.Equal([]byte(`{"n":1234}`), b)

	var s s
	err = json.Unmarshal([]byte(`{"n":null}`), &s)
	req.NoError(err)
	req.Equal(Maybe[int]{}, s.N)

	err = json.Unmarshal([]byte(`{"n":1234}`), &s)
	req.NoError(err)
	req.Equal(Unit(1234), s.N)
}

func TestFmap(t *testing.T) {
	req := require.New(t)

	req.Equal(Unit("1234"), Fmap(func(x int) string { return strconv.Itoa(x) }, Unit(1234)))

	x, err := FallibleFmap(func(x int) (string, error) { return strconv.Itoa(x), nil }, Unit(1234))
	req.NoError(err)
	req.Equal(Unit("1234"), x)

	_, err = FallibleFmap(func(x int) (string, error) { return "", errors.ErrUnsupported }, Unit(1234))
	req.Error(err)
}

func TestBind(t *testing.T) {
	req := require.New(t)

	req.Equal(Unit(1234), Join(Unit(Unit(1234))))
}

func TestGetOr(t *testing.T) {
	req := require.New(t)

	req.Equal(1234, Unit(1234).GetOr(5678))
	req.Equal(5678, Nothing[int]().GetOr(5678))
}

func TestNew(t *testing.T) {
	req := require.New(t)

	req.Equal(Unit(1234), New(pointer.To(1234)))
	req.Equal(Nothing[int](), New[int](nil))
}

func ExampleUnit() {
	m := Unit(1234)
	fmt.Println(m)
	// Output: {1234 true}
}

func ExampleNothing() {
	m := Nothing[int]()
	fmt.Println(m)
	// Output: {0 false}
}
