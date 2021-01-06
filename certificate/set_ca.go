package certificate

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/elazarl/goproxy"
	"io/ioutil"
	"os"
)

func SetCA() {
	certPath := os.Getenv("KPROXY_CERT")
	keyPath := os.Getenv("KPROXY_KEY")

	cert, _ := ioutil.ReadFile(certPath)
	key, _ := ioutil.ReadFile(keyPath)

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
