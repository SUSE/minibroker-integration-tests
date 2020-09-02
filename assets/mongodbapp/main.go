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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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
	mongodbService, err := appEnv.Services.WithName(serviceName)
	if err != nil {
		log.Fatal(err)
	}

	uriStr := mongodbService.Credentials["uri"].(string)
	fmt.Printf("Connecting to %q\n", uriStr)

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uriStr))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	databaseStr := mongodbService.Credentials["database"].(string)
	collection := client.Database(databaseStr).Collection("mits")
	expectedValue := Mits{"12345"}
	if _, err := collection.InsertOne(ctx, expectedValue); err != nil {
		log.Fatal(err)
	}

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	var results []Mits
	if err := cursor.All(ctx, &results); err != nil {
		log.Fatal(err)
	}

	if len(results) != 1 {
		log.Fatal(fmt.Errorf("Invalid result length: %d, expected 1", len(results)))
	}

	value := results[0]

	if value.MitsID != expectedValue.MitsID {
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

// Mits represents an instance to be inserted in the mits collection.
type Mits struct {
	MitsID string `json:"mits_id"`
}
