package main

import (
	"encoding/base64"
	"errors"
	"github.com/agl/ed25519"
	"log"
	"strings"
	// "fmt"
)

type SqrlData struct {
	PubKey    [32]byte
	Signature [64]byte
	Data      []byte // The signed data
	Nonce     string
}

func (sqrl *SqrlData) Verify() bool {
	return ed25519.Verify(&sqrl.PubKey, sqrl.Data, &sqrl.Signature)
}

func Decode(str string) map[string]string {
	sDec, err := base64.RawURLEncoding.DecodeString(str)
	if err != nil {
		log.Fatal(err)
	}

	m := make(map[string]string)
	pairs := strings.Split(string(sDec), "\n")
	for _, element := range pairs {
		s := strings.Split(element, "=")
		if len(s) > 1 {
			m[s[0]] = strings.Trim(s[1], "\r\n")
		}
	}

	return m
}

func GetSQRLdata(client, server, ids string) (SqrlData, error) {
	if (client == "") || (server == "") || (ids == "") {
		return SqrlData{}, errors.New("Inputs where empty")
	}

	cli := Decode(client)

	// get the nonce
	srv, err := base64.RawURLEncoding.DecodeString(server)
	if err != nil {
		log.Fatal(err)
	}
	nonce := string(srv)[strings.Index(string(srv), "?nut=")+5:]

	// decode pub key
	sDec1, err := base64.RawURLEncoding.DecodeString(cli["idk"])
	if err != nil {
		log.Fatal(err)
	}
	var pubkey [32]byte
	copy(pubkey[:], sDec1[0:32])

	sDec2, err := base64.RawURLEncoding.DecodeString(ids)
	var sign [64]byte
	copy(sign[:], sDec2[0:64])

	data := client + server

	return SqrlData{PubKey: pubkey, Signature: sign, Data: []byte(data), Nonce: nonce}, nil
}
