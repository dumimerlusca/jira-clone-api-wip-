package random

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomInt(t *testing.T) {
	min := int64(0)
	max := int64(223823)

	n1 := RandomInt(min, max)
	n2 := RandomInt(min, max)
	n3 := RandomInt(min, max)

	require.True(t, n1 >= min && n1 <= max)
	require.True(t, n2 >= min && n2 <= max)
	require.True(t, n3 >= min && n3 <= max)

	require.True(t, n1 != n2 && n1 != n3 && n2 != n3)
}

func TestRandomString(t *testing.T) {
	length := 20
	str1 := RandomString(length)
	str2 := RandomString(length)

	require.True(t, len(str1) == length)
	require.NotEqual(t, str1, str2)
}
