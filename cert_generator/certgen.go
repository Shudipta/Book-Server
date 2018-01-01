package cert_generator

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"
	"log"
	//"k8s.io/client-go/util/cert"
	//"k8s.io/client-go/util/cert"
)

const (
	organization = "AppsCode Inc."
	commonName  = "Book Server"
	duration	= 365
	caCertFilename  = "cert_generator/ca.crt"
	caKeyFilename	= "cert_generator/ca.key"
	srvCertFilename  = "cert_generator/srv.crt"
	srvKeyFilename	= "cert_generator/srv.key"
	clCertFilename  = "cert_generator/cl.crt"
	clKeyFilename	= "cert_generator/cl.key"
	isClient	= false
)

var (
	caCert, srvCert, clCert *x509.Certificate
	priv, caPriv *rsa.PrivateKey
)

func newCertificate(organization, commonName string, duration int, check int, addresses []string) *x509.Certificate {
	certificate := x509.Certificate{
		Subject: pkix.Name{
			Organization: []string{organization},
			CommonName:   commonName,
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Duration(duration) * time.Hour * 24),
		BasicConstraintsValid: true,
	}
	if check == 1 {
		certificate.KeyUsage = x509.KeyUsageCertSign | x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature
		certificate.IsCA = true
		return &certificate
	} else {
		certificate.KeyUsage = x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature
		if check == 2 {
			certificate.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
		} else {
			certificate.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
		}
	}
	//
	for i := 0; i < len(addresses); i++ {
		if ip := net.ParseIP(addresses[i]); ip != nil {
			certificate.IPAddresses = append(certificate.IPAddresses, ip)
		} else {
			certificate.DNSNames = append(certificate.DNSNames, addresses[i])
		}
	}

	return &certificate
}

func generate(certificate, parent x509.Certificate, certFilename, keyFilename string, isCA bool) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal("Failed to generate private key:", err)
	}
	var parKey *rsa.PrivateKey
	if isCA {
		caPriv = priv
		parKey = priv
	} else {
		parKey = caPriv
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	certificate.SerialNumber, err = rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatal("Failed to generate serial number:", err)
	}
	//serverCert, err := cert.NewSignedCert(cfgForServer, serverKey, caCert, caKey)
	//cert.NewSelfSignedCACert(cfg, caKey)
	//cert.NewSignedCert(cfgForServer, serverKey, caCert, caKey)
	derBytes, err := x509.CreateCertificate(rand.Reader, &certificate, &parent, &priv.PublicKey, parKey)
	if err != nil {
		log.Fatal("Failed to create certificate:", err)
	}

	certOut, err := os.Create(certFilename)
	defer certOut.Close()
	if err != nil {
		log.Fatal("Failed to open "+certFilename+" for writing:", err)
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	// permission 0600 means owner can read and write file
	keyOut, err := os.OpenFile(keyFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	defer keyOut.Close()
	if err != nil {
		log.Fatal("Failed to open key "+keyFilename+" for writing:", err)
	}
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	fmt.Println("Certificate generated successfully")
	fmt.Println("\tCertificate: ", certFilename)
	fmt.Println("\tPrivate Key: ", keyFilename)
}

func caCertPair()  {
	caCert = newCertificate(organization, commonName, duration, 1, []string{})
	generate(*caCert, *caCert, caCertFilename, caKeyFilename, true)
}

func srvCertPair()  {
	addresses := []string{"localhost", "127.0.0.1"}
	srvCert = newCertificate(organization, commonName, duration, 2, addresses)
	generate(*srvCert, *caCert, srvCertFilename, srvKeyFilename, false)
}

func clCertPair()  {
	addresses := []string{"localhost", "127.0.0.1"}
	srvCert = newCertificate(organization, commonName, duration, 3, addresses)
	generate(*srvCert, *caCert, clCertFilename, clKeyFilename, false)
}

func CertGenerate() {
	//organizaion := flag.String("org", ORGANIZATION_DEFAULT, "CA Organization nane")
	//commonName := flag.String("name", COMMON_NAME_DEFAULT, "The subject name. Usually the DNS server name")
	//duration := flag.Int("duration", DURATION_DEFAULT, "How log the certificate will be valid.")
	//certFilename := flag.String("cert", CERTIFICATE_DEFAULT, "Certificate filename.")
	//keyFilename := flag.String("key", KEY_DEFAULT, "Privake Key filename.")
	//isClientCert := flag.Bool("client", IS_CLIENT_DEFAULT, "If the certificate usage is client. Default is false (server usage)")
	//flag.Parse()

	// certificate := newCertificate(organization, commonName, duration, isClient, addresses)
	// generate(*certificate, certFilename, keyFilename)
	caCertPair()
	srvCertPair()
	clCertPair()
}
