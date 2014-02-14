Gloom - A Go (GoLang) Bloom Filter
-----

A bloom filter implemented in Go. There are others, but I just fancied building one for the sake of it. It uses [MurmurHash](http://en.wikipedia.org/wiki/MurmurHash) 32 bit hashing, and also supports removing of values from the set via bit counting. 

Author info: [@rtt](http://twitter.com/rtt) / [rsty.org](http://rsty.org)

[Further reading about bloom filters on wikipedia](http://en.wikipedia.org/wiki/Bloom_filter)

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
  if filter, err := gloom.New(0, 32, 3); err != nil {
    fmt.Printf("error", err)
    return
  }

  v, t := "some value", "some test"
  
  // add a value to the filter
  filter.Add(v)

  // test a value
  if filter.Test(t) {
    fmt.Println("hit! (probably :))")
  } else {
    fmt.Println("miss!")
  }
  
  // remove it
  filter.Remove(v)
  
}

````


License
-----

Public domain. Do as you wish with it. Pull requests (etc) welcome!


