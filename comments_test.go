package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	singleline Comment = "Test"
	multiline  Comment = "Test\n Line 2"
)

func TestToSingleLine(t *testing.T) {
	require.Equal(t, "Test", singleline.ToSingleLine())
	require.Equal(t, "Test Line 2", multiline.ToSingleLine())
}
