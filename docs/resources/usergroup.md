# UserGroup Resource

Groups associate users with authorization rules and SOLIDserver resources.

## Example Usage

Creating a Group:
```
resource "solidserver_usergroup" "t_group_01" {
  name = "group01"
  description = "descr01"
}
```

## Argument Reference

* `name` - (Required) The name of the group.
* `description` - description of the group.