package solidserver

import (
  "github.com/alexissavin/gorequest"
  "encoding/base64"
  "crypto/tls"
  "net/http"
  "net/url"
  "fmt"
)

type SOLIDserver struct {
  Host        string
  Username    string
  Password    string
  BaseUrl     string
  SSLVerify   bool
}

func New(host string, username string, password string, sslverify bool) (*SOLIDserver) {
  s := &SOLIDserver{
    Host:      host,
    Username:  username,
    Password:  password,
    BaseUrl:   "https://" + host + "/rest/",
    SSLVerify: sslverify,
  }

  return s
}

func (s *SOLIDserver) Request(method string, service string, parameters *url.Values) (*http.Response, string, []error) {
  apiclient := gorequest.New()

  switch method {
  case "post":
    return apiclient.Post(fmt.Sprintf("%s%s/%s", s.BaseUrl, service, parameters.Encode())).
    TLSClientConfig(&tls.Config{InsecureSkipVerify: !s.SSLVerify}).
    Set("X-IPM-Username", base64.StdEncoding.EncodeToString([]byte(s.Username))).
    Set("X-IPM-Password", base64.StdEncoding.EncodeToString([]byte(s.Password))).
    End()
  case "put":
    return apiclient.Put(fmt.Sprintf("%s%s/%s", s.BaseUrl, service, parameters.Encode())).
    TLSClientConfig(&tls.Config{InsecureSkipVerify: !s.SSLVerify}).
    Set("X-IPM-Username", base64.StdEncoding.EncodeToString([]byte(s.Username))).
    Set("X-IPM-Password", base64.StdEncoding.EncodeToString([]byte(s.Password))).
    End()
  case "delete":
    return apiclient.Delete(fmt.Sprintf("%s%s/%s", s.BaseUrl, service, parameters.Encode())).
    TLSClientConfig(&tls.Config{InsecureSkipVerify: !s.SSLVerify}).
    Set("X-IPM-Username", base64.StdEncoding.EncodeToString([]byte(s.Username))).
    Set("X-IPM-Password", base64.StdEncoding.EncodeToString([]byte(s.Password))).
    End()
  case "get":
    return apiclient.Get(fmt.Sprintf("%s%s/%s", s.BaseUrl, service, parameters.Encode())).
    TLSClientConfig(&tls.Config{InsecureSkipVerify: !s.SSLVerify}).
    Set("X-IPM-Username", base64.StdEncoding.EncodeToString([]byte(s.Username))).
    Set("X-IPM-Password", base64.StdEncoding.EncodeToString([]byte(s.Password))).
    End()
  default:
    return nil, "", nil
  }
}
