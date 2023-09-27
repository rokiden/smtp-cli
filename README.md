# SMTP-CLI

### Simple smtp mail sender as single static binary.

Used as notification callback.
Configured through args and/or envs.

```
Usage of ./smtpcli:
  -b string
        body
  -d    debug
  -f string
        from (default "from@smtpcli")
  -h value
        header, multiple
  -s string
        server (default "localhost:25")
  -t string
        to (default "to@smtpcli")
  -u string
        subject
```

| Flag | Env                                   | Default      | Description |
|------|---------------------------------------|--------------|-------------|
| -s   | SMTPCLI_SERVER                        | localhost:25 | SMTP server |
| -f   | SMTPCLI_FROM                          | from@smtpcli | Sender      |
| -t   | SMTPCLI_TO                            | to@smtpcli   | Recipient   |
| -u   | SMTPCLI_SUBJECT                       |              | Subject     |
| -b   | SMTPCLI_BODY                          |              | Body        |
| -h   | SMTPCLI_HEADERS_SEP + SMTPCLI_HEADERS |              | Headers     |

Flags overrides envs.

To set multiple headers:
* multiple -h flags can be set
* SMTPCLI_HEADERS can be set to separated list and SMTPCLI_HEADERS_SEP to separator

## Examples
### Healthchecks in podman
Environment:
```
SMTPCLI_SERVER=host.containers.internal:25
SMTPCLI_FROM=healthchecks@services
SMTPCLI_TO=user@local
SMTPCLI_HEADERS='Content-Type: text/plain; charset=utf-8|Content-Transfer-Encoding: quoted-printable|From: Healthchecks <healthchecks@services>'
SMTPCLI_HEADERS_SEP='|'
```
Mount: `../smtpcli/smtpcli:/smtpcli:ro`

Integration:

Execute on "down" events:
`/smtpcli -b "=F0=9F=94=B4 $NAME $STATUS $NOW"`

Execute on "up" events:
`/smtpcli -b "=E2=9C=85 $NAME $STATUS $NOW"`

Test:
```
From: Healthchecks <healthchecks@services>
To: user@local

🔴 TEST down 2023-09-27T14:59:59+00:00
```

### DIUN in podman
Environment:
```
DIUN_NOTIF_SCRIPT_CMD=sh
DIUN_NOTIF_SCRIPT_ARGS='-c,/smtpcli -b \"=F0=9F=94=B5 Image \$DIUN_ENTRY_IMAGE \$DIUN_ENTRY_STATUS at \$DIUN_ENTRY_CREATED\"'
SMTPCLI_SERVER=host.containers.internal:25
SMTPCLI_FROM=diun@services
SMTPCLI_TO=user@local
SMTPCLI_HEADERS='Content-Type: text/plain; charset=utf-8|Content-Transfer-Encoding: quoted-printable|From: Diun <diun@services>'
SMTPCLI_HEADERS_SEP='|'
```
Mount: `../smtpcli/smtpcli:/smtpcli:ro`

Test:
```
From: Diun <diun@services>
To: user@local

🔵 Image docker.io/diun/testnotif:latest new at 2020-03-26 12:23:56 +0000 UTC
```


