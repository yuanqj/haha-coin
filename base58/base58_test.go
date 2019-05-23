package base58_test

import (
	"encoding/hex"
	"haha/base58"
	"testing"
)

var pairs = []struct {
	hex string
	b58 string
}{
	{"61", "2g"},
	{"626262", "a3gV"},
	{"636363", "aPEr"},
	{"73696d706c792061206c6f6e6720737472696e67", "2cFupjhnEsSn59qHXstmK2ffpLv2"},
	{"00eb15231dfceb60925886b67d065299925915aeb172c06647", "1NS17iag9jJgTHD1VXjvLCEnZuQ3rJDE9L"},
	{"516b6fcd0f", "ABnLTmg"},
	{"bf4f89001e670274dd", "3SEo3LWLoPntC"},
	{"572e4794", "3EFU7m"},
	{"ecac89cad93923c02321", "EJDM8drfXA6uyA"},
	{"10c8511e", "Rt5zm"},
	{"00000000000000000000", "1111111111"},
}

func TestEncode(t *testing.T) {
	for _, pair := range pairs {
		raw, _ := hex.DecodeString(pair.hex)
		if res := base58.Encode(raw); res != pair.b58 {
			t.Errorf(`encode error: raw="%s", b58="%s", res="%s"`, pair.hex, pair.b58, res)
		}
	}
}

func TestDecode(t *testing.T) {
	for _, pair := range pairs {
		raw, _ := hex.DecodeString(pair.hex)
		if res := base58.Encode(raw); res != pair.b58 {
			t.Errorf(`decode error: \nb58="%s"\nraw="%s"\nres="%s"`, pair.b58, pair.hex, res)
		}
	}
}
