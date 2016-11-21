package efficientip

import (
  "crypto/tls"
  "encoding/base64"
  "github.com/go-resty/resty"
//  "log"
)

type Config struct {
  Host       string
  Username   string
  Password   string
  SSLVerify  bool
}

func (c *Config) APIClient() (*resty.Client, error) {

  Client := resty.New()

  Client.SetTLSClientConfig(&tls.Config{ InsecureSkipVerify: c.SSLVerify })
  Client.SetHostURL("https://" + c.Host + "/rpc")
  Client.SetHeaders(map[string]string{
    "X-IPM-Username": base64.StdEncoding.EncodeToString([]byte(c.Username)),
    "X-IPM-Password": base64.StdEncoding.EncodeToString([]byte(c.Password)),
  })

  return Client, nil
}
