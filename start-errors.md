wizard@wizard-server:~$ cd uDosGo/sonic-screwdriver
wizard@wizard-server:~/uDosGo/sonic-screwdriver$ go build -o sonic ./cmd/sonic
Command 'go' not found, but can be installed with:
sudo snap install go         # version 1.26.2, or
sudo apt  install golang-go  # version 2:1.21~2
sudo apt  install gccgo-go   # version 2:1.21~2
See 'snap info go' for additional versions.
wizard@wizard-server:~/uDosGo/sonic-screwdriver$ sudo snap install go
[sudo] password for wizard: 
error: This revision of snap "go" was published using classic confinement
       and thus may perform arbitrary system changes outside of the security
       sandbox that snaps are usually confined to, which may put your system at
       risk.

       If you understand and want to proceed repeat the command including
       --classic.
wizard@wizard-server:~/uDosGo/sonic-screwdriver$ sudo snap install go
error: This revision of snap "go" was published using classic confinement
       and thus may perform arbitrary system changes outside of the security
       sandbox that snaps are usually confined to, which may put your system at
       risk.

       If you understand and want to proceed repeat the command including
       --classic.
wizard@wizard-server:~/uDosGo/sonic-screwdriver$ sudo snap install go --classic
go 1.26.2 from Canonical✓ installed
wizard@wizard-server:~/uDosGo/sonic-screwdriver$ go build -o sonic ./cmd/sonic
go: downloading github.com/docker/docker v28.5.2+incompatible
go: downloading github.com/qri-io/jsonschema v0.2.1
go: downloading gopkg.in/yaml.v3 v3.0.1
go: downloading github.com/mattn/go-sqlite3 v1.14.42
go: downloading github.com/docker/go-connections v0.7.0
go: downloading github.com/docker/go-units v0.5.0
go: downloading github.com/moby/docker-image-spec v1.3.1
go: downloading github.com/opencontainers/image-spec v1.1.1
go: downloading github.com/containerd/errdefs v1.0.0
go: downloading github.com/opencontainers/go-digest v1.0.0
go: downloading github.com/containerd/errdefs/pkg v0.3.0
go: downloading github.com/distribution/reference v0.6.0
go: downloading github.com/pkg/errors v0.9.1
go: downloading go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.68.0
go: downloading go.opentelemetry.io/otel/trace v1.43.0
go: downloading go.opentelemetry.io/otel v1.43.0
go: downloading github.com/qri-io/jsonpointer v0.1.1
go: downloading github.com/felixge/httpsnoop v1.0.4
go: downloading go.opentelemetry.io/otel/metric v1.43.0
go: downloading github.com/go-logr/logr v1.4.3
go: downloading github.com/go-logr/stdr v1.2.2
go: downloading go.opentelemetry.io/auto/sdk v1.2.1
go: downloading github.com/cespare/xxhash/v2 v2.3.0
# github.com/sonic-family/sonic-screwdriver/internal/classicmodern
internal/classicmodern/readiness.go:98:34: c.MissingPackages undefined (type *CheckReadiness has no field or method MissingPackages)
internal/classicmodern/readiness.go:104:11: c.MissingPackages undefined (type *CheckReadiness has no field or method MissingPackages)
internal/classicmodern/readiness.go:105:40: c.MissingPackages undefined (type *CheckReadiness has no field or method MissingPackages)
internal/classicmodern/readiness.go:106:71: c.MissingPackages undefined (type *CheckReadiness has no field or method MissingPackages)
internal/classicmodern/readiness.go:411:7: declared and not used: data
# github.com/sonic-family/sonic-screwdriver/internal/homeassistant
internal/homeassistant/integration.go:14:2: "strings" imported and not used
internal/homeassistant/integration.go:354:2: declared and not used: httpClient
# github.com/sonic-family/sonic-screwdriver/internal/secrets
internal/secrets/secret_store.go:71:4: undefined: log
internal/secrets/secret_store.go:299:19: undefined: fmt
internal/secrets/tui.go:123:4: undefined: restoreSecrets
internal/secrets/tui.go:736:12: undefined: json

wizard@wizard-server:~/uDosGo/sonic-screwdriver$ 
wizard@wizard-server:~/uDosGo/sonic-screwdriver$ 
