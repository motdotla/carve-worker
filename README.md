# carve-worker

Background worker for converting PDFs into an array of PNGs.

I've tried to make it as easy to use as possible, but if you have any feedback please [let me know](mailto:scott@scottmotte.com).

## Installation

### Production

```
git clone https://github.com/scottmotte/carve-worker.git
cd carve-worker
heroku create -b https://github.com/ddollar/heroku-buildpack-multi.git
heroku config:set IRON_TOKEN=
heroku config:set IRON_PROJECT_ID=
heroku config:set QUEUE=
heroku config:set CARVE_PNGS_OUTPUT_DIR=
heroku config:set AWS_ACCESS_KEY_ID=
heroku config:set AWS_SECRET_ACCESS_KEY=
heroku config:set S3_BUCKET=
heroku config:set LOOP_MILLISECONDS=5000
git push heroku master
heroku ps:scale web=0 worker=1
heroku open
```

### Development
```
git clone https://github.com/scottmotte/carve-worker.git
cd carve-worker
go get
cp .env.example .env
go run worker.go
```

Make sure you edit the contents of `.env`.




