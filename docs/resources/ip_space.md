# IP Space Resource

IP Space resource allows to create spaces.

## Example Usage

Creating an IP Space:
```
resource "solidserver_ip_space" "myFirstSpace" {
  name   = "myFirstSpace"
  class  = "CUSTOM_SPACE"
  class_parameters {
    LOCATION = "PARIS"
  }
}
```

## Argument Reference

* `name` - (Required) The name of the IP Space to create.
* `class` - (Optional) An optional object class name allowing to store and display custom meta-data.
* `class_parameters` - (Optional) An optional object class parameters allowing to store and display custom meta-data as key/value.

## Attribute Reference

* `id` - The id of the IP Space.
* `name` - The name of the IP Space.
* `class` - The class name of the IP Space.
* `class_parameters` - The class parameters of the IP Space.