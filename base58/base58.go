package base58

import (
	"fmt"
	"math/big"
)

var (
	zero       = big.NewInt(0)
	base       = big.NewInt(58)
	alphabetEn = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")
	alphabetDe = makeAlphabetDecode(alphabetEn)
)

func Encode(raw []byte) (encoded string) {
	// Count of leading 0x00
	zc := 0
	for _, b := range raw {
		if b == 0x00 {
			zc++
		} else {
			break
		}
	}

	n, m := new(big.Int), new(big.Int)
	n.SetBytes(raw)
	cnt := n.BitLen()/5 + 1 + zc
	bs := make([]byte, cnt, cnt)

	i := cnt - 1
	for ; n.Cmp(zero) != 0; i-- {
		n.DivMod(n, base, m)
		bs[i] = alphabetEn[m.Int64()]
	}
	for j := 0; j < zc; j++ {
		bs[i] = alphabetEn[0]
		i--
	}
	return string(bs[i+1:])
}

func Decode(encoded string) (raw []byte, err error) {
	n := new(big.Int)
	for _, b := range encoded {
		i := alphabetDe[b]
		if i < 0 {
			return nil, fmt.Errorf(`invalid base58 byte "%x"`, b)
		}
		n.Mul(n, base)
		n.Add(n, big.NewInt(int64(i)))
	}

	// Count of leading 0x00
	zc := 0
	for _, b := range encoded {
		if byte(b) == alphabetEn[0] {
			zc++
		} else {
			break
		}
	}

	raw = make([]byte, zc, n.BitLen()/8+zc)
	raw = append(raw, n.Bytes()...)
	return raw, nil
}

func makeAlphabetDecode(alphabet []byte) [256]int {
	var alphabetDecode [256]int
	for i := 0; i < 256; i++ {
		alphabetDecode[i] = -1
	}
	for i, b := range alphabet {
		alphabetDecode[b] = i
	}
	return alphabetDecode
}
