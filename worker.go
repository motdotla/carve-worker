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

	pngs := strings.Split(s, ",")
	UploadAll(pngs)
}

func UploadAll(pngs []string) {
	keys, _ := s3gof3r.EnvKeys()
	s3 := s3gof3r.New("", keys)
	bucket := s3.Bucket(os.Getenv("S3_BUCKET"))
	u, _ := uuid.NewV4()
	folder := u.String()

	for i := range pngs {
		uri, err := Upload(pngs[i], folder, bucket)
		if err != nil {
			log.Println(err)
		}
		log.Println(uri)
	}
}

func Upload(path string, folder string, bucket *s3gof3r.Bucket) (string, error) {
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
