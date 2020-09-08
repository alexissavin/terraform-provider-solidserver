# Application Resource

Application resource allows to create an application that can be used to implement a traffic policy used by SOLIDserver GSLB(s).

## Example Usage

Creating an Application:
```
resource "solidserver_app_application" "myFirstApplicaton" {
  name         = "MyFirsApp"
  fqdn         = "myfirstapp.priv"
  gslb_members = ["ns0.priv", "ns1.priv"]
  class        = "INTERNAL_APP"
  class_parameters {
    owner = "MR. Smith"
    contact = "a.smith@mycompany.priv"
  }
}
```

## Argument Reference

* `name` - (Required) The name of the application to create.
* `fqdn` - (Optional) The Fully Qualified Domain Name of the application to create.
* `gslb_members` - (Optional) The names of the GSLB servers applying the application traffic policy.
* `class` - (Optional) An optional object class name allowing to store and display custom meta-data.
* `class_parameters` - (Optional) An optional object class parameters allowing to store and display custom meta-data as key/value.