# ldapb
LDAP phone book for SIP phones. It gathers information from other backends and translates them for hard phones
and other applications which require LDAP. Supports multi tenancy by being able to have more than one backend.

## Sample configuration
```yaml
backends:
  # user that corresponds to this backend
  grable:
    # this is a sha256 sum of user:password
    # the password is 'super-secret'
    password: 2266fb05e611d31b1ae33d8d3c39841c48f6f70b2496b8882eacb491928478a7
    url: http://127.0.0.1:8080/contacts
```

## Responses
Currently, all response records consist of three attributes: dn, cn and telephoneNumber. The IDF looks something
like this:

```
dn: uname=zinnia,
cn: Zinnia Elegans
telephoneNumber: 1000
```

## Filters
Only `cn` and `telephoneNumber` are supported in search filters. Any other field will simply be ignored. OR/AND filters
are supported as well. Here are some examples of supported filters:

```
(|(cn=1001*)(telephoneNumber=1001*))
(cn=hello)
(&(cn=John)(telephoneNumber=5000))
```

## Test Me
You can run a local server and test it by installing the OpenLDAP utilities, specifically `ldapsearch`. For example,
the sample configuration can be tested with (considering that the service is running in the same
computer): `ldapsearch -x -D grable -w super-secret -H ldap://127.0.0.1:389`.
