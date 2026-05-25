package auth

import (
	"crypto/tls"
	"fmt"
)

type ServerTLSOption struct {
	ServerCert       string
	ServerPrivatekey string
	ClientCAs        string
}

func TLSConfigProvider(tlsArtifacts *ServerTLSOption) (tlsConfig *tls.Config) {

	keyPair, err := tls.LoadX509KeyPair(tlsArtifacts.ServerCert, tlsArtifacts.ServerPrivatekey)
	if err != nil {
		msg := fmt.Sprintf("failed to parse TLS key pair - %s", err.Error())
		panic(msg)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{keyPair},
		ClientAuth:   tls.NoClientCert,
		MinVersion:   tls.VersionTLS13,
	}
}
