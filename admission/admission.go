package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"k8s.io/api/admission/v1beta1"
	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
)

var scheme = runtime.NewScheme()
var codecs = serializer.NewCodecFactory(scheme)
var deserializer = codecs.UniversalDeserializer()

func init() {
	corev1.AddToScheme(scheme)
	admissionregistrationv1beta1.AddToScheme(scheme)
}

// Config contains the server (the webhook) cert and key.
type Config struct {
	CertFile string
	KeyFile  string
}

type admitFunc func(v1beta1.AdmissionReview) *v1beta1.AdmissionResponse

func (c *Config) addFlags() {
	flag.StringVar(&c.CertFile, "tls-cert-file", c.CertFile, ""+
		"File containing the default x509 Certificate for HTTPS. (CA cert, if any, concatenated "+
		"after server cert).")
	flag.StringVar(&c.KeyFile, "tls-private-key-file", c.KeyFile, ""+
		"File containing the default x509 private key matching --tls-cert-file.")
}

func toAdmissionResponse(err error) *v1beta1.AdmissionResponse {
	return &v1beta1.AdmissionResponse{
		Result: &metav1.Status{
			Message: err.Error(),
		},
	}
}

func admitRoot(ar v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	fmt.Printf("resource: %v\n", ar.Request.Resource)
	reviewResponse := v1beta1.AdmissionResponse{}
	reviewResponse.Allowed = true
	return &reviewResponse
}

func admission(w http.ResponseWriter, r *http.Request) {
	fmt.Println("============ handling a request")
	fmt.Printf("%s\n", r.URL)

	var body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}

	// verify the content type is accurate
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		fmt.Printf("contentType=%s, expect application/json\n", contentType)
		return
	}

	var reviewResponse *v1beta1.AdmissionResponse
	ar := v1beta1.AdmissionReview{}
	if _, _, err := deserializer.Decode(body, nil, &ar); err != nil {
		reviewResponse = toAdmissionResponse(err)
	} else {
		switch ar.Request.Resource {
		case serviceAccountResource:
			reviewResponse = admitServiceAccount(ar)
		case podResource:
			reviewResponse = admitPod(ar)
		default:
			reviewResponse = admitRoot(ar)
		}
	}

	response := v1beta1.AdmissionReview{}
	if reviewResponse != nil {
		response.Response = reviewResponse
		response.Response.UID = ar.Request.UID
	}
	// reset the Object and OldObject, they are not needed in a response.
	ar.Request.Object = runtime.RawExtension{}
	ar.Request.OldObject = runtime.RawExtension{}

	resp, err := json.Marshal(response)
	if err != nil {
		fmt.Printf("marshalling: %v\n", err)
	}
	if _, err := w.Write(resp); err != nil {
		fmt.Printf("writing: %v\n", err)
	}
}

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

func main() {
	var config Config
	config.addFlags()
	flag.Parse()

	t, err := getTLSConfig()
	if err != nil {
		fmt.Printf("Failed to get TLS config: %v\n", err)
		os.Exit(1)
	}

	server := &http.Server{
		Addr:      ":443",
		TLSConfig: t,
	}

	fmt.Println("admission controller running")
	http.HandleFunc("/", admission)
	server.ListenAndServeTLS(config.CertFile, config.KeyFile)
}
