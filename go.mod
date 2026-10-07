module github.com/sirgwain/craig-stars

go 1.27.1

replace github.com/sirgwain/craig-stars/cs => ./cs

replace github.com/sirgwain/craig-stars/test => ./test

replace github.com/sirgwain/craig-stars/proto-wasm => ./proto-wasm

replace github.com/sirgwain/craig-stars/generators/protoc-gen-wasm-go => ./generators/protoc-gen-wasm-go

replace github.com/sirgwain/craig-stars/proto/gen => ./proto/gen

require (
	connectrpc.com/connect v1.21.0
	github.com/disgoorg/disgo v0.19.6
	github.com/disgoorg/snowflake/v2 v2.0.3
	github.com/fabien-marty/slog-helpers v0.0.0-20240624063600-773d61849b89
	github.com/go-chi/chi/v5 v5.3.2
	github.com/go-chi/render v1.0.3
	github.com/go-pkgz/auth/v2 v2.3.0
	github.com/go-pkgz/rest v1.24.0
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/golang-migrate/migrate/v4 v4.20.1
	github.com/golodash/godash v1.3.0
	github.com/lensesio/tableprinter v0.0.0-20201125135848-89e81fc956e7
	github.com/magefile/mage v1.17.2
	github.com/mark3labs/mcp-go v1.1.1
	github.com/mattn/go-sqlite3 v1.14.52
	github.com/phsym/console-slog v0.3.1
	github.com/planetscale/vtprotobuf v0.6.1-0.20240319094008-0393e58bdf10
	github.com/redpanda-data/protoc-gen-go-mcp v0.0.0-20260729122341-6334690a62b6
	github.com/samber/lo v1.53.0
	github.com/samber/slog-multi v1.8.0
	github.com/simukti/sqldb-logger v0.0.0-20230108155151-646c1a075551
	github.com/sirgwain/craig-stars/cs v0.0.0-20261006155015-1544c51b8442
	github.com/sirgwain/craig-stars/proto-wasm v0.0.0-20261006155015-1544c51b8442
	github.com/sirgwain/craig-stars/test v0.0.0-20261006155015-1544c51b8442
	github.com/spf13/cobra v1.10.2
	github.com/spf13/viper v1.21.0
	github.com/stretchr/testify v1.12.1
	golang.org/x/crypto v0.57.0
	golang.org/x/oauth2 v0.37.0
	golang.org/x/sync v0.23.0
	google.golang.org/grpc v1.84.0
	google.golang.org/protobuf v1.36.12
)

require (
	buf.build/gen/go/bufbuild/bufplugin/protocolbuffers/go v1.36.12-20261002170247-538130002972.2 // indirect
	buf.build/gen/go/bufbuild/protodescriptor/protocolbuffers/go v1.36.12-20250109164928-1da0de137947.2 // indirect
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.12-20260825204119-511051f7f437.2 // indirect
	buf.build/gen/go/bufbuild/registry/connectrpc/go v1.21.0-20261006190621-8c688359738f.1 // indirect
	buf.build/gen/go/bufbuild/registry/protocolbuffers/go v1.36.12-20261006190621-8c688359738f.2 // indirect
	buf.build/gen/go/pluginrpc/pluginrpc/protocolbuffers/go v1.36.12-20241007202033-cf42259fcbfc.2 // indirect
	buf.build/gen/go/redpandadata/common/protocolbuffers/go v1.36.12-20260323171043-6e06f84ad823.2 // indirect
	buf.build/go/app v0.2.1-0.20260824172350-b0e76892c61a // indirect
	buf.build/go/bufplugin v0.11.1 // indirect
	buf.build/go/bufprivateusage v0.1.0 // indirect
	buf.build/go/interrupt v1.1.0 // indirect
	buf.build/go/protovalidate v1.4.0 // indirect
	buf.build/go/protoyaml v0.7.0 // indirect
	buf.build/go/spdx v0.2.0 // indirect
	buf.build/go/standard v0.1.1-0.20260325175353-2b287e071df5 // indirect
	cel.dev/cel-go v0.32.0 // indirect
	cel.dev/expr v0.25.3 // indirect
	cloud.google.com/go/compute/metadata v0.10.0 // indirect
	connectrpc.com/otelconnect v0.10.0 // indirect
	dario.cat/mergo v1.0.2 // indirect
	filippo.io/edwards25519 v1.2.0 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/air-verse/air v1.67.4 // indirect
	github.com/ajg/form v1.9.0 // indirect
	github.com/andybalholm/brotli v1.2.6 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.1 // indirect
	github.com/bep/godartsass/v2 v2.5.0 // indirect
	github.com/bep/golibsass v1.2.0 // indirect
	github.com/bep/helpers v0.12.0 // indirect
	github.com/bitfield/gotestdox v0.2.3 // indirect
	github.com/bits-and-blooms/bitset v1.26.0 // indirect
	github.com/bufbuild/buf v1.73.0 // indirect
	github.com/bufbuild/protocompile v0.14.2-0.20260910151042-7436f7c76201 // indirect
	github.com/bufbuild/protoplugin v0.0.0-20260414125817-25d1d281b46b // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cli/browser v1.3.0 // indirect
	github.com/clipperhouse/uax29/v2 v2.7.0 // indirect
	github.com/containerd/errdefs v1.0.0 // indirect
	github.com/containerd/errdefs/pkg v0.3.0 // indirect
	github.com/coreos/go-semver v0.3.1 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.7 // indirect
	github.com/cubicdaiya/gonp v1.0.4 // indirect
	github.com/dave/jennifer v1.7.1 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dghubble/oauth1 v0.7.3 // indirect
	github.com/disgoorg/json/v2 v2.0.0 // indirect
	github.com/disgoorg/omit v1.0.0 // indirect
	github.com/distribution/reference v0.6.0 // indirect
	github.com/dnephin/pflag v1.0.7 // indirect
	github.com/docker/cli v29.8.2+incompatible // indirect
	github.com/docker/docker-credential-helpers v0.9.9 // indirect
	github.com/docker/go-connections v0.8.1 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/dustin/go-humanize v1.1.0 // indirect
	github.com/fabien-marty/tracerr v0.0.0-20240624051446-7f090eca46ee // indirect
	github.com/fatih/color v1.19.0 // indirect
	github.com/fatih/structtag v1.2.0 // indirect
	github.com/felixge/httpsnoop v1.1.0 // indirect
	github.com/fsnotify/fsnotify v1.10.1 // indirect
	github.com/go-logr/logr v1.4.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-oauth2/oauth2/v4 v4.6.0 // indirect
	github.com/go-pkgz/repeater/v2 v2.2.0 // indirect
	github.com/go-sql-driver/mysql v1.10.1 // indirect
	github.com/go-viper/mapstructure/v2 v2.5.0 // indirect
	github.com/gobwas/glob v1.0.0 // indirect
	github.com/gofrs/flock v0.13.1 // indirect
	github.com/gohugoio/hashstructure v1.1.0 // indirect
	github.com/gohugoio/hugo v0.167.0 // indirect
	github.com/golang/snappy v1.0.0 // indirect
	// SQLC still imports this path; v0.32 moved to cel.dev/cel-go.
	github.com/google/cel-go v0.31.0 // indirect
	github.com/google/go-containerregistry v0.22.1 // indirect
	github.com/google/jsonschema-go v0.4.3 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.11.0 // indirect
	github.com/jackc/puddle/v2 v2.2.3 // indirect
	github.com/jdx/go-netrc v1.0.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jmattheis/goverter v1.11.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kataras/tablewriter v0.0.0-20180708051242-e063d29b7c23 // indirect
	github.com/klauspost/compress v1.20.1 // indirect
	github.com/klauspost/pgzip v1.2.7 // indirect
	github.com/mattn/go-colorable v0.1.16 // indirect
	github.com/mattn/go-isatty v0.0.24 // indirect
	github.com/mattn/go-runewidth v0.0.30 // indirect
	github.com/moby/docker-image-spec v1.3.1 // indirect
	github.com/moby/moby/api v1.56.1 // indirect
	github.com/moby/moby/client v0.6.1 // indirect
	github.com/montanaflynn/stats v0.13.0 // indirect
	github.com/ncruces/go-sqlite3 v0.35.6 // indirect
	github.com/ncruces/go-sqlite3-wasm/v6 v6.3.35304 // indirect
	github.com/ncruces/julianday v1.0.0 // indirect
	github.com/nsf/jsondiff v0.0.0-20260207060731-8e8d90c4c0ac // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.1 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pelletier/go-toml/v2 v2.4.3 // indirect
	github.com/petermattis/goid v0.0.0-20260918085751-abfca077860b // indirect
	github.com/pganalyze/pg_query_go/v6 v6.2.5 // indirect
	github.com/pingcap/errors v0.11.5-0.20250523034308-74f78ae071ee // indirect
	github.com/pingcap/failpoint v0.0.0-20260811232634-55ac33a48e3b // indirect
	github.com/pingcap/log v1.1.0 // indirect
	github.com/pingcap/tidb/pkg/parser v0.0.0-20261006085058-3ca96b1d5df8 // indirect
	github.com/quic-go/qpack v0.6.0 // indirect
	github.com/quic-go/quic-go v0.63.0 // indirect
	github.com/redpanda-data/common-go/api v0.0.0-20260922181238-9566b7beb637 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/riza-io/grpc-go v0.2.0 // indirect
	github.com/rrivera/identicon v0.0.0-20240116195454-d5ba35832c0d // indirect
	github.com/rs/cors v1.11.1 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sagikazarmark/locafero v0.12.0 // indirect
	github.com/samber/slog-common v0.22.0 // indirect
	github.com/santhosh-tekuri/jsonschema/v6 v6.0.3 // indirect
	github.com/sasha-s/go-csync v0.0.0-20240107134140-fcbab37b09ad // indirect
	github.com/segmentio/asm v1.2.1 // indirect
	github.com/segmentio/encoding v0.5.4 // indirect
	github.com/sergi/go-diff v1.4.0 // indirect
	github.com/sirgwain/craig-stars/generators/protoc-gen-wasm-go v0.0.0-20261006155015-1544c51b8442 // indirect
	github.com/sirupsen/logrus v1.10.2 // indirect
	github.com/spf13/afero v1.15.0 // indirect
	github.com/spf13/cast v1.10.0 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	github.com/sqlc-dev/doubleclick v1.0.0 // indirect
	github.com/sqlc-dev/sqlc v1.31.1 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/tdewolff/parse/v2 v2.8.16 // indirect
	github.com/tetratelabs/wazero v1.12.0 // indirect
	github.com/tidwall/btree v1.8.2 // indirect
	github.com/wasilibs/go-pgquery v0.0.0-20260915022521-81f99195012b // indirect
	github.com/wasilibs/wazero-helpers v0.0.0-20250123031827-cd30c44769bb // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.2.0 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	github.com/ztrue/tracerr v0.4.0 // indirect
	go.etcd.io/bbolt v1.5.0 // indirect
	// Buf still requires the pre-v1 go.lsp.dev interfaces.
	go.lsp.dev/jsonrpc2 v0.10.0 // indirect
	go.lsp.dev/pkg v0.0.0-20210717090340-384b27a52fb2 // indirect
	go.lsp.dev/protocol v0.12.0 // indirect
	go.lsp.dev/uri v0.3.0 // indirect
	go.mongodb.org/mongo-driver v1.17.10 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.72.0 // indirect
	go.opentelemetry.io/otel v1.47.0 // indirect
	go.opentelemetry.io/otel/log v1.47.0 // indirect
	go.opentelemetry.io/otel/metric v1.47.0 // indirect
	go.opentelemetry.io/otel/trace v1.47.0 // indirect
	go.uber.org/atomic v1.12.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.28.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/exp v0.0.0-20261005173118-76772065c9b0 // indirect
	golang.org/x/image v0.46.0 // indirect
	golang.org/x/mod v0.41.0 // indirect
	golang.org/x/net v0.59.0 // indirect
	golang.org/x/sys v0.48.0 // indirect
	golang.org/x/telemetry v0.0.0-20260924152758-ed294f943157 // indirect
	golang.org/x/term v0.46.0 // indirect
	golang.org/x/text v0.42.0 // indirect
	golang.org/x/tools v0.51.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20261005182115-fad411399dd8 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20261005182115-fad411399dd8 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gotest.tools/gotestsum v1.13.0 // indirect
	mvdan.cc/xurls/v2 v2.6.0 // indirect
	pluginrpc.com/pluginrpc v0.6.0 // indirect
)

tool (
	connectrpc.com/connect/cmd/protoc-gen-connect-go
	github.com/air-verse/air
	github.com/bufbuild/buf/cmd/buf
	github.com/jmattheis/goverter/cmd/goverter
	github.com/planetscale/vtprotobuf/cmd/protoc-gen-go-vtproto
	github.com/redpanda-data/protoc-gen-go-mcp/cmd/protoc-gen-go-mcp
	github.com/sirgwain/craig-stars/generators/protoc-gen-wasm-go
	github.com/sqlc-dev/sqlc/cmd/sqlc
	golang.org/x/tools/cmd/goimports
	google.golang.org/protobuf/cmd/protoc-gen-go
	gotest.tools/gotestsum
)
