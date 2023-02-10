// Package keylog parses NSS key log format
// https://firefox-source-docs.mozilla.org/security/nss/legacy/key_log_format/index.html
// <Label> <space> <ClientRandom> <space> <Secret> lines
package keylog

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"strings"
)

const (
	Rsa = iota
	ClientRandom
	ClientEarlyTrafficSecret
	ClientHandshakeTrafficSecret
	ServerHandshakeTrafficSecret
	ClientTrafficSecret0
	ServerTrafficSecret0
	EarlyExporterSecret
	ExporterSecret
)

var labelToEnum = map[string]int{
	"RSA":                             Rsa,
	"CLIENT_RANDOM":                   ClientRandom,
	"CLIENT_EARLY_TRAFFIC_SECRET":     ClientEarlyTrafficSecret,
	"CLIENT_HANDSHAKE_TRAFFIC_SECRET": ClientHandshakeTrafficSecret,
	"SERVER_HANDSHAKE_TRAFFIC_SECRET": ServerHandshakeTrafficSecret,
	"CLIENT_TRAFFIC_SECRET_0":         ClientTrafficSecret0,
	"SERVER_TRAFFIC_SECRET_0":         ServerTrafficSecret0,
	"EARLY_EXPORTER_SECRET":           EarlyExporterSecret,
	"EXPORTER_SECRET":                 ExporterSecret,
}

type Entry struct {
	Label        int
	ClientRandom [32]byte
}

type Map map[Entry][]byte

func (m Map) Lookup(label int, clientRandom [32]byte) ([]byte, bool) {
	bs, ok := m[Entry{Label: label, ClientRandom: clientRandom}]
	return bs, ok
}

// Parse NSS Key Log format
//
// # comment
// <Label> <space> <ClientRandom> <space> <Secret>
func Parse(s string) (Map, error) {
	em := map[Entry][]byte{}

	lines := bufio.NewScanner(strings.NewReader(s))
	lineNr := 0
	for lines.Scan() {
		lineNr++
		line := strings.TrimSpace(lines.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Split(line, " ")
		if len(parts) != 3 {
			continue
			// return nil, fmt.Errorf("parse error %d: %s", lineNr, line)
		}

		label, labelOk := labelToEnum[parts[0]]
		if !labelOk {
			// return nil, fmt.Errorf("unknown label %d: %s", lineNr, parts[0])
			continue
		}

		clientRandom, clientRandomErr := hex.DecodeString(parts[1])
		if clientRandomErr != nil {
			return nil, fmt.Errorf("client random error %d: %w", lineNr, clientRandomErr)
		}
		if len(clientRandom) != 32 {
			return nil, fmt.Errorf("client random not 32 bytes%d: %s (%d)", lineNr, clientRandom, len(clientRandom))
		}

		value, valueErr := hex.DecodeString(parts[2])
		if valueErr != nil {
			return nil, fmt.Errorf("value random error %d: %w", lineNr, valueErr)
		}

		e := Entry{Label: label}
		copy(e.ClientRandom[:], clientRandom)

		if _, ok := em[e]; ok {
			return nil, fmt.Errorf("duplicate client random %d: %s", lineNr, line)
		}

		em[e] = value
	}

	return em, nil
}
