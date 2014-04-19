# carve-worker

Background worker for converting PDFs into an array of PNGs.

I've tried to make it as easy to use as possible, but if you have any feedback please [let me know](mailto:scott@scottmotte.com).

## Installation

```
git clone https://github.com/scottmotte/carve-worker.git
cd carve-worker
go get
cp .env.example .env
go run worker.go
```

Make sure you edit the contents of `.env`.


