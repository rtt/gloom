package main

/*
    todo: tidy up sigs, ie (a string, b string) => (a, b string)
*/

import (
    "fmt"
    "crypto/sha256"
    "io"
    "math"
)

type BitSet struct {
    sz uint
    bits []uint
}

type BloomFilter struct {
    c uint // count of items in filter
    m uint // size of bitset (in bits)
    f []BloomFilterHash
    b BitSet
}

type BloomFilterHash struct {
    algo string
    salt string
}

/*
 * Returns the amount of filters in a BloomFilter
 */
func (bf BloomFilter) FilterLen() (sz int) {
    return len(bf.f)
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
func (b *BitSet) SetVal(val uint64) error {
    x := uint64(1)
    for i := 0; uint(i) < b.sz; i++ {
        if (uint64(val) & x) == x {
            b.bits[i] = 1
        }
        // shiiiiiift and bail when ints wrap
        x = uint64(x << 1)
        if (x < 0) {
            return nil
        }
    }

    return nil; // todo!
}

/*
 * Tests whether a uint64 lies in the bitset
 */

func (b *BitSet) TestVal(val uint64) (bool, error) {
    for i := uint(0); i < b.sz; i++ {

    }
    return true, nil
}

/*
 * Adds a value (item) to a BloomFilter
 */
func (bf *BloomFilter) Add(item string) {

    // get a hash for each item
    hashes := bf.Digest(item)

    // add each hash into the bf's bitset
    for i := range(hashes) {
        // eg, 48 bits = 2 chars per byte. so 48/16 = 12 = 12 chars = 48 bits
        bf.b.SetVal(hex_uint64(hashes[i][0:((bf.m / 8) * 2)]))
    }

    // increment the count
    bf.c += 1
}

/*
 * Tests whether an item is in a bloom filter
 */
func (bf BloomFilter) Test(item string) bool {
    hashes := bf.Digest(item)
    results := make([]bool, len(bf.f))
    for i := range(hashes) {
        // eg, 48 bits = 2 chars per byte. so 48/16 = 12 = 12 chars = 48 bits
        r, _ := bf.b.TestVal(hex_uint64(hashes[i][0:((bf.m / 8) * 2)]))
        results[i] = r
    }
    return all(results)
}

/* Nicely tells us about a BloomFilter
 */
func (bf BloomFilter) String() string {
    return fmt.Sprintf("Bloom filter with %d filters, %d items and bitset size of %d", len(bf.f), bf.c, bf.m)
}

/*
 * Appends a filter to a bloom filter set
 */
func (bf *BloomFilter) AddFilter(b *BloomFilterHash) error {
    if (bf.c > 0) {
        return fmt.Errorf("Cannot append to filters while BloomFilter size is not zero")
    }

    bf.f = append(bf.f, *b)
    return nil
}

/*
 * Returns a digested value for all of the filters
 * in a BloomFilter
 */
func (bf BloomFilter) Digest(value string) []string {
    s := make([]string, 0, len(bf.f))
    for i := range(bf.f) {
        s = append(s, bf.f[i].Digest(value))
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

/*
 * Returns a new BloomFilterHash
 */
func NewFilter(f string, salt string) (BloomFilterHash, error) {
    if (len(salt) == 0) {
        return BloomFilterHash{}, fmt.Errorf("Salt must not be zero length")
    }
    return BloomFilterHash{f, salt}, nil
}

/*
 * Returns a set of BloomFilterHashes
 * according to a set of string arguments given
 * todo: make sure to unique these values
 */
func NewFilterMulti(f string, args ...string) []BloomFilterHash {
    sl := []BloomFilterHash{}
    for i := range(args) {
        f, _ := NewFilter(f, args[i])
        sl = append(sl, f)
    }
    return sl
}

/* Convenience function for constructing a new
 * bloom filter
 */
func New(c int, m int, fh []BloomFilterHash) BloomFilter {
    return BloomFilter{uint(c), uint(m), fh, NewBitSet(m)} // 48 bit bitset
}

/* Returns the salt
 */
func (b BloomFilterHash) Salt() string  {
    return b.salt
}

/* Returns a Digested value according to
 * the filter's hash and salt
 */
func (b BloomFilterHash) Digest(value string) string {
    // todo: check this is
    hash, _ := M[b.algo].(func(string, string) (string, error))(value, b.salt)
    return hash
}

/* Returns the sha256 digest of a string + salt
 * 32 bytes / 256bit
 */
func sha256_digest (value string, salt string) (string, error) {
    h := sha256.New()
    salted := fmt.Sprintf("%s%s", value, salt)
    io.WriteString(h, salted)
    return fmt.Sprintf("%x", h.Sum(nil)), nil
}

var M map[string]interface{} = make(map[string]interface{})

/* Returns a hex string as a uint64
 * oops. apparently i missed strconv.ParseInt first time round
 * todo: replace this as it is redundant
 */
func hex_uint64 (s string) uint64 {
    x := uint64(0)
    sl := len(s)
    m := make(map[string]int)

    m["0"]  = 0
    m["1"]  = 1
    m["2"]  = 2
    m["3"]  = 3
    m["4"]  = 4
    m["5"]  = 5
    m["6"]  = 6
    m["7"]  = 7
    m["8"]  = 8
    m["9"]  = 9
    m["a"]  = 10
    m["A"]  = 10
    m["b"]  = 11
    m["B"]  = 11
    m["c"]  = 12
    m["C"]  = 12
    m["d"]  = 13
    m["D"]  = 13
    m["e"]  = 14
    m["E"]  = 14
    m["f"]  = 15
    m["F"]  = 15

    for i := 0; i < sl; i++ {
        x += uint64(m[string(s[i])]) * uint64(math.Pow(16, float64(sl - i - 1)))
    }

    return x
}

func main() {
    // make a new set of filters
    fh := NewFilterMulti("sha256_digest", "asd", "asdasda", "231221")

    M["sha256_digest"] = sha256_digest

    v := "password"

    // spit a load of hashes
    for i := range(fh) {
        fmt.Println(v, ":", fh[i].Digest(v), "| salt:", fh[i].Salt())
    }

    // make a new bloom filter
    // 0 provisional items, bitsize 8, with however many filter hashes
    bf := New(0, 48, fh)

    // lets see what this is all about
    fmt.Println(bf)

    // ask a bloom filter to digest a value against all of its filters
    hashes := bf.Digest(v)
    for i := range(hashes) {
        fmt.Println("hash", hashes[i])
    }

    // lets add foo
    vv := "foo"
    bf.Add(vv)

    // and lets print that mofo
    fmt.Println(bf)

    fmt.Println(bf.b)

    bf.Add("foo2")

    fmt.Println(bf.b)

    if (bf.Test("sd")) {
        fmt.Println("sd is (probably :)) in bf")
    } else {
        fmt.Println("sd is not in bf")
    }

}
