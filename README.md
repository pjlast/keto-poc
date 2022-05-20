# Ory Keto proof of concept

This repo is just a small demonstration of how to create relations between objects in Ory Keto and how we can
use them to create links between different objects.

## Setup

Clone the Ory Keto repo:

```git clone git@github.com:ory/keto.git && cd keto```

Edit the `config/keto.yml` file and add the following lines:

```
namespaces:
  - id: 0
    name: notebooks
  - id: 1
    name: codeinsights
  - id: 2
    name: teams
```

Edit `docker-compose.yml` so that the volumes bind points to `./config/keto.yml` instead of just `config/keto.yml`.

## Run Keto

You can start Keto with an in-memory database by running

```docker-compose -f ./docker-compose.yml up```

## Run the Go code

After cloning this repository you should be able to run it. Run these commands in THIS cloned repo (not the keto repo you cloned above):

```
go mod download
go run .
```

It should end by printing the line "Steven can read Code Insight 1".

Reading the code in `main.go` along with the comments should be enough explanation as to what is happening.