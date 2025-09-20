package lang

import "embed"

// Files содержит встроенные языковые ресурсы.
//
//go:embed *.txt
var Files embed.FS
