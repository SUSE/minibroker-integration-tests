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
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/streadway/amqp"
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
	amqpEnv, err := appEnv.Services.WithName(serviceName)
	if err != nil {
		log.Fatal(err)
	}

	uriStr := amqpEnv.Credentials["uri"].(string)
	fmt.Printf("Connecting to %q\n", uriStr)
	conn, err := amqp.Dial(uriStr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	queue, err := ch.QueueDeclare(
		"foo", // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatal(err)
	}

	expectedValue := "Hello World!"

	go func() {
		if err := ch.Publish(
			"",         // exchange
			queue.Name, // routing key
			false,      // mandatory
			false,      // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(expectedValue),
			},
		); err != nil {
			log.Fatal(err)
		}
	}()

	msgs, err := ch.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		log.Fatal(err)
	}

	msg := <-msgs
	value := string(msg.Body)

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
