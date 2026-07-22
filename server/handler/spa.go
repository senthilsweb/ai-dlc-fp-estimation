package handler

import (
	"io"
	"io/fs"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// SPAHandler serves an embedded filesystem with SPA fallback.
// If the requested file exists, it is served directly.
// Otherwise, index.html is served (for client-side routing).
func SPAHandler(fsys fs.FS, prefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Strip the route prefix to get the file path within the embedded FS
		path := strings.TrimPrefix(c.Request.URL.Path, prefix)
		path = strings.TrimPrefix(path, "/")
		if path == "" {
			path = "index.html"
		}

		// Try to open the requested file
		f, err := fsys.Open(path)
		if err != nil {
			// File not found — serve index.html (SPA fallback)
			path = "index.html"
			f, err = fsys.Open(path)
			if err != nil {
				c.Status(http.StatusNotFound)
				return
			}
		}
		defer f.Close()

		// Check if it's a directory — if so, try index.html inside it
		stat, err := f.Stat()
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		if stat.IsDir() {
			f.Close()
			path = path + "/index.html"
			f, err = fsys.Open(path)
			if err != nil {
				// Directory with no index — SPA fallback
				path = "index.html"
				f, err = fsys.Open(path)
				if err != nil {
					c.Status(http.StatusNotFound)
					return
				}
			}
			defer f.Close()
		}

		// Detect content type from extension
		contentType := mime.TypeByExtension(filepath.Ext(path))
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		c.Header("Content-Type", contentType)
		c.Status(http.StatusOK)
		io.Copy(c.Writer, f.(io.Reader))
	}
}
