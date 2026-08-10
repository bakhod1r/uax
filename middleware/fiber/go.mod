module github.com/bakhod1r/uax/middleware/fiber

go 1.26.2

require (
	github.com/bakhod1r/uax v0.1.0
	github.com/gofiber/fiber/v2 v2.52.0
)

require (
	github.com/andybalholm/brotli v1.0.5 // indirect
	github.com/google/uuid v1.5.0 // indirect
	github.com/klauspost/compress v1.17.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.51.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
)

// Local development builds against the tree, not the tag. Consumers ignore a
// replace in a dependency, so the require above is what they actually resolve
// and it must name a version that exists.
replace github.com/bakhod1r/uax => ../..
