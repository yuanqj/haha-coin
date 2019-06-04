package wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"io/ioutil"
	"os"
)

const walletFile = "wallet.dat"

type Wallets struct {
	wallets map[string]*Wallet
}

func NewWallets() (*Wallets, error) {
	ws := &Wallets{make(map[string]*Wallet)}
	if err := ws.Load(); err != nil {
		return nil, err
	}
	return ws, nil
}

func (ws *Wallets) CreateWallet() (string, error) {
	w, err := NewWallet()
	if err != nil {
		return "", err
	}
	addr, err := w.GetAddr()
	if err != nil {
		return "", err
	}
	ws.wallets[addr] = w
	return addr, nil
}

func (ws *Wallets) GetAddrs() []string {
	addrs := make([]string, 0, len(ws.wallets))
	for addr := range ws.wallets {
		addrs = append(addrs, addr)
	}
	return addrs
}

func (ws *Wallets) GetWallet(addr string) *Wallet {
	return ws.wallets[addr]
}

func (ws *Wallets) Save() error {
	cont := new(bytes.Buffer)
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(cont)
	if err := encoder.Encode(ws); err != nil {
		return err
	}
	return ioutil.WriteFile(walletFile, cont.Bytes(), 0644)
}

func (ws *Wallets) Load() error {
	if _, err := os.Stat(walletFile); err != nil {
		return err
	}
	cont, err := ioutil.ReadFile(walletFile)
	if err != nil {
		return err
	}
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(cont))
	if err := decoder.Decode(ws); err != nil {
		return err
	}
	return nil
}
