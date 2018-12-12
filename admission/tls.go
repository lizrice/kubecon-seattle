package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"k8s.io/client-go/rest"
)

func getTLSConfig() (*tls.Config, error) {
	// get service account information
	inClusterConfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get in-cluster config: %v", err)
	}

	// Get the CA file from the service account.
	saCA, err := ioutil.ReadFile(inClusterConfig.TLSClientConfig.CAFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %v", inClusterConfig.TLSClientConfig.CAFile, err)
	}

	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(saCA) {
		return nil, fmt.Errorf("failed to append cert from serviceAccount CA file: %v", err)
	}

	// We will only trust requests with a certificate signed by this CA
	t := &tls.Config{
		ClientAuth: tls.RequireAndVerifyClientCert,
		// TODO: raise an issue as the server should be able to authenticate itself to the admission controller
		// ClientAuth:            tls.RequestClientCert,
		ClientCAs:             cp,
		VerifyPeerCertificate: CertificateChains,
	}

	return t, nil
}

// CertificateChains prints information about verified certificate chains
func CertificateChains(rawCerts [][]byte, chains [][]*x509.Certificate) error {
	if len(chains) > 0 {
		fmt.Println("Verified certificate chain from peer:")

		for _, v := range chains {
			for i, cert := range v {
				fmt.Printf("  Cert %d:\n", i)
				fmt.Printf(CertificateInfo(cert))
			}
		}
	}

	return nil
}

// CertificateInfo returns a string describing the certificate
func CertificateInfo(cert *x509.Certificate) string {
	if cert.Subject.CommonName == cert.Issuer.CommonName {
		return fmt.Sprintf("    Self-signed certificate %v\n", cert.Issuer.CommonName)
	}

	s := fmt.Sprintf("    Subject %v\n", cert.DNSNames)
	s += fmt.Sprintf("    Issued by %s\n", cert.Issuer.CommonName)
	return s
}
