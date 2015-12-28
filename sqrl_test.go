package main

import (
	"testing"
	// "fmt"
)

func TestDecode(t *testing.T) {
	testStr := "dmVyPTENCmlkaz1OZ3JRRkN0NVJfSWg0ajNiZzlMZlFUOS0ySVRPYnAzODhET00tR25NdHV3DQpjbWQ9bG9naW4NCg"
	m := Decode(testStr)

	if m["ver"] != "1" {
		t.Fatal("Decode failed")
	}

	if m["idk"] != "NgrQFCt5R_Ih4j3bg9LfQT9-2ITObp388DOM-GnMtuw" {
		t.Fatal("Decode failed")
	}

	if m["cmd"] != "login" {
		t.Fatal("Decode failed")
	}
}

func TestGetSQRLdata(t *testing.T) {
	_, err := GetSQRLdata("", "", "")
	if err.Error() != "Inputs where empty" {
		t.Fatal(err)
	}

	client := "dmVyPTENCmlkaz1OZ3JRRkN0NVJfSWg0ajNiZzlMZlFUOS0ySVRPYnAzODhET00tR25NdHV3DQpjbWQ9bG9naW4NCg"
	server := "MTkyLjE2OC4xLjE1NTo4MDgwL2F1dGgvc3FybD9udXQ9MTIz"
	sign := "ccABaAJqsaA8nABNksEXGUOiC5c9UzToizquXhFSLrXBw7a1pVo6Zzh5GZRKKFFK9b2g0VIECdazjNskE7xtBw"
	auth, err := GetSQRLdata(client, server, sign)
	if err != nil {
		t.Fatal(err)
	}

	if auth.Nonce != "123" {
		t.Fatal("Nonce did not match")
	}

	if !auth.Verify() {
		t.Fatal("Did not verify")
	}
}
