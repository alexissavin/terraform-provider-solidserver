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

const regexpIPPort = `^!?(([0-9]{1,3})\.){3}[0-9]{1,3}:[0-9]{1,5}$`
const regexpHostname = `^(([a-z0-9]|[a-z0-9][a-z0-9\-]*[a-z0-9])\.)*([a-z0-9]|[a-z0-9][a-z0-9\-]*[a-z0-9])$`
const regexpNetworkAcl = `^!?(([0-9]{1,3})\.){3}[0-9]{1,3}/[0-9]{1,2}$`

type SOLIDserver struct {
	Host                     string
	Username                 string
	Password                 string
	BaseUrl                  string
	SSLVerify                bool
	AdditionalTrustCertsFile string
	Version                  int
	Authenticated            bool
}

func NewSOLIDserver(host string, username string, password string, sslverify bool, certsfile string, version string) (*SOLIDserver, error) {
	s := &SOLIDserver{
		Host:                     host,
		Username:                 username,
		Password:                 password,
		BaseUrl:                  "https://" + host,
		SSLVerify:                sslverify,
		AdditionalTrustCertsFile: certsfile,
		Version:                  0,
		Authenticated:            false,
	}

	if err := s.GetVersion(version); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *SOLIDserver) GetVersion(version string) error {
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

		if rversion, rversionExist := buf[0]["member_version"].(string); rversionExist {
			log.Printf("[DEBUG] SOLIDServer - Version: %s\n", rversion)

			StrVersion := strings.Split(rversion, ".")

			// Computing version number
			for i := 0; i < len(StrVersion) && i < 3; i++ {
				num, numErr := strconv.Atoi(StrVersion[i])
				if numErr == nil {
					s.Version = s.Version*10 + num
				} else {
					s.Version = s.Version*10 + 0
				}
			}

			// Handling new branch version
			if s.Version < 100 {
				s.Version = s.Version * 10
			}

			log.Printf("[DEBUG] SOLIDServer - server version retrieved from remote SOLIDserver: %d\n", s.Version)

			return nil
		}
	}

	if err == nil && resp.StatusCode == 401 && version != "" {
		StrVersion := strings.Split(version, ".")

		for i := 0; i < len(StrVersion) && i < 3; i++ {
			num, numErr := strconv.Atoi(StrVersion[i])
			if numErr == nil {
				s.Version = s.Version*10 + num
			} else {
				s.Version = s.Version*10 + 0
			}
		}

		log.Printf("[DEBUG] SOLIDServer - server version retrived from local provider parameter: %d\n", s.Version)

		return nil
	}

	return fmt.Errorf("SOLIDServer - Error retrieving SOLIDserver Version (No Answer)\n")
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
	if s.Authenticated == false {
		apiclient.Retry(3, time.Duration(rand.Intn(15)+1)*time.Second, http.StatusTooManyRequests, http.StatusInternalServerError)
	} else {
		apiclient.Retry(3, time.Duration(rand.Intn(15)+1)*time.Second, http.StatusTooManyRequests, http.StatusUnauthorized, http.StatusInternalServerError)
	}

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
			Set("Cache-Control", "no-cache").
			End()
	default:
		return nil, "", fmt.Errorf("SOLIDServer - Error initiating API call, unsupported HTTP request\n")
	}

	if err != nil {
		return nil, "", fmt.Errorf("SOLIDServer - Error initiating API call (%q)\n", err)
	}

	if len(body) > 0 && body[0] == '{' && body[len(body)-1] == '}' {
		log.Printf("[DEBUG] Repacking HTTP JSON Body\n")
		body = "[" + body + "]"
	}

	if s.Authenticated == false && (200 <= resp.StatusCode && resp.StatusCode <= 204) {
		s.Authenticated = true
	}

	return resp, body, nil
}
