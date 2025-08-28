module github.com/sirgwain/craig-stars

go 1.25

replace github.com/sirgwain/craig-stars/cs => ./cs

replace github.com/sirgwain/craig-stars/test => ./test

replace github.com/sirgwain/craig-stars/proto-wasm => ./proto-wasm

replace github.com/sirgwain/craig-stars/generators/protoc-gen-wasm-go => ./generators/protoc-gen-wasm-go

replace github.com/sirgwain/craig-stars/proto/gen => ./proto/gen

require (
	github.com/spf13/cobra v1.9.1
	github.com/spf13/viper v1.20.1
	github.com/stretchr/testify v1.11.1
)

require (
	connectrpc.com/connect v1.18.1
	github.com/disgoorg/disgo v0.18.16
	github.com/disgoorg/snowflake/v2 v2.0.3
	github.com/fabien-marty/slog-helpers v0.0.0-20240624063600-773d61849b89
	github.com/go-chi/chi/v5 v5.2.3
	github.com/go-chi/render v1.0.3
	github.com/go-pkgz/auth/v2 v2.0.0
	github.com/go-pkgz/rest v1.20.4
	github.com/golang-jwt/jwt/v5 v5.3.0
	github.com/golang-migrate/migrate/v4 v4.18.3
	github.com/golodash/godash v1.3.0
	github.com/lensesio/tableprinter v0.0.0-20201125135848-89e81fc956e7
	github.com/magefile/mage v1.15.0
	github.com/phsym/console-slog v0.3.1
	github.com/samber/lo v1.51.0
	github.com/samber/slog-multi v1.4.1
	github.com/simukti/sqldb-logger v0.0.0-20230108155151-646c1a075551
	github.com/sirgwain/craig-stars/cs v0.0.0-00010101000000-000000000000
	github.com/sirgwain/craig-stars/proto-wasm v0.0.0-00010101000000-000000000000
	github.com/sirgwain/craig-stars/test v0.0.0-20250827144401-54aa2fd68ef6
	google.golang.org/protobuf v1.36.8
)

require (
	buf.build/gen/go/bufbuild/bufplugin/protocolbuffers/go v1.36.6-20250718181942-e35f9b667443.1 // indirect
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.6-20250717185734-6c6e0d3c608e.1 // indirect
	buf.build/gen/go/bufbuild/registry/connectrpc/go v1.18.1-20250721151928-2b7ae473b098.1 // indirect
	buf.build/gen/go/bufbuild/registry/protocolbuffers/go v1.36.6-20250721151928-2b7ae473b098.1 // indirect
	buf.build/gen/go/pluginrpc/pluginrpc/protocolbuffers/go v1.36.6-20241007202033-cf42259fcbfc.1 // indirect
	buf.build/go/app v0.1.0 // indirect
	buf.build/go/bufplugin v0.9.0 // indirect
	buf.build/go/interrupt v1.1.0 // indirect
	buf.build/go/protovalidate v0.14.0 // indirect
	buf.build/go/protoyaml v0.6.0 // indirect
	buf.build/go/spdx v0.2.0 // indirect
	buf.build/go/standard v0.1.0 // indirect
	cel.dev/expr v0.24.0 // indirect
	connectrpc.com/otelconnect v0.7.2 // indirect
	dario.cat/mergo v1.0.1 // indirect
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20250102033503-faa5f7b0171c // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/air-verse/air v1.61.7 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.1 // indirect
	github.com/bep/godartsass v1.2.0 // indirect
	github.com/bep/godartsass/v2 v2.1.0 // indirect
	github.com/bep/golibsass v1.2.0 // indirect
	github.com/bitfield/gotestdox v0.2.2 // indirect
	github.com/bufbuild/buf v1.56.0 // indirect
	github.com/bufbuild/protocompile v0.14.1 // indirect
	github.com/bufbuild/protoplugin v0.0.0-20250218205857-750e09ce93e1 // indirect
	github.com/cli/safeexec v1.0.1 // indirect
	github.com/containerd/errdefs v1.0.0 // indirect
	github.com/containerd/errdefs/pkg v0.3.0 // indirect
	github.com/containerd/stargz-snapshotter/estargz v0.17.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.7 // indirect
	github.com/creack/pty v1.1.23 // indirect
	github.com/cubicdaiya/gonp v1.0.4 // indirect
	github.com/dave/jennifer v1.6.0 // indirect
	github.com/disgoorg/json v1.2.0 // indirect
	github.com/distribution/reference v0.6.0 // indirect
	github.com/dnephin/pflag v1.0.7 // indirect
	github.com/docker/cli v28.3.3+incompatible // indirect
	github.com/docker/distribution v2.8.3+incompatible // indirect
	github.com/docker/docker v28.3.3+incompatible // indirect
	github.com/docker/docker-credential-helpers v0.9.3 // indirect
	github.com/docker/go-connections v0.5.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/fatih/color v1.17.0 // indirect
	github.com/fatih/structtag v1.2.0 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-openapi/jsonpointer v0.21.0 // indirect
	github.com/go-sql-driver/mysql v1.9.2 // indirect
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/gofrs/flock v0.12.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/gohugoio/hugo v0.134.3 // indirect
	github.com/google/cel-go v0.26.0 // indirect
	github.com/google/go-containerregistry v0.20.6 // indirect
	github.com/google/pprof v0.0.0-20250607225305-033d6d78b36a // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.22.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.7.4 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jdx/go-netrc v1.0.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jmattheis/goverter v1.9.0 // indirect
	github.com/kataras/tablewriter v0.0.0-20180708051242-e063d29b7c23 // indirect
	github.com/klauspost/pgzip v1.2.6 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/moby/docker-image-spec v1.3.1 // indirect
	github.com/moby/term v0.5.2 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/ncruces/go-strftime v0.1.9 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.1 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pganalyze/pg_query_go/v6 v6.1.0 // indirect
	github.com/pingcap/errors v0.11.5-0.20240311024730-e056997136bb // indirect
	github.com/pingcap/failpoint v0.0.0-20240528011301-b51a646c7c86 // indirect
	github.com/pingcap/log v1.1.0 // indirect
	github.com/pingcap/tidb/pkg/parser v0.0.0-20250324122243-d51e00e5bbf0 // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10
	github.com/quic-go/qpack v0.5.1 // indirect
	github.com/quic-go/quic-go v0.54.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/riza-io/grpc-go v0.2.0 // indirect
	github.com/rrivera/identicon v0.0.0-20240116195454-d5ba35832c0d // indirect
	github.com/rs/cors v1.11.1 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sagikazarmark/locafero v0.9.0 // indirect
	github.com/sasha-s/go-csync v0.0.0-20240107134140-fcbab37b09ad // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	github.com/segmentio/encoding v0.5.3 // indirect
	github.com/sirgwain/craig-stars/generators/protoc-gen-wasm-go v0.0.0-00010101000000-000000000000 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/sqlc-dev/sqlc v1.29.0 // indirect
	github.com/stoewer/go-strcase v1.3.1 // indirect
	github.com/tdewolff/parse/v2 v2.7.15 // indirect
	github.com/tetratelabs/wazero v1.9.0 // indirect
	github.com/vbatts/tar-split v0.12.1 // indirect
	github.com/wasilibs/go-pgquery v0.0.0-20250409022910-10ac41983c07 // indirect
	github.com/wasilibs/wazero-helpers v0.0.0-20240620070341-3dff1577cd52 // indirect
	go.lsp.dev/jsonrpc2 v0.10.0 // indirect
	go.lsp.dev/pkg v0.0.0-20210717090340-384b27a52fb2 // indirect
	go.lsp.dev/protocol v0.12.0 // indirect
	go.lsp.dev/uri v0.3.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.62.0 // indirect
	go.opentelemetry.io/otel v1.37.0 // indirect
	go.opentelemetry.io/otel/metric v1.37.0 // indirect
	go.opentelemetry.io/otel/trace v1.37.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/mock v0.5.2 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/exp v0.0.0-20250819193227-8b4c13bb791b // indirect
	golang.org/x/mod v0.27.0 // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/term v0.34.0 // indirect
	golang.org/x/tools v0.36.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250728155136-f173205681a0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250728155136-f173205681a0 // indirect
	google.golang.org/grpc v1.72.2 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gotest.tools/gotestsum v1.12.0 // indirect
	modernc.org/libc v1.62.1 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.9.1 // indirect
	modernc.org/sqlite v1.37.0 // indirect
	pluginrpc.com/pluginrpc v0.5.0 // indirect
)

require (
	cloud.google.com/go/compute/metadata v0.8.0 // indirect
	github.com/ajg/form v1.5.1 // indirect
	github.com/dghubble/oauth1 v0.7.3 // indirect
	github.com/go-oauth2/oauth2/v4 v4.5.4 // indirect
	github.com/go-pkgz/repeater v1.2.0 // indirect
	github.com/golang/snappy v1.0.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	go.etcd.io/bbolt v1.4.3 // indirect
	go.mongodb.org/mongo-driver v1.17.4 // indirect
	golang.org/x/image v0.30.0 // indirect
	golang.org/x/oauth2 v0.30.0
	golang.org/x/sync v0.16.0
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/fabien-marty/tracerr v0.0.0-20240624051446-7f090eca46ee // indirect
	github.com/nsf/jsondiff v0.0.0-20230430225905-43f6cf3098c1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/samber/slog-common v0.19.0 // indirect
	github.com/ztrue/tracerr v0.4.0 // indirect
)

require (
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/mattn/go-sqlite3 v1.14.32
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/spf13/afero v1.14.0 // indirect
	github.com/spf13/cast v1.9.2 // indirect
	github.com/spf13/pflag v1.0.7 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	golang.org/x/crypto v0.41.0
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/text v0.28.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

tool (
	connectrpc.com/connect/cmd/protoc-gen-connect-go
	github.com/air-verse/air
	github.com/bufbuild/buf/cmd/buf
	github.com/jmattheis/goverter/cmd/goverter
	github.com/planetscale/vtprotobuf/cmd/protoc-gen-go-vtproto
	github.com/sirgwain/craig-stars/generators/protoc-gen-wasm-go
	github.com/sqlc-dev/sqlc/cmd/sqlc
	google.golang.org/protobuf/cmd/protoc-gen-go
	gotest.tools/gotestsum
)
