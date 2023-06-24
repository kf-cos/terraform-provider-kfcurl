package provider

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

const TerraformProviderProductUserAgent = "terraform-provider-kfcurl"

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	Client HTTPClient
)

func init() {
	Client = &http.Client{}
}

type TLSClient struct {
	client *http.Client
}

func (tc *TLSClient) Do(req *http.Request) (*http.Response, error) {
	return tc.client.Do(req)
}

func NewTLSClient(certFile, keyFile, caCert, caDir string, insecureSkipVerify bool) (HTTPClient, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	var rootCAs *x509.CertPool
	if caCert != "" {
		rootCAs = x509.NewCertPool()
		caCertBytes, err := ioutil.ReadFile(caCert)
		if err != nil {
			return nil, err
		}
		if !rootCAs.AppendCertsFromPEM(caCertBytes) {
			return nil, errors.New("failed to append CA certificate")
		}
	} else if caDir != "" {
		rootCAs = x509.NewCertPool()
		files, err := ioutil.ReadDir(caDir)
		if err != nil {
			return nil, err
		}

		for _, file := range files {
			if !strings.HasSuffix(file.Name(), ".pem") {
				continue
			}

			caCert, err := ioutil.ReadFile(filepath.Join(caDir, file.Name()))
			if err != nil {
				return nil, err
			}

			if !rootCAs.AppendCertsFromPEM(caCert) {
				return nil, errors.New("failed to append CA certificate")
			}
		}
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates:       []tls.Certificate{cert},
			RootCAs:            rootCAs,
			InsecureSkipVerify: insecureSkipVerify,
		},
	}

	return &TLSClient{&http.Client{Transport: tr}}, nil
}

func setClient(certFile, keyFile, caCert, caDir string, insecureSkipVerify bool) error {
	if certFile == "" {
		return nil
	}

	tlsClient, err := NewTLSClient(certFile, keyFile, caCert, caDir, insecureSkipVerify)
	if err != nil {
		return err
	}

	Client = tlsClient
	return nil
}

func Provider() *schema.Provider {
	provider := &schema.Provider{

		DataSourcesMap: map[string]*schema.Resource{
			"kfcurl_request": dataSourceCurlRequest(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"kfcurl_request": resourceCurl(),
		},
	}

	return provider
}

type apiClient struct {
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {

		return &apiClient{}, nil
	}
}
