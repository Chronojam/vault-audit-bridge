Vault Audit Bridge
---

Lets you write your vault audit logs to google cloud datastore.

You'll need to get a google cloud keyfile, attached to a service account with the datastore user permission

``` bash
$ GOOGLE_APPLICATION_CREDENTIALS=/path/to/keyfile.json ./vault-audit-bridge -google.project=MYPROJECT -datastore.entity=MY_ENTITY
$ vault audit-enable socket address="127.0.0.1:3333" socket_type="tcp"
```

or from docker

```bash
$ docker run -d -e GOOGLE_APPLICATION_CREDENTIALS=/path/to/keyfile.json -v /path/to/keyfile.json:/path/to/keyfile.json:ro quay.io/chronojam/vault-audit-bridge:latest
```
