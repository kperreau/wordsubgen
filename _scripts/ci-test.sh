#!/bin/bash

go test -race -vet=off ./... -tags test | grep -v '\[no test files\]'
exit "${PIPESTATUS[0]}"