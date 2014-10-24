#!/bin/sh

GOARCH=arm go build $2
scp memdump root@$1:. 
