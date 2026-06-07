package systemconfig

//go:generate go tool oapi-codegen --include-operation-ids getSystemConfigs,getSystemConfig,putSystemConfig,postSystemConfigReset --generate types --package systemconfig -o zz_generated.systemconfig.go ../../../../../openapi/openapi.yaml
