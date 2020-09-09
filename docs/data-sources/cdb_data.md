# Custom DB Data Source

Getting information from a Custom DB, based on its name.

## Example Usage

```
data "solidserver_cdb_data" "myCustomData" {
  custom_db        = "myFirstCustomDB"
  value1           = "FR"
}
```

## Argument Reference

* `custom_db` - (Required) The name of the Custom DB.
* `value1` - (Required) The value of the first column, used as a key.
## Attribute Reference

* `custom_db` - (Required)The name of the Custom DB.
* `value1` - The value of the first column, used as a key.
* `value2` - The value of the second column.
* ...
* `value10` - The value of the tenth column.