# HaxxMap

[![Main Actions Status](https://github.com/aezhar/haxxmap/workflows/Go/badge.svg)](https://github.com/aezhar/haxxmap/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/aezhar/haxxmap)](https://goreportcard.com/report/github.com/aezhar/haxxmap)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](./LICENSE.md)
> A lightning fast concurrent hashmap

This is a fork of the original [haxmap](https://github.com/alphadose/haxmap) package developed by [alphadose](https://github.com/alphadose). The goal of this fork is to allow for more customization of the hashmap's behavior at the expense of execution time.  

The default hashing algorithm for strings used is [xxHash](https://github.com/Cyan4973/xxHash) and the hashmap's buckets are implemented using [Harris lock-free list](https://www.cl.cam.ac.uk/research/srg/netos/papers/2001-caslists.pdf)

## Installation

You need Golang [1.18.x](https://go.dev/dl/) or above

```bash
$ go get github.com/aezhar/haxxmap
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/aezhar/haxxmap"
)

func main() {
	// initialize map with key type `int` and value type `string`
	mep := haxxmap.New[int, string]()

	// set a value (overwrites existing value if present)
	mep.Set(1, "one")

	// get the value and print it
	val, ok := mep.Get(1)
	if ok {
		println(val)
	}

	mep.Set(2, "two")
	mep.Set(3, "three")
	mep.Set(4, "four")

	// ForEach loop to iterate over all key-value pairs and execute the given lambda
	mep.ForEach(func(key int, value string) bool {
		fmt.Printf("Key -> %d | Value -> %s\n", key, value)
		return true // return `true` to continue iteration and `false` to break iteration
	})

	mep.Del(1) // delete a value
	mep.Del(0) // delete is safe even if a key doesn't exists

	// bulk deletion is supported too in the same API call
	// has better performance than deleting keys one by one
	mep.Del(2, 3, 4)

	fmt.Println("entries left: ", mep.Len())
}
```

## Benchmarks

Benchmarks are performed against other implementations of thread-safe hashmaps:
* [sync.Map](https://pkg.go.dev/sync#Map)
* [github.com/cornelk/hashmap](https://github.com/cornelk/hashmap)
* [github.com/puzpuzpuz/xsync](https://github.com/puzpuzpuz/xsync)
* [github.com/alphadose/haxmap](https://github.com/alphadose/haxmap)

All results are computed from [benchstat](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat) of 20 runs (code available [here](./benchmarks))

1. Concurrent Reads Only
```
cpu: AMD Ryzen 7 5800X 8-Core Processor
                       │   sec/op    │
HaxMapReadsOnly-16       2.627µ ± 2%
HaxxMapReadsOnly-16      3.850µ ± 0%
GoSyncMapReadsOnly-16    10.28µ ± 1%
CornelkMapReadsOnly-16   2.997µ ± 2%
XsyncMapReadsOnly-16     1.929µ ± 1%
```

2. Concurrent Reads with Writes
```
cpu: AMD Ryzen 7 5800X 8-Core Processor
                             │   sec/op    │
HaxMapReadsWithWrites-16       3.216µ ± 7%
HaxxMapReadsWithWrites-16      3.778µ ± 4%
GoSyncMapReadsWithWrites-16    11.74µ ± 1%
CornelkMapReadsWithWrites-16   3.545µ ± 5%
XsyncMapReadsWithWrites-16     2.373µ ± 1%

                             │     B/op      │
HaxMapReadsWithWrites-16         339.5 ± 2%
HaxxMapReadsWithWrites-16        444.5 ± 2%
GoSyncMapReadsWithWrites-16    2.408Ki ± 3%
CornelkMapReadsWithWrites-16     398.5 ± 2%
XsyncMapReadsWithWrites-16       396.5 ± 3%

                             │  allocs/op  │
HaxMapReadsWithWrites-16       42.00 ± 2%
HaxxMapReadsWithWrites-16      55.00 ± 2%
GoSyncMapReadsWithWrites-16    228.5 ± 3%
CornelkMapReadsWithWrites-16   49.00 ± 2%
XsyncMapReadsWithWrites-16     24.00 ± 4%

```

## Tips

1. HaxxMap by default uses [xxHash](https://github.com/cespare/xxhash) algorithm and compares each value directly, but this behavior can be overriden by specifying a different hash and comparison function before using the hashmap. Beneath lies an example for the same.

```go
package main

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

func main() {
	// initialize a string-string map
	m := haxxmap.New[string, string](
		// this overrides the default xxHash algorithm
		haxxmap.WithHasher[string, string](customStringHasher),
		// this overrides the default key comparison function
		haxxmap.WithComparator[string, string](customStringCompare),
	)

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
}
```

2. You can pre-allocate the size of the map which will improve performance in some cases.

```go
package main

import (
	"fmt"

	"github.com/aezhar/haxxmap"
)

func main() {
	const initialSize = 1 << 10

	// pre-allocating the size of the map will prevent all grow operations
	// until that limit is hit thereby improving performance
	m := haxxmap.New[int, string](haxxmap.WithInitialSize[int, string](initialSize))

	m.Set(1, "1")

	if val, ok := m.Get(1); ok {
		fmt.Println(val)
	}
}
```
