package solidserver

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type SOLIDserver struct {
	Host                     string
	Username                 string
	Password                 string
	BaseUrl                  string
	SSLVerify                bool
	AdditionalTrustCertsFile string
	Version                  int
}

func NewSOLIDserver(host string, username string, password string, sslverify bool, certsfile string) *SOLIDserver {
	s := &SOLIDserver{
		Host:                     host,
		Username:                 username,
		Password:                 password,
		BaseUrl:                  "https://" + host,
		SSLVerify:                sslverify,
		AdditionalTrustCertsFile: certsfile,
		Version:                  0,
	}

	s.GetVersion()

	return s
}

func (s *SOLIDserver) GetVersion() error {
	// Get the SystemCertPool, continue with an empty pool on error
	rootCAs, x509err := x509.SystemCertPool()

	if rootCAs == nil || x509err != nil {
		rootCAs = x509.NewCertPool()
	}

	if s.AdditionalTrustCertsFile != "" {
		certs, readErr := ioutil.ReadFile(s.AdditionalTrustCertsFile)
		log.Printf("[DEBUG] Certificates = %s\n", certs)

		if readErr != nil {
			log.Fatalf("Failed to append %q to RootCAs: %v\n", s.AdditionalTrustCertsFile, readErr)
		}

		log.Printf("[DEBUG] Cert Subjects Before Append = %d\n", len(rootCAs.Subjects()))

		if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
			log.Printf("No certs appended, using system certs only\n")
		}
		log.Printf("[DEBUG] Cert Subjects After Append = %d\n", len(rootCAs.Subjects()))
	}

	apiclient := gorequest.New()

	parameters := url.Values{}
	parameters.Add("WHERE", "member_is_me='1'")

	resp, body, err := apiclient.Get(fmt.Sprintf("%s/%s?%s", s.BaseUrl, "rest/member_list", parameters.Encode())).
		TLSClientConfig(&tls.Config{InsecureSkipVerify: !s.SSLVerify, RootCAs: rootCAs}).
		Set("X-IPM-Username", base64.StdEncoding.EncodeToString([]byte(s.Username))).
		Set("X-IPM-Password", base64.StdEncoding.EncodeToString([]byte(s.Password))).
		End()

	if err == nil && resp.StatusCode == 200 {
		var buf [](map[string]interface{})
		json.Unmarshal([]byte(body), &buf)

		if version, versionExist := buf[0]["member_version"].(string); versionExist {
			log.Printf("[DEBUG] SOLIDServer - Version: %s\n", version)

			StrVersion := strings.Split(version, ".")

			for i := 0; i < 3; i++ {
				num, numErr := strconv.Atoi(StrVersion[i])
				if numErr == nil {
					s.Version = s.Version*10 + num
				} else {
					s.Version = s.Version*10 + 0
				}
			}

			log.Printf("[DEBUG] SOLIDServer - Version: %d\n", s.Version)

			return nil
		}
	}

	return fmt.Errorf("SOLIDServer - Error retrieving SOLIDserver Version\n")
}

func (s *SOLIDserver) Request(method string, service string, parameters *url.Values) (*http.Response, string, error) {
	var resp *http.Response = nil
	var body string = ""
	var err []error = nil

	// Get the SystemCertPool, continue with an empty pool on error
	rootCAs, x509err := x509.SystemCertPool()

	if rootCAs == nil || x509err != nil {
		rootCAs = x509.NewCertPool()
	}

	if s.AdditionalTrustCertsFile != "" {
		certs, readErr := ioutil.ReadFile(s.AdditionalTrustCertsFile)
		log.Printf("[DEBUG] Certificates = %s\n", certs)

		if readErr != nil {
			log.Fatalf("Failed to append %q to RootCAs: %v\n", s.AdditionalTrustCertsFile, readErr)
		}

		log.Printf("[DEBUG] Cert Subjects Before Append = %d\n", len(rootCAs.Subjects()))

		if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
			log.Printf("No certs appended, using system certs only\n")
		}
		log.Printf("[DEBUG] Cert Subjects After Append = %d\n", len(rootCAs.Subjects()))
	}

	apiclient := gorequest.New()

	// Set gorequest options
	apiclient.Timeout(16 * time.Second)
	apiclient.Retry(3, time.Duration(rand.Intn(128)) * time.Millisecond, http.StatusTooManyRequests)

	switch method {
	case "post":
		// Random Delay for write operation to distribute the load
		time.Sleep(time.Duration(rand.Intn(16)) * time.Millisecond)
		resp, body, err = apiclient.Post(fmt.Sprintf("%s/%s?%s", s.BaseUrl, service, parameters.Encode())).
			TLSClientConfig(&tls.Config{InsecureSkipVerify: !s.SSLVerify, RootCAs: rootCAs}).
			Set("X-IPM-Username", base64.StdEncoding.EncodeToString([]byte(s.Username))).
			Set("X-IPM-Password", base64.StdEncoding.EncodeToString([]byte(s.Password))).
			End()
	case "put":
		// Random Delay for write operation to distribute the load
		time.Sleep(time.Duration(rand.Intn(16)) * time.Millisecond)
		resp, body, err = apiclient.Put(fmt.Sprintf("%s/%s?%s", s.BaseUrl, service, parameters.Encode())).
			TLSClientConfig(&tls.Config{InsecureSkipVerify: !s.SSLVerify, RootCAs: rootCAs}).
			Set("X-IPM-Username", base64.StdEncoding.EncodeToString([]byte(s.Username))).
			Set("X-IPM-Password", base64.StdEncoding.EncodeToString([]byte(s.Password))).
			End()
	case "delete":
		// Random Delay for write operation to distribute the load
		time.Sleep(time.Duration(rand.Intn(16)) * time.Millisecond)
		resp, body, err = apiclient.Delete(fmt.Sprintf("%s/%s?%s", s.BaseUrl, service, parameters.Encode())).
			TLSClientConfig(&tls.Config{InsecureSkipVerify: !s.SSLVerify, RootCAs: rootCAs}).
			Set("X-IPM-Username", base64.StdEncoding.EncodeToString([]byte(s.Username))).
			Set("X-IPM-Password", base64.StdEncoding.EncodeToString([]byte(s.Password))).
			End()
	case "get":
		resp, body, err = apiclient.Get(fmt.Sprintf("%s/%s?%s", s.BaseUrl, service, parameters.Encode())).
			TLSClientConfig(&tls.Config{InsecureSkipVerify: !s.SSLVerify, RootCAs: rootCAs}).
			Set("X-IPM-Username", base64.StdEncoding.EncodeToString([]byte(s.Username))).
			Set("X-IPM-Password", base64.StdEncoding.EncodeToString([]byte(s.Password))).
			End()
	default:
		return nil, "", fmt.Errorf("SOLIDServer - Error initiating API call, unsupported HTTP request\n")
	}

	if err != nil {
		return nil, "", fmt.Errorf("SOLIDServer - Error initiating API call (%q)\n", err)
	}

	return resp, body, nil
}
