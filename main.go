package main

import (
	"crypto/ecdsa"
	"flag"
	"fmt"
	"os"
	"strings"

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

func findVanityAddress(firstChars string, caseSensitive bool, result chan pair) {
	for {
		pair := generatePair()
		toCheck := pair.address[2 : len(firstChars)+2]
		if (!caseSensitive && strings.EqualFold(toCheck, firstChars)) || toCheck == firstChars {
			result <- *pair
		}
	}
}

func main() {

	start := flag.String("start", "", "starting chars in address, or * for any address")
	caseSensitive := flag.Bool("caseSensitive", true, "default true. when true, ignore case of result")
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

	concurrency := 15
	for i := 0; i < concurrency; i++ {
		go findVanityAddress(*start, *caseSensitive, foundPair)
	}
	res := <-foundPair

	fmt.Printf("Private Key:           %s\n", res.privateKey)
	fmt.Printf("Address:               %s\n", res.address)

}
