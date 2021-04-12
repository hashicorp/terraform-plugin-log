package tflog

import "github.com/hashicorp/go-hclog"

// this is just a placeholder file to make CircleCI happy with an empty
// repository.
//
// We created a package so Go vet will be happy.
//
// We're using hclog so we have an import, so our cache behavior will stop
// breaking on the lack of a go.sum file.
var _ = hclog.NoLevel
