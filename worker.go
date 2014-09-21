package main

import (
	"bytes"
	"encoding/json"
	"github.com/iron-io/iron_go/mq"
	"github.com/joho/godotenv"
	carve "github.com/motdotla/carve"
	"github.com/nu7hatch/gouuid"
	"github.com/rlmcpherson/s3gof3r"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Payload struct {
	Documents []Document `json:"documents"`
	Meta      *struct {
		Postscript string `json:"postscript,omitempty"`
	} `json:"meta,omitempty"`
}
type Document struct {
	Pages   []Page `json:"pages"`
	Status  string `json:"status"`
	Url     string `json:"url"`
	Webhook string `json:"webhook"`
}

type Page struct {
	Number int    `json:"number"`
	Url    string `json:"url"`
}

func main() {
	godotenv.Load()

	loop_milliseconds, _ := time.ParseDuration(os.Getenv("LOOP_MILLISECONDS") + "ms")

	for x := range time.Tick(loop_milliseconds) {
		log.Println(x)

		queue := mq.New(os.Getenv("QUEUE"))
		msg, err := queue.Get()
		if err != nil {
			log.Println(err)
		}
		if msg != nil {
			s := msg.Body
			var payload Payload
			err := json.Unmarshal([]byte(s), &payload)
			if err != nil {
				log.Println(err)
			}
			msg.Delete()

			go Process(payload)
		}
	}
}

func Process(payload Payload) {
	s, err := carve.Convert(payload.Documents[0].Url, os.Getenv("CARVE_PNGS_OUTPUT_DIR"))
	if err != nil {
		log.Println(err)
	}

	pngs := strings.Split(s, ",")
	png_urls, err := Upload(pngs)
	if err != nil {
		log.Println(err)
	}
	log.Println(png_urls)

	Webhook(png_urls, payload)
}

func Webhook(pages []Page, payload Payload) {
	payload.Documents[0].Pages = pages
	payload.Documents[0].Status = "processed"
	marshaled_payload, _ := json.Marshal(payload)
	payload_string := string(marshaled_payload)

	req, err := http.NewRequest("POST", payload.Documents[0].Webhook, bytes.NewBufferString(payload_string))
	if err != nil {
		log.Println(err)
	} else {
		req.Header.Set("Content-Type", "application/json")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println(err)
		} else {
			defer res.Body.Close()
		}
	}
}

func Upload(pngs []string) ([]Page, error) {
	keys, _ := s3gof3r.EnvKeys()
	s3 := s3gof3r.New("", keys)
	bucket := s3.Bucket(os.Getenv("S3_BUCKET"))
	u, _ := uuid.NewV4()
	folder := u.String()

	pages := []Page{}

	for i := range pngs {
		uri, err := put(pngs[i], folder, bucket)
		if err != nil {
			return pages, err
		}

		var page Page
		page.Number = i + 1
		page.Url = uri

		pages = append(pages, page)
	}

	return pages, nil
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
