package resource

import "embed"

//go:embed public
//go:embed templates
var FS embed.FS
