## Chirpy

Chirpy is a backend webserver built to learn how to build webservers in Go.

## 🛠️ Installation

`go get github.com/TokiLoshi/chirpy`

## ⚡️ Quick Start

Install Goose and get a postres instance going to connect to. In your env file you'll want to have the following:

- "DB_URL" the name of your database to connect to
- "JWT_SECRET" a randomly generated secret for your auth to check agains
- "POLKA_KEY" another randomly generated secret for your auth to check againse.

To run the server `go run .`

## Exposed Endpoints

### Admin endpoints

- "/admin/metrics" - it takes a GET request to see how many times the page was visited, this is managed by metricsHandler
- "/admin/reset-metrics" - a POST request to reset the metrics managed by resetHandler
- "/admin/reset" - a POST request managed by adminReset

### Static file Serving

- "/app/" a midleware metric
- "/assets/" a file server

### User API endpoints

- "/api/users" a POST request handled by validateUser
- "/api/login" a POST request handled by handleLogin
- "/api/refres" a POST request handled by handleRefresh
- "/api/revoke" a POST request handled by handleRevoke
- "/api/users" a PUT request handled by handleAuthentication

### Chirp API endpoints

- "/api/chirps" a POST request handled by createChipr
- "/api/chirps" a GET Request handled by getAllChirps
- "/api/chiprs/{chirpId}" a GET request handled by getSingleChirp
- "/api/chirps/{chirpId}" a DELETE request handled by handleDelete

### Webhooks

- "api/polka/webhooks" an imaginary webhook to handle upgraded users, it is a POST request handled by handleUpgrade
