package desc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSet(t *testing.T) {
	s := NewImports()

	s.SetDefaultQual("test", "test")

	s.AddQual("ext", "github.com/test/ext")

	require.Equal(t, "string", s.GoString("string", true))
	require.Equal(t, "[]string", s.GoString("[]string", true))
	require.Equal(t, "[]test.Type", s.GoString("[]Type", true))
	require.Equal(t, "[]*test.Type", s.GoString("[]*Type", true))
	require.Equal(t, "map[string]test.Test", s.GoString("map[string]Test", true))
	require.Equal(t, "map[*string]*test.Test", s.GoString("map[*string]*Test", true))
	require.Equal(t, "[]*test.Type", s.GoString("[]*Type", true))
	require.Equal(t, "[]*alfa.Type", s.GoString("[]*alfa.Type", true))
	require.Equal(t, "map[test.Key]test.Value", s.GoString("map[test.Key]Value", true))

	require.Equal(t, "ext.SomeType", s.GoString("github.com/test/ext.SomeType", false))
	require.Equal(t, "ext.SomeType", s.GoString("ext.SomeType", false))

	require.Equal(t, "other_ext.SomeType", s.GoString("other_ext.SomeType", false))
}
