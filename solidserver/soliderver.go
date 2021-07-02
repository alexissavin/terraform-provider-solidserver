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
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"reflect"
)


type HttpRequestFunc func(*gorequest.SuperAgent, string) *gorequest.SuperAgent

var httpRequestMethods = map[string]HttpRequestFunc{
	"post": (*gorequest.SuperAgent).Post,
	"put": (*gorequest.SuperAgent).Put,
	"delete": (*gorequest.SuperAgent).Delete,
	"get": (*gorequest.SuperAgent).Get,
}

var httpRequestTimings = map[string]struct {
	msSweep int
	sTimeout int
	maxTry int
}{
	"post": { msSweep:16, sTimeout:10, maxTry:1 },
	"put": { msSweep:16, sTimeout:10, maxTry:1 },
	"delete": { msSweep:16, sTimeout:10, maxTry:1 },
	"get": { msSweep:16, sTimeout:3, maxTry:6 },
}

const maxTry = 6
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

func SubmitRequest(s *SOLIDserver, apiclient *gorequest.SuperAgent, method string, service string, parameters string) (*http.Response, string, error) {
	var resp *http.Response = nil
	var body string = ""
	var errs []error = nil
	var requestUrl string = ""

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

	t := httpRequestTimings[method]
	log.Printf("[DEBUG] timings for method '%s' : {%v}\n", method, t)

	apiclient.Timeout(time.Duration(t.sTimeout) * time.Second)

	retryCount := 0

KeepTrying:
	for retryCount < t.maxTry {

		log.Printf("[DEBUG] request retryCount=%d\n", retryCount)

		httpFunc, ok := httpRequestMethods[method]
		if !ok {
			return nil, "", fmt.Errorf("SOLIDServer - Error initiating API call, unsupported HTTP request '%s'\n", method)
		}
		// Random Delay for write operation to distribute the load
		time.Sleep(time.Duration(rand.Intn(t.msSweep)) * time.Millisecond)
		requestUrl = fmt.Sprintf("%s/%s?%s", s.BaseUrl, service, parameters)
		resp, body, errs = httpFunc(apiclient, requestUrl).
			TLSClientConfig(&tls.Config{InsecureSkipVerify: !s.SSLVerify, RootCAs: rootCAs}).
			Set("X-IPM-Username", base64.StdEncoding.EncodeToString([]byte(s.Username))).
			Set("X-IPM-Password", base64.StdEncoding.EncodeToString([]byte(s.Password))).
			End()

		log.Printf("[DEBUG] checking for errors\n")
		if errs != nil {
			log.Printf("[DEBUG] '%s' API request '%s' failed with errors...\n", method, requestUrl)
			for i, err := range errs {
				log.Printf("[DEBUG] errs[%d] / (%s) = '%v'\n", i, reflect.TypeOf(err), err)
				// https://stackoverflow.com/questions/23494950/specifically-check-for-timeout-error/23497404
				if err, ok := err.(net.Error) ; ok && err.Timeout() {
					log.Printf("[DEBUG] timeout error: retrying...\n")
					retryCount++
					continue KeepTrying
				}
				log.Printf("[DEBUG] non retryable error: bailing out...\n")
			}
		}
		break KeepTrying
	}

	if retryCount >= maxTry {
		return nil, "", fmt.Errorf("SOLIDServer - [ERROR] '%s' API request '%s' : timeout retry count exceeded (maxTry = %d) !\n", method, requestUrl, maxTry)
	}

	if errs != nil {
		return nil, "", fmt.Errorf("SOLIDServer - Error initiating API call (%q)\n", errs)
	}

	return resp, body, nil
}

func (s *SOLIDserver) GetVersion(version string) error {

	apiclient := gorequest.New()

	parameters := url.Values{}
	parameters.Add("WHERE", "member_is_me='1'")

	resp, body, errs := SubmitRequest(s, apiclient, "get", "rest/member_list", parameters.Encode())

	if errs == nil && resp.StatusCode == 200 {
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

	if errs == nil && resp.StatusCode == 401 && version != "" {
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
	var err error = nil

	apiclient := gorequest.New()

	if s.Authenticated == false {
		apiclient.Retry(3, time.Duration(rand.Intn(15)+1)*time.Second, http.StatusTooManyRequests, http.StatusInternalServerError)
	} else {
		apiclient.Retry(3, time.Duration(rand.Intn(15)+1)*time.Second, http.StatusRequestTimeout, http.StatusTooManyRequests, http.StatusInternalServerError, http.StatusUnauthorized)
	}

	resp, body, err = SubmitRequest(s, apiclient, method, service, parameters.Encode())

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
