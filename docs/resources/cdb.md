# Custom DB Resource

Custom DB resource allows to write custom data in the SOLIDserver database. It can be reused by classes.

## Example Usage

Creating a Custom DB:
```
resource "solidserver_cdb" "myFirstCustomDB" {
  name = "myFirstCustomDB"
  label1 = "Country Code"
  label2 = "Country Name"
}
```

## Argument Reference

* `name` - (Required) The name of the Custom DB.
* `label1` - (Optional) The name of the first column.
* `label2` - (Optional) The name of the second column.
* ...
* `label10` - (Optional) The name of the tenth column.