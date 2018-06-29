package solidserver

import (
  "github.com/alexissavin/gorequest"
  "encoding/base64"
  "crypto/tls"
  "net/http"
  "net/url"
  "fmt"
  //"log"
)

type SOLIDserver struct {
  Host        string
  Username    string
  Password    string
  BaseUrl     string
  SSLVerify   bool
}

func NewSOLIDserver(host string, username string, password string, sslverify bool) (*SOLIDserver) {
  s := &SOLIDserver{
    Host:      host,
    Username:  username,
    Password:  password,
    BaseUrl:   "https://" + host,
    SSLVerify: sslverify,
  }

  return s
}

func (s *SOLIDserver) Request(method string, service string, parameters *url.Values) (*http.Response, string, error) {
  var response *http.Response = nil
  var body string = ""
  var err []error = nil

  apiclient := gorequest.New()

  switch method {
  case "post":
    response, body, err = apiclient.Post(fmt.Sprintf("%s/%s?%s", s.BaseUrl, service, parameters.Encode())).
    TLSClientConfig(&tls.Config{InsecureSkipVerify: !s.SSLVerify}).
    Set("X-IPM-Username", base64.StdEncoding.EncodeToString([]byte(s.Username))).
    Set("X-IPM-Password", base64.StdEncoding.EncodeToString([]byte(s.Password))).
    End()
  case "put":
    response, body, err = apiclient.Put(fmt.Sprintf("%s/%s?%s", s.BaseUrl, service, parameters.Encode())).
    TLSClientConfig(&tls.Config{InsecureSkipVerify: !s.SSLVerify}).
    Set("X-IPM-Username", base64.StdEncoding.EncodeToString([]byte(s.Username))).
    Set("X-IPM-Password", base64.StdEncoding.EncodeToString([]byte(s.Password))).
    End()
  case "delete":
    response, body, err = apiclient.Delete(fmt.Sprintf("%s/%s?%s", s.BaseUrl, service, parameters.Encode())).
    TLSClientConfig(&tls.Config{InsecureSkipVerify: !s.SSLVerify}).
    Set("X-IPM-Username", base64.StdEncoding.EncodeToString([]byte(s.Username))).
    Set("X-IPM-Password", base64.StdEncoding.EncodeToString([]byte(s.Password))).
    End()
  case "get":
    response, body, err = apiclient.Get(fmt.Sprintf("%s/%s?%s", s.BaseUrl, service, parameters.Encode())).
    TLSClientConfig(&tls.Config{InsecureSkipVerify: !s.SSLVerify}).
    Set("X-IPM-Username", base64.StdEncoding.EncodeToString([]byte(s.Username))).
    Set("X-IPM-Password", base64.StdEncoding.EncodeToString([]byte(s.Password))).
    End()
  default:
    return nil, "", fmt.Errorf("SOLIDServer - Error initiating API call, unsupported HTTP request")
  }

  if (err != nil) {
    return nil, "", fmt.Errorf("SOLIDServer - Error initiating API call")
  }
  
  return response, body, nil
}