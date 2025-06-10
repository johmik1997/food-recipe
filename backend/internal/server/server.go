package server

import (
	"context"
	"foodrecipe/internal/auth"
	"foodrecipe/internal/middleware"
	"foodrecipe/internal/recipe"
	"foodrecipe/internal/uploadRecipeimage"
	"foodrecipe/internal/uploadimage"
	"foodrecipe/internal/utils"
	"net/http"
	"foodrecipe/handlers"
)

type Server struct {
	router        *http.ServeMux
	cfg          *utils.Config
	server       *http.Server
	authHandler  *auth.Handler
	recipeHandler *recipe.Handler
	imageHandler *uploadimage.Handler
	recipeImageHandler *uploadRecipeimage.Handler
	//paymentHandler    *handlers.PaymentHandler
	
}

func NewServer(
	cfg *utils.Config,
	authHandler *auth.Handler,
	recipeHandler *recipe.Handler,
	imageHandler *uploadimage.Handler,
	recipeImageHandler  *uploadRecipeimage.Handler,
	//paymentHandler *handlers.Paymenthandler,
) *Server {
	return &Server{
		router:        http.NewServeMux(),
		cfg:          cfg,
		authHandler:  authHandler,
		recipeHandler: recipeHandler,
		imageHandler: imageHandler,
		recipeImageHandler: recipeImageHandler,
		//paymentHandler:    paymentHandler,
	}
}

func (s *Server) SetupRoutes() {
	// Common middleware stack
	commonMiddleware := func(h http.Handler) http.Handler {
		return middleware.RecoveryMiddleware(
			middleware.LoggingMiddleware(
				middleware.CorsMiddleware(h),
			),
		)
	}

	// Public routes
	s.router.HandleFunc("GET /api/health", s.handleHealth)
	s.router.HandleFunc("POST /api/auth/register", s.authHandler.HandleRegister)
	s.router.HandleFunc("POST /api/auth/login", s.authHandler.HandleLogin)
	s.router.HandleFunc("POST /api/auth/refresh", s.authHandler.HandleRefreshToken)

	// Protected routes
	protected := http.NewServeMux()
	// protected.HandleFunc("POST /api/recipes", s.recipeHandler.CreateRecipe)
	protected.HandleFunc("/api/recipes", s.recipeHandler.CreateRecipe)
	//protected.HandleFunc("GET /api/recipes/{id}", s.recipeHandler.GetRecipe)
	protected.HandleFunc("/api/recipes/upload-image", s.recipeImageHandler.UploadRecipeImageHandler)
	protected.HandleFunc("/api/upload-profile-image", s.imageHandler.UploadProfileImage)
  s.router.Handle("POST /api/payments", commonMiddleware(http.HandlerFunc(handlers.HandlePayment)))	// Apply auth middleware to protected routes
	protectedWithAuth := middleware.AuthMiddleware(s.cfg.JWTSecret)(protected)

	// Mount all routes with proper middleware
	s.router.Handle("/api/", commonMiddleware(protectedWithAuth))
	s.router.Handle("/", commonMiddleware(http.NotFoundHandler()))
}

func (s *Server) Start() error {
	s.server = &http.Server{
		Addr:    ":" + s.cfg.HTTPPort,
		Handler: s.router,
	}
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.server != nil {
		return s.server.Shutdown(ctx)
	}
	return nil
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "OK"})
}