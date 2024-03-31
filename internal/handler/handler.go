package handler

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/syarifid/bankx/internal/cfg"
	"github.com/syarifid/bankx/internal/service"
)

// var (
// 	requestsTotal = prometheus.NewCounterVec(
// 		prometheus.CounterOpts{
// 			Name: "http_requests_total",
// 			Help: "Total number of HTTP requests.",
// 		},
// 		[]string{"method", "path", "status"},
// 	)
// 	requestDuration = prometheus.NewHistogramVec(
// 		prometheus.HistogramOpts{
// 			Name:    "http_request_duration_seconds",
// 			Help:    "Histogram of request duration in seconds.",
// 			Buckets: prometheus.DefBuckets,
// 		},
// 		[]string{"method", "path", "status"},
// 	)
// )

// Define the histogram metric.

var (
	httpRequestProm = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_histogram",
		Help:    "Histogram of the http request duration.",
		Buckets: prometheus.LinearBuckets(1, 1, 10),
	}, []string{"path", "method", "status"})
)

// type responseWriter struct {
// 	http.ResponseWriter
// 	statusCode int
// }

// func newResponseWriter(w http.ResponseWriter) *responseWriter {
// 	return &responseWriter{w, http.StatusOK}
// }

// func (rw *responseWriter) WriteHeader(code int) {
// 	rw.statusCode = code
// 	rw.ResponseWriter.WriteHeader(code)
// }

// func (rw *responseWriter) Status() int {
// 	return rw.statusCode
// }

type Handler struct {
	router  *chi.Mux
	service *service.Service
	cfg     *cfg.Cfg
}

func NewHandler(router *chi.Mux, service *service.Service, cfg *cfg.Cfg) *Handler {
	handler := &Handler{router, service, cfg}
	handler.registRoute()

	return handler
}

func (h *Handler) registRoute() {
	// prometheus.MustRegister(requestsTotal, requestDuration)

	r := h.router
	var tokenAuth *jwtauth.JWTAuth = jwtauth.New("HS256", []byte(h.cfg.JWTSecret), nil, jwt.WithAcceptableSkew(30*time.Second))

	userH := newUserHandler(h.service.User)
	fileH := newFileHandler(h.cfg)
	transactionH := newTransactionHandler(h.service.Transaction)

	// r.Use(middleware.RedirectSlashes)
	// r.Use(prometheusMiddleware)

	// r.Get("/metrics", func(h http.Handler) http.HandlerFunc {
	// 	return func(w http.ResponseWriter, r *http.Request) {
	// 		h.ServeHTTP(w, r)
	// 	}
	// }(promhttp.Handler()))

	c := chi.NewRouter()
	c.Use(ChiPrometheusMiddleware)
	c.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(w, r)
	})

	// GET /healthz -> 200 OK
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"service": "ok"}`))
	})

	r.Post("/v1/user/register", userH.Register)
	r.Post("/v1/user/login", userH.Login)

	// protected route
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))

		r.Patch("/v1/user", userH.UpdateAccount)

		r.Post("/v1/balance", transactionH.AddBalance)
		r.Get("/v1/balance", transactionH.GetBalance)
		r.Get("/v1/balance/history", transactionH.GetBalanceHistory)

		r.Post("/v1/transaction", transactionH.AddTransaction)

		r.Post("/v1/image", fileH.Upload)
	})
}

// func prometheusMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		startTime := time.Now()
// 		rw := newResponseWriter(w)
// 		defer func() {
// 			status := rw.Status()
// 			duration := time.Since(startTime).Seconds()
// 			requestsTotal.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(status)).Inc()
// 			requestDuration.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(status)).Observe(duration)
// 		}()
// 		next.ServeHTTP(rw, r)
// 	})
// }

func ChiPrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r) // Process request

		status := http.StatusOK // Assuming status OK, customize as needed
		httpRequestProm.WithLabelValues(r.URL.Path, r.Method, http.StatusText(status)).Observe(float64(time.Since(start).Milliseconds()))
	})
}
