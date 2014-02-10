package main

import (
    "fmt"
    "crypto/sha256"
    "io"
)

type BitSet struct {
    sz uint
    bits []uint
}

type BloomFilter struct {
    c uint // count of items in filter
    m uint // size of bitset
    f []BloomFilterHash
    b BitSet
}

type BloomFilterHash struct {
    algo string
    salt string
}

func (bf BloomFilter) FilterLen() (sz int) {
    return len(bf.f)
}

/***************/
func InitBitSet(l int) BitSet {
    b := BitSet{}
    b.bits = make([]uint, 0, l)
    return b
}

/***************/

/* Adds a value (item) to a BloomFilter
 */
func (bf *BloomFilter) Add(item string) {

    // get a hash for each item
    hashes := bf.Digest(item)
    for i := range(hashes) {
        fmt.Println("adding...", item, hashes[i])
        // todo
    }

    // increment the count
    bf.c += 1
}

func (bf BloomFilter) Test (item string) bool {
    // todo
    return true
}

/* Nicely tells us about a BloomFilter
 */
func (bf BloomFilter) String() string {
    return fmt.Sprintf("Bloom filter with %d filters, %d items and bitset size of %d", len(bf.f), bf.c, bf.m)
}

/* Appends a filter to a bloom filter set
 */
func (bf *BloomFilter) AddFilter(b *BloomFilterHash) error {
    if (bf.c > 0) {
        return fmt.Errorf("Cannot append to filters while BloomFilter size is not zero")
    }

    bf.f = append(bf.f, *b)
    return nil
}

func (bf BloomFilter) Digest(value string) []string {
    s := make([]string, 0, len(bf.f))
    for i := range(bf.f) {
        s = append(s, bf.f[i].Digest(value))
    }
    return s
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
    return BloomFilter{uint(c), uint(m), fh, InitBitSet(48)} // 48 bit bitset
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

func sha256_digest (value string, salt string) (string, error) {
    h := sha256.New()
    salted := fmt.Sprintf("%s%s", value, salt)
    io.WriteString(h, salted)
    return fmt.Sprintf("%x", h.Sum(nil)), nil
}

var M map[string]interface{} = make(map[string]interface{})

func main() {
    // make a new set of filters
    fh := NewFilterMulti("sha256_digests", "asd", "asd2", "11313", "1313", "207yd", "13213", "1223123", "131313")

    M["sha256_digest"] = sha256_digest

    v := "password"

    // spit a load of hashes
    for i := range(fh) {
        fmt.Println(v, ":", fh[i].Digest(v), "| salt:", fh[i].Salt())
    }

    // make a new bloom filter
    // 0 provisional items, bitsize 8, with however many filter hashes
    bf := New(0, 8, fh)

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

    if (bf.Test(vv)) {
        fmt.Println(v, "is in bf")
    } else {
        fmt.Println(v, "is not in bf")
    }

}
