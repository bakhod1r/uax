module github.com/bakhod1r/uax/middleware/fiber

go 1.26.2

require (
	github.com/bakhod1r/uax v0.1.0
	github.com/gofiber/fiber/v2 v2.52.14
)

require (
	github.com/andybalholm/brotli v1.2.2 // indirect
	github.com/clipperhouse/uax29/v2 v2.7.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/klauspost/compress v1.19.2 // indirect
	github.com/mattn/go-colorable v0.1.15 // indirect
	github.com/mattn/go-isatty v0.0.24 // indirect
	github.com/mattn/go-runewidth v0.0.27 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.73.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
)

// Local development builds against the tree, not the tag. Consumers ignore a
// replace in a dependency, so the require above is what they actually resolve
// and it must name a version that exists.
replace github.com/bakhod1r/uax => ../..
