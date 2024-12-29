package job

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUniqueSet(t *testing.T) {
	set := make(UniqueSet)
	set.Insert(5, 7, 9)
	assert.True(t, set.Contain(5))
	assert.True(t, set.Contain(7))
	assert.True(t, set.Contain(9))
	assert.True(t, set.Contain(5, 7, 9))

	containFive := false
	set.Range(func(i UniqueID) bool {
		if i == 5 {
			containFive = true
			return false
		}
		return true
	})
	assert.True(t, containFive)

	set.Remove(7)
	assert.True(t, set.Contain(5))
	assert.False(t, set.Contain(7))
	assert.True(t, set.Contain(9))
	assert.False(t, set.Contain(5, 7, 9))

	count := 0
	set.Range(func(element UniqueID) bool {
		count++
		return true
	})
	assert.Equal(t, set.Len(), count)
	count = 0
	set.Range(func(element UniqueID) bool {
		count++
		return false
	})
	assert.Equal(t, 1, count)
}

func TestUniqueSetClear(t *testing.T) {
	set := make(UniqueSet)
	set.Insert(5, 7, 9)
	assert.True(t, set.Contain(5))
	assert.True(t, set.Contain(7))
	assert.True(t, set.Contain(9))
	assert.Equal(t, 3, set.Len())

	set.Clear()
	assert.False(t, set.Contain(5))
	assert.False(t, set.Contain(7))
	assert.False(t, set.Contain(9))
	assert.Equal(t, 0, set.Len())
}

func TestConcurrentSet(t *testing.T) {
	set := NewConcurrentSet[int]()
	set.Insert(1)
	set.Insert(2)
	set.Insert(3)

	count := 0
	set.Range(func(element int) bool {
		count++
		return true
	})

	assert.Len(t, set.Collect(), count)
}
