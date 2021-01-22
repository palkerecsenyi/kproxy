package certificate

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/elazarl/goproxy"
	"os"
)

func GetPublicKey() []byte {
	certPath := os.Getenv("KPROXY_CERT")
	cert, _ := os.ReadFile(certPath)
	return cert
}

func SetCA() {
	cert := GetPublicKey()

	keyPath := os.Getenv("KPROXY_KEY")
	key, _ := os.ReadFile(keyPath)

	goproxyCa, err := tls.X509KeyPair(cert, key)
	if err != nil {
		panic(err)
	}

	goproxyCa.Leaf, err = x509.ParseCertificate(goproxyCa.Certificate[0])
	if err != nil {
		panic(err)
	}

	tlsConfig := goproxy.TLSConfigFromCA(&goproxyCa)
	goproxy.GoproxyCa = goproxyCa
	goproxy.OkConnect = &goproxy.ConnectAction{Action: goproxy.ConnectAccept, TLSConfig: tlsConfig}
	goproxy.MitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: tlsConfig}
	goproxy.HTTPMitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectHTTPMitm, TLSConfig: tlsConfig}
	goproxy.RejectConnect = &goproxy.ConnectAction{Action: goproxy.ConnectReject, TLSConfig: tlsConfig}
}
