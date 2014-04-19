package main

import (
	"bytes"
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
	Pngs    []string `json:"pngs"`
	Status  string   `json:"status"`
	Url     string   `json:"url"`
	Webhook string   `json:"webhook"`
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

			go Process(document)
		}
	}
}

func Process(document Document) {
	s, err := carve.Convert(document.Url, os.Getenv("CARVE_PNGS_OUTPUT_DIR"))
	if err != nil {
		log.Println(err)
	}

	pngs := strings.Split(s, ",")
	png_urls, err := Upload(pngs)
	if err != nil {
		log.Println(err)
	}

	Webhook(png_urls, document)
}

func Webhook(png_urls []string, document Document) {
	log.Println(png_urls)

	document.Pngs = png_urls
	document.Status = "processed"
	payload := map[string]interface{}{"success": true, "document": document}

	marshaled_payload, _ := json.Marshal(payload)
	payload_string := string(marshaled_payload)

	// Prepare the request and headers
	log.Println(document.Webhook)
	log.Println(document)
	log.Println(payload_string)

	req, err := http.NewRequest("POST", document.Webhook, bytes.NewBufferString(payload_string))
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()
}

func Upload(pngs []string) ([]string, error) {
	keys, _ := s3gof3r.EnvKeys()
	s3 := s3gof3r.New("", keys)
	bucket := s3.Bucket(os.Getenv("S3_BUCKET"))
	u, _ := uuid.NewV4()
	folder := u.String()

	png_urls := []string{}

	for i := range pngs {
		uri, err := put(pngs[i], folder, bucket)
		if err != nil {
			return png_urls, err
		}
		png_urls = append(png_urls, uri)
	}

	return png_urls, nil
}

func put(path string, folder string, bucket *s3gof3r.Bucket) (string, error) {
	r, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer r.Close()

	header := make(http.Header)
	header.Add("x-amz-acl", "public-read")
	base := filepath.Base(path)
	fullpath := folder + "/" + base

	w, err := bucket.PutWriter(fullpath, header, nil)
	if err != nil {
		return "", err
	}
	if _, err = io.Copy(w, r); err != nil {
		return "", err
	}
	if err = w.Close(); err != nil {
		return "", err
	}

	uri := "https://" + os.Getenv("S3_BUCKET") + ".s3.amazonaws.com/" + fullpath

	return uri, nil
}
