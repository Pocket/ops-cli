package util

import (
	"testing"

	"gotest.tools/assert"
)

func TestStringInSlice(t *testing.T) {
	strings := []string{"billy", "joel", "New York State of Mind", "Piano Man", "You May Be Right"}
	assert.Assert(t, StringInSlice("billy", strings))
	assert.Assert(t, StringInSlice("joel", strings))
	assert.Assert(t, !StringInSlice("Piano Ma", strings))
	assert.Assert(t, StringInSlice("You May Be Right", strings))
	assert.Assert(t, StringInSlice("Piano Man", strings))
	assert.Assert(t, !StringInSlice("laslsad", strings))
}

func TestRemoveDuplicatesFromSlice(t *testing.T) {
	strings := []string{"billy", "billy", "billy", "New York State of Mind", "Piano Man", "billy", "You May Be Right", "Piano Man"}
	assert.DeepEqual(t, RemoveDuplicatesFromSlice(strings), []string{"billy", "New York State of Mind", "Piano Man", "You May Be Right"})
}
