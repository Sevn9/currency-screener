package healthz

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type HealthController struct {
	logger *zap.Logger
}

func NewHealthController(logger *zap.Logger) *HealthController {
	return &HealthController{logger: logger}
}

func (h *HealthController) Healthz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]string{"status": "ok"}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("Error writing simple healthz response:", zap.Error(err))
	}
}

func (h *HealthController) Index(w http.ResponseWriter, r *http.Request) {
	const html = `<html>
		<h1>Currency service status</h1>
		<p>Hostname: %s; Version: %s</p>
	</html>`
	fmt.Fprintf(w, html, "Currency", "1")
}
