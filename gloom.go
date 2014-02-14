package main

/*
 * Author Rich Taylor, 2014
 * twitter.com/rtt  -- http://rsty.org
 * License: Public domain
 */

import (
    "bytes"
    "encoding/binary"
    "fmt"
)

// this could also be called a bit vector?
type BitSet struct {
    sz uint
    bits []uint
}

type BloomFilter struct {
    c uint // count of items in filter
    m uint // size of bitset (in bits)
    k uint // amt of hashes
    b BitSet
}

/*
 * Implements a 32-bit murmur hash
 * This is a straight implementation from:
 *    http://en.wikipedia.org/wiki/MurmurHash
 */
func Murmur32(key []byte) uint32 {
    length := len(key)
    if length == 0 {
        return 0
    }
    c1, c2 := uint32(0xcc9e2d51), uint32(0x1b873593)
    blocks := length / 4
    var h, k uint32
    buf := bytes.NewBuffer(key)
    for i := 0; i < blocks; i++ {
        binary.Read(buf, binary.LittleEndian, &k)
        k *= c1
        k = (k << 15) | (k >> (32 - 15))
        k *= c2
        h ^= k
        h = (h << 13) | (h >> (32 - 13))
        h = (h * 5) + 0xe6546b64
    }
    k = 0
    remaining := blocks * 4
    switch length & 3 {
        case 3:
            k ^= uint32(key[remaining + 2]) << 16
            fallthrough
        case 2:
            k ^= uint32(key[remaining + 1]) << 8
            fallthrough
        case 1:
            k ^= uint32(key[remaining])
            k *= c1
            k = (k << 15) | (k >> (32 - 15))
            k *= c2
            h ^= k
    }
    h ^= uint32(length)
    h ^= h >> 16
    h *= 0x85ebca6b
    h ^= h >> 13
    h *= 0xc2b2ae35
    h ^= h >> 16
    return h
}


/*
 * Initialises a BitSet to a given length
 * (all bits set to zero)
 */
func NewBitSet(l int) BitSet {
    b := BitSet{uint(l), make([]uint, l, l)}
    return b
}

/*
 * Adds a uint64 value into a bitset
 */
func (b *BitSet) SetVal(val uint32) error {
    x, i := uint32(1), uint(0)
    for ; i < b.sz; i++ {
        if (val & x) == x {
            b.bits[i] = 1
        }
        x = x << 1
    }
    return nil; // todo!
}

/*
 * Tests whether a uint64 lies in the bitset
 */

func (b *BitSet) TestVal(val uint32) (bool, error) {
    results := make([]bool, 0)

    x, i := uint32(1), uint(0)
    for ; i < b.sz; i++ {
        // if this bit is meant to be on, record whether it was or not
        if (val & x) == x {
            results = append(results, b.bits[i] == 1)
        }
        x = x << 1
    }

    return all(results), nil
}

/*
 * Adds a value (item) to a BloomFilter
 */
func (bf *BloomFilter) Add(item string) {

    // get a hash for each item
    hashes := bf.Digest(item)

    // add each hash into the bf's bitset
    for i := range(hashes) {
        bf.b.SetVal(hashes[i])
    }

    // increment the count
    bf.c += 1
}

/*
 * Tests whether an item is in a bloom filter
 */
func (bf BloomFilter) Test(item string) bool {
    hashes := bf.Digest(item)
    results := make([]bool, bf.k)
    for i := range(hashes) {
        r, _ := bf.b.TestVal(hashes[i])
        results[i] = r
    }
    return all(results)
}

/* Nicely tells us about a BloomFilter
 */
func (bf BloomFilter) String() string {
    return fmt.Sprintf(
        "BloomFilter with %d filters, %d items and bitset size of %d",
        bf.k, bf.c, bf.m)
}

/*
 * Returns a digested value for all of the filters
 * in a BloomFilter
 */
func (bf BloomFilter) Digest(value string) []uint32 {
    s, i := make([]uint32, 0, bf.k), uint32(0)
    for ; uint(i) < bf.k; i++ {
        vl := len(value)
        h := Murmur32(
            []byte(value[0:vl / 2])) +
            (i * Murmur32([]byte(value[vl / 2:]))) % (1 << bf.m - 1)

        s = append(s, h)
    }
    return s
}

/*
 * returns true if all items in the supplied []bool are true, otherwise false
 */
func all(b []bool) bool {
    l, c, i := len(b), 0, 0
    for ; i < l; i++ {
        if b[i] {
            c++ // lol
        }
    }
    return l == c
}

/* Convenience function for constructing a new
 * bloom filter
 */
func New(c int, m int, k int) (BloomFilter, error) {
    if m > 32 {
        return BloomFilter{}, fmt.Errorf("Bits must not be greater than 32")
    }
    return BloomFilter{uint(c), uint(m), uint(k), NewBitSet(m)}, nil
}


func main() {

    bf, _ := New(0, 32, 3)

    fmt.Println(bf)

    v, t, t2 := "some value", "some test value", "some value"

    fmt.Println("adding", v)
    bf.Add(v)

    fmt.Println(bf.b)

    if (bf.Test(t)) {
        fmt.Println(fmt.Sprintf("%s is (probably :)) in bf", t))
    } else {
        fmt.Println(fmt.Sprintf("%s is not in bf", t))
    }

    if (bf.Test(t2)) {
        fmt.Println(fmt.Sprintf("%s is (probably :)) in bf", t2))
    } else {
        fmt.Println(fmt.Sprintf("%s is not in bf", t2))
    }
}
