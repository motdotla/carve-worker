# carve-worker

<img src="https://raw.githubusercontent.com/motdotla/carve-worker/master/carve-worker.gif" alt="carve-worker" align="right" width="280" />

Background worker for converting PDFs into an array of PNGs. Works in tandem with [carve-worker](https://github.com/motdotla/carve-worker).

I've tried to make it as easy to use as possible, but if you have any feedback please [let me know](mailto:mot@mot.la).

## Installation

[![Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy)

```
heroku ps:scale web=0 worker=1
```

## Development
```
git clone https://github.com/motdotla/carve-worker.git
cd carve-worker
go get
cp .env.example .env
go run worker.go
```
