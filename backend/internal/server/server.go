package server

import (
	"context"
	"foodrecipe/internal/auth"
	"foodrecipe/internal/middleware"
	"foodrecipe/internal/recipe"
	"foodrecipe/internal/utils"
	"net/http"
)

	type Server struct {
		router        *http.ServeMux
		cfg          *utils.Config
		server       *http.Server
		authHandler  *auth.Handler
		recipeHandler *recipe.RecipeHandler
	}

	func NewServer(cfg *utils.Config, authHandler *auth.Handler, recipeHandler *recipe.RecipeHandler) *Server {
		return &Server{
			router:        http.NewServeMux(),
			cfg:           cfg,
			authHandler:   authHandler,
			recipeHandler: recipeHandler,
		}
	}

	func (s *Server) SetupRoutes() {
		// Apply common middleware to all routes
		commonMiddleware := func(h http.Handler) http.Handler {
			return middleware.RecoveryMiddleware(
				middleware.LoggingMiddleware(h),
			)
		}

		// Public routes
		s.router.HandleFunc("GET /api/health", s.handleHealth)
		s.router.HandleFunc("POST /api/auth/register", s.authHandler.HandleRegister)
		s.router.HandleFunc("POST /api/auth/login", s.authHandler.HandleLogin)
		s.router.HandleFunc("POST /api/auth/refresh", s.authHandler.HandleRefreshToken)

		// Protected routes
		protected := http.NewServeMux()
		protected.HandleFunc("POST /api/recipes", s.recipeHandler.CreateRecipe)
		protected.HandleFunc("GET /api/recipes", s.handleGetRecipes)

		// Apply auth middleware to protected routes
		protectedWithAuth := middleware.AuthMiddleware(s.cfg.JWTSecret)(protected)

		// Mount all routes
		s.router.Handle("/api/", commonMiddleware(protectedWithAuth)) // Protected routes
		s.router.Handle("/", commonMiddleware(s.router))             // Public routes
	}

	func (s *Server) Start() error {
		s.server = &http.Server{
			Addr:    ":" + s.cfg.HTTPPort,
			Handler: s.router, // Use the router directly now
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

	func (s *Server) handleGetRecipes(w http.ResponseWriter, r *http.Request) {
		// Delegate to recipe handler if you have one
		// s.recipeHandler.GetRecipes(w, r)
		w.Header().Set("Content-Type", "application/json")
		utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Protected recipes endpoint"})
	}