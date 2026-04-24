package common

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const FrontendUnavailablePageHTML = `<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Frontend unavailable</title>
    <style>
      body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif; margin: 0; background: #0b1020; color: #e5e7eb; }
      main { max-width: 720px; margin: 0 auto; padding: 64px 24px; }
      code { background: rgba(255,255,255,0.08); padding: 2px 6px; border-radius: 6px; }
      .card { background: rgba(255,255,255,0.05); border: 1px solid rgba(255,255,255,0.1); border-radius: 16px; padding: 24px; }
      a { color: #7dd3fc; }
    </style>
  </head>
  <body>
    <main>
      <div class="card">
        <h1>Frontend is not available on this backend node</h1>
        <p>Build and mount <code>web/dist</code>, or set <code>FRONTEND_BASE_URL</code> to your deployed frontend origin.</p>
      </div>
    </main>
  </body>
</html>`

func ResolveFrontendDistDir() string {
	candidates := []string{}
	if envPath := strings.TrimSpace(os.Getenv("FRONTEND_DIST_DIR")); envPath != "" {
		candidates = append(candidates, envPath)
	}
	candidates = append(candidates, "web/dist")
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		candidates = append(candidates,
			filepath.Join(exeDir, "web", "dist"),
			filepath.Join(exeDir, "dist"),
		)
	}

	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}
		indexPath := filepath.Join(candidate, "index.html")
		if stat, err := os.Stat(indexPath); err == nil && !stat.IsDir() {
			return candidate
		}
	}
	SysLog(fmt.Sprintf("frontend dist not found; checked %v", candidates))
	return ""
}
