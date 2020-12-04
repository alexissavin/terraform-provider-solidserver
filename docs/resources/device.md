# Device Resource

Device resource allows to track devices on the network and link them with IP addresses.

## Example Usage

Creating a Device:
```
resource "solidserver_device" "my_first_device" {
  name   = "my_device"
  class  = "AWS_EC2_INSTANCE"
  class_parameters {
    cloudaz = "eu-west-1a"
    instanceid = "i-03d4bd36f915b0322"
    instancetype = "t2.micro"
  }
}
```

## Argument Reference

* `name` - (Required) The name of the Device to create.
* `class` - (Optional) An optional object class name allowing to store and display custom meta-data.
* `class_parameters` - (Optional) An optional object class parameters allowing to store and display custom meta-data as key/value.

## Attribute Reference

* `id` - An internal id.
