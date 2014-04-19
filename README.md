# carve-worker

Background worker for converting PDFs into an array of PNGs.

I've tried to make it as easy to use as possible, but if you have any feedback please [let me know](mailto:scott@scottmotte.com).

## Installation

```
git clone https://github.com/scottmotte/carve-worker.git
cd carve-worker
go get github.com/scottmotte/carve
go get github.com/joho/godotenv
go get github.com/iron-io/iron_go/mq
go get github.com/rlmcpherson/s3gof3r
go get github.com/nu7hatch/gouuid
cp .env.example .env
go run worker.go
```

Make sure you edit the contents of `.env`.


