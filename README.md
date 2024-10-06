# ldapb
LDAP phone book for SIP phones. It gathers information from other backends and translates them for hard phones
and other applications which require LDAP. Supports multi tenancy by being able to have more than one backend.

## Sample configuration
```yaml
backends:
  # user that corresponds to this backend
  grable:
    # this is a sha256 sum of user:password
    password: 80fd33ed819217ddc85d30cb3ea16942610196c46037e85300844de22fbccb23
    url: http://127.0.0.1:8080/contacts
```
