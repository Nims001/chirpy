package main

// - `package main` declares this file belongs to the `main` package - required for any Go program that's meant to be run directly (not imported as a library).
// - `import` brings in two standard library packages:
//     - `net/http` gives you everything for building an HTTP server (ServeMux, Server, FileServer, etc.)
//     - `log` gives you tools for printing log messages, optionally with a fatal exit.
import (
	"log"
	"net/http"
)

func main() {
	// This creates a new **ServeMux** - your router. Right now it's empty; no paths are registered yet. `:=` declares and assigns in one step (shorthand for `var sermux = ...`).
	sermux := http.NewServeMux()

	// 	This creates an `http.Server` struct:

	// - `Addr: ":8080"` tells it to listen on port 8080 on all network interfaces.
	// - `Handler: sermux` tells the server "whenever a request comes in, hand it off to `sermux` to figure out what to do."
	// - The `&` means you're taking a **pointer** to this struct - `s` holds the memory address of the server, not a copy of it. This matters because things like `ListenAndServe` need to be called on the actual server, not some duplicate copy.
	s := &http.Server{
		Addr:    ":8080",
		Handler: sermux,
	}

	// This line registers a new route with the ServeMux:
	// 	Notice: at this point, `sermux` has no routes yet, but that's fine because Go doesn't execute this code top-to-bottom in the sense of "checking" anything - it's just building objects in memory.

	// ```go
	// 	sermux.Handle("/app/",http.FileServer(http.Dir("./")))
	// ```

	// Breaking this down from the inside out:

	// - `http.Dir("./")` - treats the current directory as the root folder to serve files from.
	// - `http.FileServer(...)` - creates a handler that knows how to serve files/directories as HTTP responses (handling things like content-type detection, directory listings, etc.).
	// - `sermux.Handle("/app/", ...)` - registers that fileserver handler on the mux, so that any request path starting with `/app/` gets routed to it.
	sermux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("./"))))
	// we added the `http.StripPrefix` function to remove the `/app` prefix from the request path before passing it to the file server. This is necessary because the
	// file server expects paths relative to the root of the directory it's serving, and without stripping the prefix, it would look for files in a non-existent `/app` subdirectory.

	// This makes sure that when a request comes in for `/healthz`, it gets routed to `handlerfunction`. The `handlerfunction` is defined below and simply responds with a 200 OK and the text "OK".
	sermux.HandleFunc("/healthz", handlerfunction)

	// 	Note: this line runs _after_ `s` was created, but that's fine in Go - `s.Handler` holds a _reference_ to `sermux`, not a snapshot. So even though you registered the route after building the server struct, the server will still see it because it's looking at the same `sermux` object in memory.

	// ```go
	// 	err := s.ListenAndServe()
	// ```

	// This is the line that actually starts the server - it blocks (pauses execution here) and listens for incoming TCP connections on port 8080,
	// routing each one through `sermux`. It only returns when the server stops, usually due to an error.

	// ```go
	// 	if err != nil{
	// 		log.Printf("Error occured")
	// 		return
	// 	}
	// 	return
	// }
	// ```

	// If `ListenAndServe` returns an error (e.g., the port is already in use), you log a message and exit `main`.
	// The final `return` at the bottom is actually unreachable in the sense that `ListenAndServe` almost never returns `nil` - it blocks forever until something goes wrong.

	err := s.ListenAndServe()
	if err != nil {
		log.Printf("Error occured")
		return
	}
	// there was a return here that never got executed because ListenAndServe blocks until an error occurs, so the program would exit after logging the error.
}

// http.ResponseWriter is an interface that allows you to construct an HTTP response. It has methods for setting headers, writing the status code, and writing the body of the response.
// http.Request is a struct that represents an incoming HTTP request. It contains information like the method (GET, POST, etc.), URL, headers, and body of the request.
// http.Request is used as a pointer because it can be large and you want to avoid copying it around. Passing a pointer allows the handler to read from the request
// without making a full copy of it, which is more efficient.
func handlerfunction(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}
