package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"haha/base58"
)

const version = byte(0x00)
const walletFile = "wallet.dat"
const lenAddrChecksum = 4

type Wallet struct {
	PrvKey *ecdsa.PrivateKey
	PubKey []byte
}

type Wallets struct {
	Wallets map[string]*Wallet
}

func NewWallet() (*Wallet, error) {
	prv, pub, err := genKeyPair()
	if err != nil {
		return nil, err
	}
	return &Wallet{prv, pub}, nil
}

func (w *Wallet) GetAddr() (addr string, err error) {
	pubKeyHash, err := HashPubKey(w.PubKey)
	if err != nil {
		return
	}

	cont := make([]byte, 0, 1+len(pubKeyHash)+lenAddrChecksum)
	cont = append(cont, version)
	cont = append(cont, pubKeyHash...)
	checksum := checksum(cont)
	cont = append(cont, checksum...)
	return base58.Encode(cont), nil
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

func checksum(payload []byte) []byte {
	sum := sha256.Sum256(payload)
	sum = sha256.Sum256(sum[:])
	return sum[:lenAddrChecksum]
}
