## topogen server

The web server for UI of Topology Generator

### Synopsis

The web server for the UI of Topology Generator.

```
topogen server [flags]
```

### Options

```
  -e, --enable-injection string   Enable injection or not (default "true")
  -h, --help                      help for server
  -i, --image string              Image tag name (default "quay.io/leandroberetta/topogen")
  -j, --injection-label string    Injection Label (default "istio-injection:enabled")
      --lcpu string               Mimik Limit CPU (default "200m")
      --lmem string               Mimik Limit Memory (default "256Mi")
  -N, --name string               Name of the instance (default "mimik")
  -p, --port int                  The port of Web Server (default 8080)
  -C, --proxycpu string           IstioProxy Request CPU (default "50m")
  -M, --proxymem string           IstioProxy Request Memory (default "128Mi")
      --rcpu string               Mimik Request CPU (default "25m")
      --replica int               Number of Replicas created (default 1)
      --rmem string               Mimik Request Memory (default "64Mi")
  -v, --version string            Image version (default "v0.0.1")
```

### SEE ALSO

* [topogen](topogen.md)	 - Generate Complex Topology by UI or Commands

###### Auto generated by spf13/cobra on 5-Apr-2022
