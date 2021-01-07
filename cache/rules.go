package cache

// these MIME types are cached, and nothing else is
var allowedContentTypes = []string{
	"text/html",
	"text/css",
	"application/javascript",
	"text/javascript", // for compatibility with bad servers and websites

	"image/png",
	"image/jpeg",
	"image/jpg",
	"image/webp",
	"image/svg+xml",
	"image/x-icon", // an unofficial type primarily used by .ico files, which almost all websites use (favicon.ico)

	"application/pdf",
	"font/ttf",
	"font/woff",
	"font/woff2",
	"application/font-woff2",
	"font/otf",
}
