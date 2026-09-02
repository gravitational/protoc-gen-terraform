package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToSingleLine(t *testing.T) {
	test := func(input, expected string) func(t *testing.T) {
		return func(t *testing.T) {
			require.Equal(t, expected, Comment(input).ToSingleLine())
		}
	}

	t.Run("Single Line", test("Test", "Test"))
	t.Run("Multiline Line", test("Test\n Line 2", "Test Line 2"))
	t.Run("Empty String", test("", ""))
	t.Run("CRLF", test("Test\r\n Line 2", "Test Line 2"))
	t.Run("Leading", test("\nLine 1\n Line 2", "Line 1 Line 2"))
	t.Run("Trailing", test("Line 1\n Line 2\n", "Line 1 Line 2"))
	t.Run("Trailing", test("Line 1\n   Line 2\n", "Line 1 Line 2"))
}
