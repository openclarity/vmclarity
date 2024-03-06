module github.com/openclarity/vmclarity/uibackend/server

go 1.21.4

require (
	github.com/Portshift/go-utils v0.0.0-20220421083203-89265d8a6487
	github.com/deepmap/oapi-codegen/v2 v2.1.0
	github.com/getkin/kin-openapi v0.123.0
	github.com/google/go-cmp v0.6.0
	github.com/labstack/echo/v4 v4.11.4
	github.com/mitchellh/mapstructure v1.5.0
	github.com/oapi-codegen/echo-middleware v1.0.1
	github.com/oapi-codegen/runtime v1.1.1
	github.com/onsi/gomega v1.31.1
	github.com/openclarity/vmclarity/api/client v0.0.0-00010101000000-000000000000
	github.com/openclarity/vmclarity/api/types v0.0.0-00010101000000-000000000000
	github.com/openclarity/vmclarity/cli v0.0.0-00010101000000-000000000000
	github.com/openclarity/vmclarity/core v0.0.0-00010101000000-000000000000
	github.com/openclarity/vmclarity/uibackend/types v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/cobra v1.8.0
	github.com/spf13/viper v1.18.2
	gotest.tools/v3 v3.5.1
)

require (
	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/go-openapi/jsonpointer v0.20.2 // indirect
	github.com/go-openapi/swag v0.22.10 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/invopop/yaml v0.2.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/pelletier/go-toml/v2 v2.1.1 // indirect
	github.com/perimeterx/marshmallow v1.1.5 // indirect
	github.com/sagikazarmark/locafero v0.4.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.21.0 // indirect
	golang.org/x/exp v0.0.0-20240222234643-814bf88cf225 // indirect
	golang.org/x/mod v0.15.0 // indirect
	golang.org/x/net v0.22.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	golang.org/x/tools v0.18.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/openclarity/vmclarity/api/client => ../../api/client
	github.com/openclarity/vmclarity/api/types => ../../api/types
	github.com/openclarity/vmclarity/cli => ../../cli
	github.com/openclarity/vmclarity/core => ../../core
	github.com/openclarity/vmclarity/uibackend/types => ../types
	github.com/openclarity/vmclarity/utils => ../../utils
)
