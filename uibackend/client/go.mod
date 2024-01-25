module github.com/openclarity/vmclarity/uibackend/client

go 1.21.4

require (
	github.com/oapi-codegen/runtime v1.1.1
	github.com/openclarity/vmclarity/uibackend/types v0.0.0-00010101000000-000000000000
)

require (
	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
	github.com/google/uuid v1.5.0 // indirect
)

replace github.com/openclarity/vmclarity/uibackend/types => ../types
