module my-plugin

go 1.12

require (
	github.com/bogus/my-abstraction v0.0.0
	github.com/hashicorp/go-plugin v1.0.0
	github.com/pkg/errors v0.8.1
)

replace github.com/bogus/my-abstraction => ../my-abstraction
