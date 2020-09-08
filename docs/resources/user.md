# User Resource

Users can connect through Web GUI and use API(s).

## Example Usage

Creating a User:
```
resource "solidserver_user" "myFirstUser" {
   login = "jsmith"
   password = "a_very_c0mpl3x_P@ssw0rd"
   description = "My Very First User Resource"
   last_name = "Smith"
   first_name = "John"
   email = "j.smith@efficientip.com"
   groups = [ "${solidserver_usergroup.grp_admin.id}" ]
}
```

## Argument Reference

* `login` - (Required) The login of the user
* `password` - (Required) The password of the user
* `groups` - (Required) A list of groups the user belongs to
* `description` - The description of the user
* `last_name` - The last name of the user
* `first_name` - The first name of the user
* `email` - The email address of the user