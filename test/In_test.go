// +build test
package generationk

import (
	genk "generationk/internal"
	"testing"
)

func TestInterface(t *testing.T) {
	var xInt genk.IntTest
	var xImpl genk.IntImpl
	xInt = xImpl
	xInt.Update()
	xImpl.NotInInterface()
}
