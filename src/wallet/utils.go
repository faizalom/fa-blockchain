package wallet

import (
	"log"

	"github.com/mr-tron/base58"
)

func Base58Encode(input []byte) []byte {
	encode := base58.Encode(input)

	return []byte(encode)
}

func Base58Decode(input string) ([]byte, error) {
	decode, err := base58.Decode(input[:])
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return decode, nil
}

// 0 O l I + /
