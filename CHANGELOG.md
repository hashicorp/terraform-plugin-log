# 0.2.0 (Unreleased)

BREAKING CHANGES:
* Provider log outputs now default to being named "provider" unless another name is provided. ([#9](https://github.com/hashicorp/terraform-plugin-log/issues/9))
* The `tflog` package has been moved to `github.com/hashicorp/terraform-plugin-log/tflog` to make it work better with goimports and other tooling. ([#7](https://github.com/hashicorp/terraform-plugin-log/issues/7))
* The `tfsdklog` package has been moved to `github.com/hashicorp/terraform-plugin-log/tfsdklog` to make it work better with goimports and other tooling. ([#7](https://github.com/hashicorp/terraform-plugin-log/issues/7))
* With the release of Go 1.17, Go 1.16 is now the lowest supported version of Go to use with terraform-plugin-log. ([#8](https://github.com/hashicorp/terraform-plugin-log/issues/8))
* `tflog.New` has been moved to `tfsdklog.NewRootProviderLogger`. Provider developers should not need this functionality. ([#10](https://github.com/hashicorp/terraform-plugin-log/issues/10))
* `tflog.Option` has been moved to an internal package. Consumers can no longer reference the type. ([#10](https://github.com/hashicorp/terraform-plugin-log/issues/10))
* `tflog.WithLogName` has been removed. SDK developers should use `tfsdklog.WithLogName`. Provider developers should not need this functionality. ([#10](https://github.com/hashicorp/terraform-plugin-log/issues/10))
* `tflog.WithStderrFromInit` has been removed. SDK developers should use `tfsdklog.WithStderrFromInit`. Provider developers should not need this functionality. ([#10](https://github.com/hashicorp/terraform-plugin-log/issues/10))
* `tfsdklog.New` has been renamed `tfsdklog.NewRootLogger`. ([#10](https://github.com/hashicorp/terraform-plugin-log/issues/10))
* `tfsdklog.Option` has been moved to an internal package. Consumers can no longer reference the type. ([#10](https://github.com/hashicorp/terraform-plugin-log/issues/10))

FEATURES:
* Added a `tfsdklog.RegisterTestSink` function that will turn on a logging sink appropriate for use in testing. When a logging sink is activated, all downstream provider and SDK loggers (i.e., all loggers created after `tfsdklog.RegisterTestSink` is called) will be sub-loggers of the sink. The sink respects the `TF_LOG`, `TF_LOG_PATH`, `TF_ACC_LOG_PATH`, and `TF_LOG_PATH_MASK` environment variables, by default discards logs, and when only `TF_LOG` is set writes to stderr. It is meant to replicate Terraform's log aggregation and filtering capabilities for test frameworks. ([#9](https://github.com/hashicorp/terraform-plugin-log/issues/9))

# 0.1.0 (June 24, 2021)

FEATURES:
* Build out the beginnings of the module, allowing for providers, SDKs, and their subsystems to use an opinionated interface to log data about their execution. ([#2](https://github.com/hashicorp/terraform-plugin-log/issues/2))
