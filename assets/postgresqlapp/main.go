/*
   Copyright 2020 SUSE

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cloudfoundry-community/go-cfenv"
	pgx "github.com/jackc/pgx/v4"
)

func main() {
	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		log.Fatal("SERVICE_NAME not set")
	}

	appEnv, err := cfenv.Current()
	if err != nil {
		log.Fatal(err)
	}
	postgresqlService, err := appEnv.Services.WithName(serviceName)
	if err != nil {
		log.Fatal(err)
	}

	uriStr := postgresqlService.Credentials["uri"].(string)
	fmt.Printf("Connecting to %q\n", uriStr)

	ctx := context.Background()

	db, err := pgx.Connect(ctx, uriStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(ctx)

	if err := db.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(ctx, createTableStatement); err != nil {
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

const createTableStatement = `
CREATE Table mits(
	id SERIAL
);
`
