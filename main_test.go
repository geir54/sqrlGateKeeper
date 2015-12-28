package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type mockDB struct {
}

func (db mockDB) get(pubkey []byte) error {
	if reflect.DeepEqual(pubkey, []byte{54, 10, 208, 20, 43, 121, 71, 242, 33, 226, 61, 219, 131, 210, 223, 65, 63, 126, 216, 132, 206, 110, 157, 252, 240, 51, 140, 248, 105, 204, 182, 236}) {
		return nil
	}

	return errors.New("some error")
}

func (db mockDB) set(pubkey []byte) error {
	return nil
}

func TestHandleAuth(t *testing.T) {
	db := mockDB{}
	nuts := initnutList(100)
	nuts.add("123")

	AuthRoute := handleAuth(db, nuts)

	client := "dmVyPTENCmlkaz1OZ3JRRkN0NVJfSWg0ajNiZzlMZlFUOS0ySVRPYnAzODhET00tR25NdHV3DQpjbWQ9bG9naW4NCg"
	server := "MTkyLjE2OC4xLjE1NTo4MDgwL2F1dGgvc3FybD9udXQ9MTIz"
	sign := "ccABaAJqsaA8nABNksEXGUOiC5c9UzToizquXhFSLrXBw7a1pVo6Zzh5GZRKKFFK9b2g0VIECdazjNskE7xtBw"
	req, _ := http.NewRequest("GET", "/?client="+client+"&server="+server+"&ids="+sign, nil)
	w := httptest.NewRecorder()

	AuthRoute(w, req)

	ret, _ := nuts.get("123")

	if ret.PubKey != [32]byte{54, 10, 208, 20, 43, 121, 71, 242, 33, 226, 61, 219, 131, 210, 223, 65, 63, 126, 216, 132, 206, 110, 157, 252, 240, 51, 140, 248, 105, 204, 182, 236} {
		t.Errorf("PubKey was not correct")
	}

	if ret.UnknownUser {
		t.Errorf("Should not be unknown user")
	}

	if w.Code != http.StatusOK {
		t.Errorf("Home page didn't return %v", http.StatusOK)
	}

}

func TestGetQR(t *testing.T) {
	nuts := initnutList(100)

	qrRoute := getQR(nuts)
	req, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()
	qrRoute(w, req)

	session, err := store.Get(req, "sessions")
	if err != nil {
		t.Errorf(err.Error())
	}

	nut, ok := session.Values["nut"].(string)

	if !ok {
		t.Errorf("Session not set")
	}

	_, ok = nuts.get(nut)

	if !ok {
		t.Errorf("Nut not set")
	}

	if w.Code != http.StatusOK {
		t.Errorf("Home page didn't return %v", http.StatusOK)
	}

}

func TestHome(t *testing.T) {
	nuts := initnutList(100)

	indexRoute := indexRoute(nuts)
	req, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()
	indexRoute(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Home page didn't return %v", http.StatusOK)
	}
}
