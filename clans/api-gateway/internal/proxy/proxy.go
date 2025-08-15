package proxy

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/AlexGuo43/clans/api-gateway/config"
)

type Gateway struct {
	config *config.Config
	client *http.Client
}

func NewGateway(cfg *config.Config) *Gateway {
	return &Gateway{
		config: cfg,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (g *Gateway) RouteRequest(w http.ResponseWriter, r *http.Request) {
	service := g.getServiceForPath(r.URL.Path)
	if service == nil {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	targetURL := g.buildTargetURL(service, r)
	g.forwardRequest(w, r, targetURL)
}

func (g *Gateway) getServiceForPath(path string) *config.ServiceConfig {
	switch {
	case strings.HasPrefix(path, "/api/auth/") || strings.HasPrefix(path, "/api/users/"):
		return &g.config.UserService
	case strings.HasPrefix(path, "/api/posts"):
		return &g.config.PostService
	case strings.HasPrefix(path, "/api/comments"):
		return &g.config.CommentService
	case strings.HasPrefix(path, "/api/clans"):
		return &g.config.ClanService
	default:
		return nil
	}
}

func (g *Gateway) buildTargetURL(service *config.ServiceConfig, r *http.Request) string {
	targetPath := r.URL.Path
	
	switch service.Name {
	case "user-service":
		if strings.HasPrefix(targetPath, "/api/auth/") {
			targetPath = strings.TrimPrefix(targetPath, "/api/auth")
		} else if strings.HasPrefix(targetPath, "/api/users/") {
			targetPath = strings.TrimPrefix(targetPath, "/api/users")
		}
	case "post-service":
		// Keep the full path for post-service as it expects /api/posts
		// targetPath = strings.TrimPrefix(targetPath, "/api")
	case "comment-service":
		// Keep the full path for comment-service as it expects /api/comments
		// targetPath = strings.TrimPrefix(targetPath, "/api")
	case "clan-service":
		targetPath = strings.TrimPrefix(targetPath, "/api")
	}

	targetURL := fmt.Sprintf("%s%s", service.URL, targetPath)
	if r.URL.RawQuery != "" {
		targetURL += "?" + r.URL.RawQuery
	}
	
	return targetURL
}

func (g *Gateway) forwardRequest(w http.ResponseWriter, r *http.Request, targetURL string) {
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		http.Error(w, "Invalid target URL", http.StatusInternalServerError)
		return
	}

	// Read and copy the body to avoid consumption issues
	var bodyBytes []byte
	if r.Body != nil {
		bodyBytes, err = io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
	}

	proxyReq, err := http.NewRequest(r.Method, parsedURL.String(), bytes.NewReader(bodyBytes))
	if err != nil {
		http.Error(w, "Failed to create proxy request", http.StatusInternalServerError)
		return
	}

	for key, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(key, value)
		}
	}

	proxyReq.Header.Set("X-Forwarded-For", r.RemoteAddr)
	proxyReq.Header.Set("X-Forwarded-Proto", "http")

	resp, err := g.client.Do(proxyReq)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (g *Gateway) HealthCheck(w http.ResponseWriter, r *http.Request) {
	status := make(map[string]string)
	
	for _, service := range g.config.Services {
		healthURL := fmt.Sprintf("%s/health", service.URL)
		resp, err := g.client.Get(healthURL)
		
		if err != nil || resp.StatusCode != http.StatusOK {
			status[service.Name] = "unhealthy"
		} else {
			status[service.Name] = "healthy"
		}
		
		if resp != nil {
			resp.Body.Close()
		}
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"services": %v}`, status)
}