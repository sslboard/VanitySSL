package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func verifySignature(customer, domain, sig, secret string) bool {
	if secret == "" || sig == "" {
		return false
	}
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(customer + "\n" + domain))
	expected := hex.EncodeToString(h.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(sig))
}

func handler(secret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		customer := r.Header.Get("X-Customer-ID")
		domain := r.Header.Get("X-Customer-Domain")
		sig := r.Header.Get("X-Vanity-Signature")
		valid := verifySignature(customer, domain, sig, secret)

		fmt.Fprintf(w, "<html><body>")
		fmt.Fprintf(w, "<h1>Backend Response</h1>")
		fmt.Fprintf(w, "<p>Customer ID: %s</p>", customer)
		fmt.Fprintf(w, "<p>Domain: %s</p>", domain)
		fmt.Fprintf(w, "<p>Signature Valid: %t</p>", valid)
		fmt.Fprintf(w, "<p>Method: %s</p>", r.Method)
		fmt.Fprintf(w, "<p>Path: %s</p>", r.URL.Path)
		fmt.Fprintf(w, "<p>Body Length: %d</p>", len(body))
		fmt.Fprintf(w, "<h2>Headers</h2><pre>")
		for name, vals := range r.Header {
			for _, v := range vals {
				fmt.Fprintf(w, "%s: %s\n", name, v)
			}
		}
		fmt.Fprintf(w, "</pre></body></html>")
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	secret := os.Getenv("PROXY_SECRET")
	http.HandleFunc("/", handler(secret))
	log.Printf("Dummy backend running on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
