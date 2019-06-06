package wallet

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"haha/base58"
)

const lenAddrChecksum = 4

type Addr struct {
	version              byte
	pubKeyHash, checksum []byte
	content              string
}

func NewAddr(version byte, pubKeyHash []byte) *Addr {
	cont := make([]byte, 0, 1+len(pubKeyHash)+lenAddrChecksum)
	cont = append(cont, version)
	cont = append(cont, pubKeyHash...)
	checksum := checksum(cont)
	cont = append(cont, checksum...)
	return &Addr{version: version, pubKeyHash: pubKeyHash, checksum: checksum, content: base58.Encode(cont)}
}

func DecodeAddr(addr string) (*Addr, error) {
	cont, err := base58.Decode(addr)
	if err != nil {
		return nil, err
	}
	lenCont := len(cont)
	actualChecksum := cont[lenCont-lenAddrChecksum:]
	targetChecksum := checksum(cont[:lenCont-lenAddrChecksum])
	if bytes.Compare(actualChecksum, targetChecksum) != 0 {
		return nil, fmt.Errorf("invalid address")
	}
	addrObj := &Addr{version: cont[0], pubKeyHash: cont[1 : lenCont-lenAddrChecksum], checksum: actualChecksum, content: addr}
	return addrObj, nil
}

func (addr *Addr) String() string {
	return addr.content
}

func (addr *Addr) PubKeyHash() []byte {
	return addr.pubKeyHash
}

func checksum(payload []byte) []byte {
	sum := sha256.Sum256(payload)
	sum = sha256.Sum256(sum[:])
	return sum[:lenAddrChecksum]
}
