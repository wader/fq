package ciphersuites

//go:generate sh -c "cat tls-parameters-4.csv other.csv | ./ciphersuites.jq | gofmt -s > ciphersuites_gen.go"
