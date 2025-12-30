package logintrap

import (
	"context"
	"fmt"
	"time"

	"github.com/Kartikey2011yadav/voidsink/internal/telemetry"
	"github.com/Kartikey2011yadav/voidsink/pkg/notifier"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

const loginPageHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Administration Login</title>
    <style>
        body { font-family: sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; background-color: #f0f2f5; }
        .login-container { background: white; padding: 2rem; border-radius: 8px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); width: 300px; }
        h2 { text-align: center; color: #333; }
        input { width: 100%; padding: 10px; margin: 10px 0; border: 1px solid #ddd; border-radius: 4px; box-sizing: border-box; }
        button { width: 100%; padding: 10px; background-color: #007bff; color: white; border: none; border-radius: 4px; cursor: pointer; }
        button:hover { background-color: #0056b3; }
        .error { color: red; text-align: center; margin-bottom: 10px; font-size: 0.9em; }
    </style>
</head>
<body>
    <div class="login-container">
        <h2>Admin Panel</h2>
        <form method="POST">
            <input type="text" name="username" placeholder="Username" required>
            <input type="password" name="password" placeholder="Password" required>
            <button type="submit">Login</button>
        </form>
    </div>
</body>
</html>`

// LoginTrap implements the Trap interface for a fake login page.
type LoginTrap struct {
	addr       string
	serverName string
	server     *fasthttp.Server
	notifier   *notifier.Notifier
}

// New creates a new instance of LoginTrap.
func New(addr, serverName string, n *notifier.Notifier) *LoginTrap {
	return &LoginTrap{
		addr:       addr,
		serverName: serverName,
		notifier:   n,
	}
}

// Start starts the HTTP server.
func (t *LoginTrap) Start(ctx context.Context) error {
	t.server = &fasthttp.Server{
		Handler:      t.requestHandler,
		Name:         t.serverName,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	log.Info().Str("address", t.addr).Msg("Starting Login Trap")

	errChan := make(chan error, 1)
	go func() {
		errChan <- t.server.ListenAndServe(t.addr)
	}()

	select {
	case <-ctx.Done():
		return t.Shutdown(context.Background())
	case err := <-errChan:
		return err
	}
}

// Shutdown gracefully shuts down the server.
func (t *LoginTrap) Shutdown(ctx context.Context) error {
	log.Info().Msg("Shutting down Login Trap")
	if t.server != nil {
		return t.server.Shutdown()
	}
	return nil
}

func (t *LoginTrap) requestHandler(ctx *fasthttp.RequestCtx) {
	path := string(ctx.Path())
	remoteIP := ctx.RemoteAddr().String()
	userAgent := string(ctx.UserAgent())
	method := string(ctx.Method())

	log.Info().Str("path", path).Str("method", method).Str("remote_addr", remoteIP).Msg("Login Trap hit")
	telemetry.TrapsTriggered.WithLabelValues("login_trap", path).Inc()

	if method == "POST" {
		// Capture Credentials
		username := string(ctx.FormValue("username"))
		password := string(ctx.FormValue("password"))

		if username != "" || password != "" {
			log.Warn().
				Str("username", username).
				Str("password", password).
				Str("remote_addr", remoteIP).
				Msg("Credentials Captured")

			telemetry.CredentialsCaptured.Inc()

			if t.notifier != nil {
				// High priority alert
				// We pass the credentials in the "trapName" field so they appear in the alert
				t.notifier.SendAlert(fmt.Sprintf("LoginTrap (User: %s, Pass: %s)", username, password), remoteIP, userAgent)
			}
		}

		// Return failure to encourage more tries
		ctx.SetContentType("text/html")
		// Inject error message into the HTML
		// Simple string replacement for this demo
		response := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Administration Login</title>
    <style>
        body { font-family: sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; background-color: #f0f2f5; }
        .login-container { background: white; padding: 2rem; border-radius: 8px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); width: 300px; }
        h2 { text-align: center; color: #333; }
        input { width: 100%; padding: 10px; margin: 10px 0; border: 1px solid #ddd; border-radius: 4px; box-sizing: border-box; }
        button { width: 100%; padding: 10px; background-color: #007bff; color: white; border: none; border-radius: 4px; cursor: pointer; }
        button:hover { background-color: #0056b3; }
        .error { color: red; text-align: center; margin-bottom: 10px; font-size: 0.9em; }
    </style>
</head>
<body>
    <div class="login-container">
        <h2>Admin Panel</h2>
        <div class="error">Invalid username or password</div>
        <form method="POST">
            <input type="text" name="username" placeholder="Username" required>
            <input type="password" name="password" placeholder="Password" required>
            <button type="submit">Login</button>
        </form>
    </div>
</body>
</html>`
		ctx.WriteString(response)
		return
	}

	// GET Request - Serve Login Page
	ctx.SetContentType("text/html")
	ctx.WriteString(loginPageHTML)
}
