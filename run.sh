#!/bin/bash

go get -t github.com/vslovik/form3
cd tests
go test -v -tags=integration ./integration