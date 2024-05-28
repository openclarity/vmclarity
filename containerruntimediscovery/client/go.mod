module github.com/openclarity/vmclarity/containerruntimediscovery/client

go 1.22.2

require (
	github.com/openclarity/vmclarity/api/types v0.7.0
	github.com/openclarity/vmclarity/containerruntimediscovery/types v0.7.0
)

require (
	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
	github.com/deepmap/oapi-codegen/v2 v2.1.0 // indirect
	github.com/getkin/kin-openapi v0.123.0 // indirect
	github.com/go-openapi/jsonpointer v0.20.2 // indirect
	github.com/go-openapi/swag v0.22.8 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/invopop/yaml v0.2.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/oapi-codegen/runtime v1.1.1 // indirect
	github.com/openclarity/vmclarity/core v0.7.0 // indirect
	github.com/openclarity/vmclarity/plugins/sdk-go v0.7.0 // indirect
	github.com/perimeterx/marshmallow v1.1.5 // indirect
	golang.org/x/mod v0.14.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/tools v0.17.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/openclarity/vmclarity/api/types => ../../api/types
	github.com/openclarity/vmclarity/containerruntimediscovery/types => ../types
	github.com/openclarity/vmclarity/core => ../../core
	github.com/openclarity/vmclarity/plugins/sdk-go => ../../plugins/sdk-go
)
