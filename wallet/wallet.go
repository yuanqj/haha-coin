package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
)

const version = byte(0x00)

type Wallet struct {
	PrvKey *ecdsa.PrivateKey
	PubKey []byte
	Addr   string
}

func NewWallet() (*Wallet, error) {
	prv, pub, err := genKeyPair()
	if err != nil {
		return nil, err
	}
	pubKeyHash, err := HashPubKey(pub)
	if err != nil {
		return nil, err
	}
	return &Wallet{PrvKey: prv, PubKey: pub, Addr: NewAddr(version, pubKeyHash).String()}, nil
}

func HashPubKey(key []byte) ([]byte, error) {
	keySHA256 := sha256.Sum256(key)
	hasher := ripemd160.New()
	if _, err := hasher.Write(keySHA256[:]); err != nil {
		return nil, err
	}
	return hasher.Sum(nil), nil
}

func genKeyPair() (prv *ecdsa.PrivateKey, pub []byte, err error) {
	curve := elliptic.P256()
	prv, err = ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return
	}

	pub = make([]byte, 64)
	x, y := prv.PublicKey.X.Bytes(), prv.PublicKey.Y.Bytes()
	copy(pub[32-len(x):], x)
	copy(pub[64-len(x):], y)
	return
}
