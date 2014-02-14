Gloom - A Go (GoLang) Bloom Filter
-----

A bloom filter implemented in Go.

Author info: @rtt / rsty.org / github.com/rtt

Further reading: http://en.wikipedia.org/wiki/Bloom_filter

Usage
-----

Example code:

````Go
import (
  "github.com/rtt/gloom"
  "fmt"
)

func main() {
  // init a new filter with 0 initial items, a 32bit bitset, and 3 hash functions
  filter, _ := gloom.New(0, 32, 3)

  // add a value to the filter
  filter.Add("test value")

  // test a value
  if filter.Test("some value") {
    fmt.Println("hit! (probably :))")
  } else {
    fmt.Println("miss!")
  }

}

````


License
-----

Public domain. Do as you wish with it. Pull requests (etc) welcome!


