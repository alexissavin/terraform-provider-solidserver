# Custom DB Data Data Source

Getting information from a Custom DB Data, based on the first column value.

## Example Usage

```
data "solidserver_cdb" "myFirstCustomDB" {
  name             = "myFirstCustomDB"
}
```

## Argument Reference

* `name` - (Required) The name of the Custom DB.

## Attribute Reference

* `name` - The name of the Custom DB.
* `label1` - The name of the first column.
* `label2` - The name of the second column.
* ...
* `label10` - The name of the tenth column.