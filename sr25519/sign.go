package main

import (
	"fmt"
	"os"
	"encoding/hex"
	schnorrkel "github.com/ChainSafe/go-schnorrkel"
)

func main() {
	message := os.Args[1]
	msg := []byte(message)
	signingCtx :=  []byte("substrate")
	signingTranscript := schnorrkel.NewSigningContext(signingCtx, msg)

	priv, pub, _ := schnorrkel.GenerateKeypair()

	sig, _ := priv.Sign(signingTranscript)

	pubEncoded := pub.Encode()
	sigEncoded := sig.Encode()

	fmt.Printf("\n%q %q %q\n", message, "0x"+hex.EncodeToString(pubEncoded[:]), "0x"+hex.EncodeToString(sigEncoded[:]))
}
