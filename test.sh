#!/bin/sh

GODEBUG=httpmuxgo121=0 go test ./... -v -race
