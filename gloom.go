package main

/*
    todo: tidy up sigs, ie (a string, b string) => (a, b string)
*/

import (
    "bytes"
    "encoding/binary"
    "fmt"
    "strconv"
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
 * This is a straight implementation from: http://en.wikipedia.org/wiki/MurmurHash
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
        // shiiiiiift and bail when ints wrap
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
        // if the bit is 1 and the set at the same position is 1, this is a hit
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
    return fmt.Sprintf("Bloom filter with %d filters, %d items and bitset size of %d", bf.k, bf.c, bf.m)
}

/*
 * Returns a digested value for all of the filters
 * in a BloomFilter
 */
func (bf BloomFilter) Digest(value string) []uint32 {
    s, i := make([]uint32, 0, bf.k), 0
    for ; i < int(bf.k); i++ {
        vl := len(value)
        h := (Murmur32([]byte(value[0:vl / 2])) + uint32(i) + Murmur32([]byte(value[vl / 2:])))
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
func New(c int, m int, k int) BloomFilter {
    if m > 32 {
        m =  32
    }
    return BloomFilter{uint(c), uint(m), uint(k), NewBitSet(m)}
}

//var M map[string]interface{} = make(map[string]interface{})

/* Returns a hex string as a uint64
 * oops. apparently i missed strconv.ParseInt first time round
 * todo: replace this as it is redundant
 */
func hex_uint64 (s string) uint64 {
    i, _ := strconv.ParseInt(s, 16, 64)
    return uint64(i)
}

func main() {

    bf := New(0, 32, 3)

    fmt.Println(bf)

    bf.Add("test")
    fmt.Println(bf.b)
    t := "1231"

    if (bf.Test(t)) {
        fmt.Println(fmt.Sprintf("%s is (probably :)) in bf", t))
    } else {
        fmt.Println(fmt.Sprintf("%s is not in bf", t))
    }

}
