# IP Space Data Source

Getting information from an IP Space, based on its name.

## Example Usage

```
data "solidserver_ip_space" "myFirstSpaceData" {
  depends_on = [solidserver_ip_space.myFirstSpace]
  name       = solidserver_ip_space.myFirstSpace.name
}
```

## Argument Reference

* `name` - (Required) The name of the IP Space.

## Attribute Reference

* `name` - The name of the IP Space.
* `class` -  The name of the class associated with the IP Space.
* `class_parameters` - The class parameters associated with the IP Space class, as key/value.