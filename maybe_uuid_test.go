//go:build go1.27

package maybe

import (
	"database/sql/driver"
	"testing"
	"uuid"

	"github.com/stretchr/testify/require"
)

func TestValueUUID(t *testing.T) {
	t.Run("uuid is converted to its string form", func(t *testing.T) {
		req := require.New(t)

		id := uuid.MustParse("0a65a93e-1ba8-48d7-a149-6350d05acdbb")
		v, err := Unit(id).Value()
		req.NoError(err)
		req.Equal("0a65a93e-1ba8-48d7-a149-6350d05acdbb", v)
	})

	t.Run("uuid passes the database/sql argument converter", func(t *testing.T) {
		req := require.New(t)

		id := uuid.MustParse("0a65a93e-1ba8-48d7-a149-6350d05acdbb")
		v, err := driver.DefaultParameterConverter.ConvertValue(Unit(id))
		req.NoError(err)
		req.Equal("0a65a93e-1ba8-48d7-a149-6350d05acdbb", v)
	})

	t.Run("nothing passes the database/sql argument converter as nil", func(t *testing.T) {
		req := require.New(t)

		v, err := driver.DefaultParameterConverter.ConvertValue(Nothing[uuid.UUID]())
		req.NoError(err)
		req.Nil(v)
	})
}
