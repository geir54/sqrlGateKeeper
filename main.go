package main

import (
	"code.google.com/p/rsc/qr" // TODO: should be updated
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"text/template"
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))

func RandString(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal("error:", err)

	}

	return base64.StdEncoding.EncodeToString(b)
}

func getQR(nuts *nutList) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "sessions")
		if err != nil {
			log.Println(err)
			return
		}

		nut := RandString(10)

		session.Values["nut"] = nut
		session.Save(r, w)

		nuts.add(nut)

		q, err := qr.Encode("qrl://192.168.1.155:8080/sqrl/auth?nut="+nut, qr.H)
		if err != nil {
			log.Fatal(err)
		}

		i := q.PNG()
		w.Write(i)
	}
}

type databaseInter interface {
	get([]byte) error
	set([]byte) error
}

func handleAuth(db databaseInter, nuts *nutList) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		client := r.Form.Get("client")
		server := r.Form.Get("server")
		sign := r.Form.Get("ids")

		auth, err := GetSQRLdata(client, server, sign)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Someone tried to log in")

		if auth.Verify() {
			nut, ok := nuts.get(auth.Nonce)
			if ok { // Check if the nonce comes from us
				nut.PubKey = auth.PubKey

				err := db.get(auth.PubKey[:])
				if err != nil {
					nut.UnknownUser = true
					nuts.update(auth.Nonce, nut)
				} else {
					nut.Autenticated = true
					nuts.update(auth.Nonce, nut)
				}
			}
			fmt.Fprintf(w, "good")
		}
	}
}

func handleLogout() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "sessions")
		if err != nil {
			log.Println(err)
			return
		}

		delete(session.Values, "pubkey")
		delete(session.Values, "nut")
		session.Save(r, w)
		fmt.Fprintf(w, "You have been logged out")

	}
}

func handleAdd(db *database, nuts *nutList, adminPwd *string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		password := r.Form.Get("password")
		if password != *adminPwd {
			fmt.Fprintf(w, "Wrong password")
			return
		}
		session, err := store.Get(r, "sessions")
		if err != nil {
			log.Println(err)
			return
		}

		nut, _ := session.Values["nut"].(string)
		data, ok := nuts.get(nut)
		if ok {
			err := db.set(data.PubKey[:])
			if err != nil {
				log.Fatal(err)
			}

			nuts.delete(nut)
			fmt.Fprintf(w, "Key added")
		}
	}
}

func indexRoute(nuts *nutList) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		session, err := store.Get(r, "sessions")
		if err != nil {
			log.Println(err)
			return
		}

		pubkey, ok := session.Values["pubkey"].(string)

		if !ok || pubkey == "" {
			nut, ok := session.Values["nut"].(string)

			if ok && nut != "" {
				data, _ := nuts.get(nut)
				if data.Autenticated {
					log.Println("I know this person add session")
					session.Values["pubkey"] = "test"
					session.Save(r, w)
					goto proxy
					return
				}
				if data.UnknownUser {
					t, err := template.ParseFiles("public/admin.html")
					if err != nil {
						log.Fatal(err)
					}
					type Response struct {
						PubKey string
					}

					user := Response{PubKey: base64.StdEncoding.EncodeToString(data.PubKey[:])}
					err = t.Execute(w, user)
					if err != nil {
						log.Fatal(err)
					}
					return
				}
			}

			serveSingle(w, r, "public/index.html")
			return
		}

	proxy:
		url, err := url.Parse("http://127.0.0.1:80/")
		if err != nil {
			log.Fatal(err)
		}

		target := httputil.NewSingleHostReverseProxy(url)

		director := target.Director
		target.Director = func(req *http.Request) {
			director(req)

			req.URL.Opaque = req.RequestURI
			req.URL.RawQuery = ""
		}

		target.ServeHTTP(w, r)
	}
}

func serveSingle(w http.ResponseWriter, r *http.Request, filename string) {
	http.ServeFile(w, r, filename)
}

func main() {
	adminPwdPtr := flag.String("p", "", "The admin password")
	remoteServerPtr := flag.String("r", "", "The remote server")
	flag.Parse()

	if *adminPwdPtr == "" {
		fmt.Println("Please set an admin password")
		os.Exit(0)
	}

	if *remoteServerPtr == "" {
		fmt.Println("Please set a remote server")
		os.Exit(0)
	}

	nuts := initnutList(100)

	db := initDB("users.db")
	defer db.Close()

	r := mux.NewRouter()

	r.HandleFunc("/sqrl/qr", getQR(nuts))
	r.HandleFunc("/sqrl/auth", handleAuth(&db, nuts))
	r.HandleFunc("/sqrl/logout", handleLogout())
	r.HandleFunc("/sqrl/add", handleAdd(&db, nuts, adminPwdPtr))
	r.PathPrefix("/").HandlerFunc(indexRoute(nuts))

	server := &http.Server{
		Addr:    ":8080", //os.Getenv("IP")+":"+os.Getenv("PORT"),
		Handler: r,
	}

	fmt.Println("Starting server")
	server.ListenAndServe()
}
