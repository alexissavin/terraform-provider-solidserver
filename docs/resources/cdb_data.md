# Custom DB Data Resource

Custom DB Data resource allows to write values in a Custom DB.

## Example Usage

Creating Custom DB Data
```
resource "solidserver_cdb_data" "myFirstCustomData" {
  custom_db = "myFirstCustomDB"
  value1 = "FR"
  value2 = "France"
}
```

## Argument Reference

* `custom_db` - (Required) The name of the Custom DB into which writing the data.
* `value1` - (Required) The value of the first column.
* `value2` - (Optional) The value of the second column.
* ...
* `value10` - (Optional) The value of the tenth column.