package haxxmap_test

import (
	"fmt"
	"strings"

	"github.com/aezhar/haxxmap"
)

// your custom hash function
// the hash function signature must adhere to `func(keyType) uintptr`
func customStringHasher(s string) uintptr {
	return uintptr(len(s))
}

// your custom comparison function
// This allows for more complex key types to be compared
func customStringCompare(l, r string) bool {
	return strings.ToLower(l) == strings.ToLower(r)
}

func Example() {
	m := haxxmap.New[string, string]()   // initialize a string-string map
	m.SetHasher(customStringHasher)      // this overrides the default xxHash algorithm
	m.SetComparator(customStringCompare) // this overrides the default key comparison function

	m.Set("one", "1")
	m.Set("Two", "2")
	m.Set("three", "3")

	if val, ok := m.Get("One"); ok {
		fmt.Println(val)
	}

	if val, ok := m.Get("two"); ok {
		fmt.Println(val)
	}

	if val, ok := m.Get("three"); ok {
		fmt.Println(val)
	}
	// Output:
	// 1
	// 2
	// 3
}
