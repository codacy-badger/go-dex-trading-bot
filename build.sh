#!/usr/bin/env bash
GOOS="linux" GOARCH="amd64" go build -o bot
GOOS="windows" GOARCH="386" go build -o bot_win

