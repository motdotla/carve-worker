package main

import (
	"encoding/json"
	"github.com/iron-io/iron_go/mq"
	"github.com/joho/godotenv"
	"github.com/nu7hatch/gouuid"
	"github.com/rlmcpherson/s3gof3r"
	carve "github.com/scottmotte/carve"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
	s, err := carve.Convert(url, os.Getenv("CARVE_PNGS_OUTPUT_DIR"))
	if err != nil {
		log.Println(err)
	}

	keys, _ := s3gof3r.EnvKeys()
	s3 := s3gof3r.New("", keys)
	bucket := s3.Bucket("carvedevelopment")
	u, _ := uuid.NewV4()
	folder := u.String()

	pngs := strings.Split(s, ",")
	for i := range pngs {
		Upload(pngs[i], folder, bucket)
		//go Upload(pngs[i], folder, bucket) // this should work, but make sure it is ok for large amounts of stuff
	}
}

func Upload(path string, folder string, bucket *s3gof3r.Bucket) {
	r, err := os.Open(path)
	if err != nil {
		log.Println(err)
	}
	defer r.Close()

	header := make(http.Header)
	base := filepath.Base(path)

	w, err := bucket.PutWriter(folder+"/"+base, header, nil)
	if err != nil {
		log.Println(err)
	}
	if _, err = io.Copy(w, r); err != nil {
		log.Println(err)
	}
	if err = w.Close(); err != nil {
		log.Println(err)
	}
	log.Println(path)
}
