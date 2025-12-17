package middleware

import (
	"bytes"
	"net/http"
	"strings"
)

// liveReloadScript is the client-side script that connects to the WebSocket server.
const liveReloadScript = `
<!-- Goku Framework Live-Reload Script -->
<script>
    (function() {
        let socket = new WebSocket("ws://" + window.location.host + "/goku-ws");
        socket.onmessage = function(event) {
            if (event.data === 'reload') {
                console.log("Goku received reload signal. Refreshing...");
                window.location.reload();
            }
        };
        socket.onclose = function(event) {
            console.log("Goku live-reload socket closed. Manual refresh may be required.");
        };
        socket.onerror = function(error) {
            console.error("Goku live-reload socket error: ", error);
        };
    })();
</script>
`

// responseWriterInterceptor is a wrapper around http.ResponseWriter to capture the response.
type responseWriterInterceptor struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriterInterceptor) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

// LiveReload is a middleware that injects a script for live reloading.
func LiveReload(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Intercept the response
		interceptor := &responseWriterInterceptor{
			ResponseWriter: w,
			body:           &bytes.Buffer{},
		}

		next.ServeHTTP(interceptor, r)

		// Check if the response is HTML
		contentType := w.Header().Get("Content-Type")
		isHTML := strings.Contains(contentType, "text/html")

		// Get the original body
		originalBody := interceptor.body.Bytes()

		if isHTML {
			// Inject the script just before the closing </body> tag
			injectedBody := bytes.Replace(originalBody, []byte("</body>"), []byte(liveReloadScript+"</body>"), 1)
			
			// If body tag was not found, just append the script.
			// This is a fallback for malformed or partial HTML.
			if len(injectedBody) == len(originalBody) {
				injectedBody = append(originalBody, []byte(liveReloadScript)...)
			}

			// Write the modified body to the original response writer
			w.Write(injectedBody)
		} else {
			// If not HTML, write the original body without modification
			w.Write(originalBody)
		}
	})
}
