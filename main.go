package main

import (
	"crypto/ecdsa"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type pair struct {
	privateKey string
	address    string
}

var (
	start         = flag.String("start", "", "starting chars in address, or * for any address")
	caseSensitive = flag.Bool("caseSensitive", true, "when true, case of wallet string must match 'start'. defaults to true")
	concurrency   = flag.Int("concurrency", runtime.NumCPU(), "number of goroutines to use. defaults to number of cpus")
	numResults    = flag.Int("numResults", 1, "number of results to wait for. defaults to 1")
	singleLine    = flag.Bool("singleLine", false, "when true, address then private key are output on a single line separated by a space. defaults to false")
	verbosity     = flag.Int("verbosity", 0, "verbosity of output. 0=results only, 1=also output start message, 2=also output settings")
)

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

	flag.Parse()

	if len(*start) == 0 && *start != "*" {
		flag.Usage()
		os.Exit(1)
	}

	if *start == "*" {
		*start = ""
	}

	if *verbosity >= 1 {
		fmt.Printf("Looking for ETH wallet with address starting '0x%v'\n", *start)
		fmt.Printf("\n")
	}
	if *verbosity >= 2 {
		fmt.Printf("Settings:\n")
		fmt.Printf("  start:           %v\n", *start)
		fmt.Printf("  caseSensitive:   %v\n", *caseSensitive)
		fmt.Printf("  concurrency:     %v\n", *concurrency)
		fmt.Printf("  numResults:      %v\n", *numResults)
		fmt.Printf("  verbosity:       %v\n", *verbosity)
		fmt.Printf("  singleLine:      %v\n", *singleLine)
		fmt.Printf("\n")
	}

	foundPair := make(chan pair)

	for i := 0; i < *concurrency; i++ {
		go findVanityAddress(*start, *caseSensitive, foundPair)
	}

	if *numResults < 1 {
		*numResults = 1
	}

	for i := 0; i < *numResults; i++ {
		res := <-foundPair

		if *singleLine {
			fmt.Printf("%s %s\n", res.address, res.privateKey)
		} else {
			fmt.Printf("Address:               %s\n", res.address)
			fmt.Printf("Private Key:           %s\n", res.privateKey)

			if *numResults > 1 {
				fmt.Printf("\n")
			}
		}

	}

}
