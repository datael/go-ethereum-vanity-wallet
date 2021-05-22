package main

import (
	"crypto/ecdsa"
	"flag"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type pair struct {
	privateKey string
	address    string
}

func generatePair() *pair {

	privateKey, _ := crypto.GenerateKey()

	privateKeyBytes := crypto.FromECDSA(privateKey)
	private := hexutil.Encode(privateKeyBytes)[2:]

	publicKey := privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	return &pair{privateKey: private, address: address}
}

func findVanityAddress(firstChars string, result chan pair) {
	for {
		pair := generatePair()
		if pair.address[2:len(firstChars)+2] == firstChars {
			result <- *pair
		}
	}
}

func main() {

	start := flag.String("start", "", "starting chars in address, or * for any address")
	flag.Parse()

	if len(*start) == 0 && *start != "*" {
		flag.Usage()
		os.Exit(1)
	}

	if *start == "*" {
		*start = ""
	}

	fmt.Printf("Looking for ETH wallet with address starting '0x%v'\n\n", *start)

	foundPair := make(chan pair)

	concurrency := 32
	for i := 0; i < concurrency; i++ {
		go findVanityAddress(*start, foundPair)
	}
	res := <-foundPair

	fmt.Printf("Private Key:           %s\n", res.privateKey)
	fmt.Printf("Address:               %s\n", res.address)

}
