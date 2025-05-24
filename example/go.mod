module example

go 1.22

toolchain go1.24.0

replace github.com/swaggest/jsonform-go => ../

require (
	github.com/swaggest/jsonform-go v0.1.0
	github.com/swaggest/openapi-go v0.2.58
	github.com/swaggest/rest v0.2.74
	github.com/swaggest/swgui v1.8.4
	github.com/swaggest/usecase v1.3.1
)

require (
	github.com/go-chi/chi/v5 v5.2.1 // indirect
	github.com/santhosh-tekuri/jsonschema/v3 v3.1.0 // indirect
	github.com/swaggest/form/v5 v5.1.1 // indirect
	github.com/swaggest/jsonschema-go v0.3.78 // indirect
	github.com/swaggest/refl v1.4.0 // indirect
	github.com/vearutop/statigz v1.5.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
