package solidserver

import (
  "crypto/tls"
  "encoding/base64"
  "github.com/go-resty/resty"
  "log"
)

type Config struct {
  Host       string
  Username   string
  Password   string
  SSLVerify  bool
}

func (c *Config) APIClient() (*resty.Client, error) {

  client := resty.New()

  client.SetTLSClientConfig(&tls.Config{ InsecureSkipVerify: !c.SSLVerify })
  client.SetHostURL("http://" + c.Host)
  // Trying to force header case - not working
  //client.Header["X-IPM-Username"] = []string{base64.StdEncoding.EncodeToString([]byte(c.Username))}
  //client.Header["X-IPM-Password"] = []string{base64.StdEncoding.EncodeToString([]byte(c.Password))}

  client.SetHeaders(map[string]string{
    "X-IPM-Username": base64.StdEncoding.EncodeToString([]byte(c.Username)),
    "X-IPM-Password": base64.StdEncoding.EncodeToString([]byte(c.Password)),
  })

  log.Printf("[DEBUG] SOLIDserver Client : %#v", client)

  return client, nil
}
