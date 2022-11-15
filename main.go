package main

import (
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
	"pault.ag/go/sshsig"
)

func main() {
	var algo sshsig.HashAlgo

	pubKey := flag.String("pub", "", "Public key")
	file := flag.String("f", "", "")
	nameSpace := flag.String("n", "file", "Namespace")
	signature := flag.String("sig", fmt.Sprintf("%s.sig", *file), "Signature of file")
	flag.Parse()

	pubKeyData, err := os.ReadFile(*pubKey)
	if err != nil {
		log.Fatal(err)
	}

	fileData, err := os.ReadFile(*file)
	if err != nil {
		log.Fatal(err)
	}

	sigData, err := os.ReadFile(*signature)
	if err != nil {
		log.Fatal(err)
	}

	pk, _, _, _, err := ssh.ParseAuthorizedKey(pubKeyData)
	if err != nil {
		log.Fatal(err)
	}

	block, _ := pem.Decode(sigData)
	sig, err := sshsig.ParseSignature(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	switch sig.HashAlgorithm {
	case sshsig.HashAlgoSHA512:
		algo = sshsig.HashAlgoSHA512
	default:
		log.Fatalln(fmt.Errorf("%q not implimented", sig.HashAlgorithm))
	}

	ch, _ := algo.Hash()
	h := ch.New()
	h.Write(fileData)
	hash := h.Sum(nil)

	err = sshsig.Verify(pk, []byte(*nameSpace), sig.HashAlgorithm, hash, sig)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Signature OK")
}
