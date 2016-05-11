package main

import (
	"os"
	"net/http"
	"log"
	"fmt"
	"strconv"
)

func main() {
	port := "8080"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	uaaUrl := os.Getenv("UAA_URL")
	if uaaUrl == "" {
		fatal("You have to set the env var UAA_URL")
	}
	skipInsecure := false
	if os.Getenv("SKIP_INSECURE") != "" {
		tempInsecure, err := strconv.ParseBool(os.Getenv("SKIP_INSECURE"))
		if err == nil {
			skipInsecure = tempInsecure
		}
	}
	proxy := NewCustomProxy(uaaUrl, skipInsecure)

	http.HandleFunc("/", proxy.handle)
	log.Fatal(http.ListenAndServe(":" + port, nil))
}
func fatalIf(doing string, err error) {
	if err != nil {
		fatal(doing + ": " + err.Error())
	}
}
func fatal(message string) {
	fmt.Fprintln(os.Stdout, message)
	os.Exit(1)
}