package main

import (
	"encoding/json"
	"github.com/iron-io/iron_go/mq"
	"github.com/joho/godotenv"
	carve "github.com/scottmotte/carve"
	"log"
	"os"
	"time"
)

type Document struct {
	Pngs    []string
	Status  string
	Url     string
	Webhook string
}

func main() {
	godotenv.Load()

	for x := range time.Tick(500 * time.Millisecond) {
		log.Println(x)

		queue := mq.New(os.Getenv("QUEUE"))
		msg, err := queue.Get()
		if err != nil {
			log.Println(err)
		}
		if msg != nil {
			s := msg.Body
			var document Document
			err := json.Unmarshal([]byte(s), &document)
			if err != nil {
				log.Println(err)
			}
			msg.Delete()

			go Convert(document.Url)
		}
	}
}

func Convert(url string) {
	pngs, err := carve.Convert(url, os.Getenv("CARVE_PNGS_OUTPUT_DIR"))
	if err != nil {
		log.Println(err)
	}

	log.Println(pngs)
}
