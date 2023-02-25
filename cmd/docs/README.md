# CLI

## Upotreba

```sh
wgcli -h

wgcli [komanda] -h
```

```sh
wgcli ... --username {username} --password {password} --saveSession

wgcli ... -U {username} -P {password} -S
```

```sh
wgcli submit --problemId 1 --grader C --fileName subm.c file.c

wgcli submit -p 1 -g C -f subm.c file.c
```

- `fileName` - ime fajla koje će biti poslato (ako nije definisano uzima se ime izvornog fajla)
- `grader` - jezik koji grader koristi. Trenutno dozvoljava: C, OS, C++, Pascal, Haskell, Python, Java

## Konfiguracija

Parametri zadati kroz komandnu liniju override-uju istoimene parametre u konfiguracionim fajlovima.

### "Korisnička" konfiguracija

Default lokacija (Linux): `~/.config/wgcli/conf.yaml`

Default lokacija (Windows): `%AppData%/wgcli/conf.yaml`

```yaml
auth:
  username: "username"
  password: "password"
```

### Konfiguracija "po folderu"

Default lokacija: `./wgcli.yaml`

```yaml
auth:
# ...
problem:
  problemId: 1
  sourceFile: "..."
  fileName: "..."
  grader: "..."
```
