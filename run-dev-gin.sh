#!/bin/sh

go install github.com/codegangsta/gin@latest
gin --appPort 8000 run main.go