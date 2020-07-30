package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/cloudfoundry-community/go-cfenv"
	_ "github.com/go-sql-driver/mysql"
)

var serviceName string

func init() {
	serviceName = os.Getenv("SERVICE_NAME")
}

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	appEnv, err := cfenv.Current()
	if err != nil {
		log.Fatal(err)
	}
	mariadb, err := appEnv.Services.WithName(serviceName)
	if err != nil {
		log.Fatal(err)
	}

	uri, err := url.Parse(mariadb.Credentials["uri"].(string))
	if err != nil {
		log.Fatal(err)
	}
	uriStr := fmt.Sprintf("%s@tcp(%s)/%s", uri.User.String(), uri.Hostname(), strings.TrimPrefix(uri.Path, "/"))
	fmt.Printf("Connecting to %q\n", uriStr)

	db, err := sql.Open("mysql", uriStr)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "8080"
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ok")
	})
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
