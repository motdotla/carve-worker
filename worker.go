package main

import (
	"encoding/json"
	"github.com/iron-io/iron_go/mq"
	"github.com/joho/godotenv"
	"github.com/rlmcpherson/s3gof3r"
	carve "github.com/scottmotte/carve"
	"io"
	"log"
	"net/http"
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

	r, err := os.Open("./README.md")
	if err != nil {
		log.Println(err)
	}
	defer r.Close()

	keys, _ := s3gof3r.EnvKeys()
	s3 := s3gof3r.New("", keys)
	bucket := s3.Bucket("carvedevelopment")

	log.Println(bucket)
	header := make(http.Header)
	w, err := bucket.PutWriter("README.md", header, nil)
	if err != nil {

	}
	if _, err = io.Copy(w, r); err != nil {
		log.Println(err)
	}
	if err = w.Close(); err != nil {
		log.Println(err)
	}

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
