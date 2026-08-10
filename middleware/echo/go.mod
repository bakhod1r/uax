module github.com/bakhod1r/uax/middleware/echo

go 1.26.2

require (
	github.com/bakhod1r/uax v0.1.0
	github.com/labstack/echo/v4 v4.15.4
)

require (
	github.com/labstack/gommon v0.5.0 // indirect
	github.com/mattn/go-colorable v0.1.15 // indirect
	github.com/mattn/go-isatty v0.0.24 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	golang.org/x/crypto v0.54.0 // indirect
	golang.org/x/net v0.57.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.40.0 // indirect
)

// Local development builds against the tree, not the tag. Consumers ignore a
// replace in a dependency, so the require above is what they actually resolve
// and it must name a version that exists.
replace github.com/bakhod1r/uax => ../..
