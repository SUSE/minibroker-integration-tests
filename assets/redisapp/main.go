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
	redis "github.com/go-redis/redis/v8"
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
	redisEnv, err := appEnv.Services.WithName(serviceName)
	if err != nil {
		log.Fatal(err)
	}

	uriStr := redisEnv.Credentials["uri"].(string)
	opt, err := redis.ParseURL(uriStr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Connecting to %q\n", uriStr)
	db := redis.NewClient(opt)
	defer db.Close()

	ctx := context.Background()

	if _, err := db.Ping(ctx).Result(); err != nil {
		log.Fatal(err)
	}

	const key = "foo"
	const expectedValue = "bar"

	if err := db.Set(ctx, key, expectedValue, 0).Err(); err != nil {
		log.Fatal(err)
	}

	value, err := db.Get(ctx, key).Result()
	if err != nil {
		log.Fatal(err)
	}

	if value != expectedValue {
		log.Fatal(fmt.Errorf("Value %q is not the expected %q", value, expectedValue))
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
