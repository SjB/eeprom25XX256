#!/bin/sh

GOARCH=arm make all
scp memload memdump root@$1:. 
