#!/bin/sh -eux

cd `dirname $0`

goapp get -u golang.org/x/tools/cmd/goimports
goapp get -u github.com/golang/lint/golint
goapp get -u github.com/constabulary/gb/...
goapp get -u github.com/PalmStoneGames/gb-gae

gb vendor restore