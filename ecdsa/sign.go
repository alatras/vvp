package main

import (
	"fmt"
	"os"
	"encoding/hex"
	"github.com/miguelmota/go-ethereum-hdwallet"
)

const DerivationPath string = "m/44'/60'/0'/0/0"

func main() {
	message := os.Args[1]
	msg := []byte(message)
	mnemonic := "tag volcano eight thank tide danger coast health above argue embrace heavy"
	wallet, _ := hdwallet.NewFromMnemonic(mnemonic)
	path := hdwallet.MustParseDerivationPath(DerivationPath)
	account, _ := wallet.Derive(path, true)
	signature, _ := wallet.SignText(account, msg)
	address, _ := wallet.AddressHex(account)
	fmt.Printf("\n%q %q %q\n", message, address, "0x"+hex.EncodeToString(signature))
}
