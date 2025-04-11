package main

import (
	"testing"
	"time"

	"github.com/mishuma1/pokemon/cache"
)

func TestCleanInput(t *testing.T) {
	// ...

	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello     world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "",
			expected: []string{},
		},
		{
			input:    "BoB",
			expected: []string{"bob"},
		},
		{
			input:    "This is a    new day   is it",
			expected: []string{"this", "is", "a", "new", "day", "is", "it"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Length Expected: %v, Actual: %v", len(c.expected), len(actual))
			continue
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("Expected: %v, Actual: %v", expectedWord, word)
			}
		}
	}

	var restTime int64 = 2
	cacheTest := cache.Cache{}
	cacheTest.NewCache(restTime, restTime)
	key := "test"
	val := []byte("This is a test")
	cacheTest.AddCache(key, val)
	resultSize := cacheTest.CacheSize()
	if resultSize == 0 {
		t.Errorf("Cache should not be empty")
	}
	lookup := cacheTest.GetCache(key)
	if len(lookup) == 0 {
		t.Errorf("Cache result should not be empty")
	}
	time.Sleep(time.Second * time.Duration(restTime+2))
	resultSize = cacheTest.CacheSize()
	if resultSize != 0 {
		t.Errorf("Cache should  be empty")
	}
}
