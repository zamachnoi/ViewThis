package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/zamachnoi/viewthis/handlers"
	"github.com/zamachnoi/viewthis/lib"
	auth "github.com/zamachnoi/viewthis/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	lib.InitDB()
	lib.InitRD()
	// CORS middleware to handle cross-origin requests
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://example.com"}, // Adjust this to your frontend's address
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	r := chi.NewRouter()
	r.Use(corsMiddleware.Handler)

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Prefix all routes with "/api"
	api := chi.NewRouter()

	// Define your API routes here
	api.Route("/users", func(r chi.Router) {
		r.Get("/{id}", handlers.GetUserByIDHandler) // Get user by ID
		r.Post("/", handlers.CreateUserHandler)     // Create user
	})

	api.Route("/queues", func(r chi.Router) {
		r.Get("/", handlers.GetAllQueuesHandler)                // Get all queues
		r.Post("/", handlers.CreateQueueHandler)                // Create a new queue
		r.Patch("/{id}", handlers.UpdateQueueHandler)           // Update queue by ID
		r.Delete("/{id}", handlers.DeleteQueueHandler)          // Delete queue by ID
		r.Patch("/{id}/clear", handlers.ClearQueueByIDHandler)  // Clear queue by ID
	})

	api.Route("/queues/{queueID}/submissions", func(r chi.Router) {
		r.Get("/", handlers.GetSubmissionsByQueueIDHandler)      // Get all submissions for a queue
		r.Post("/", handlers.CreateSubmissionHandler)            // Create a new submission
		r.Delete("/{id}", handlers.DeleteSubmissionByIDHandler)  // Delete submission by ID
		r.Patch("/{id}", handlers.UpdateSubmissionHandler)       // Update submission by ID
	})

	api.Route("/auth/discord", func(r chi.Router) {
		r.Get("/login", handlers.DiscordAuthLoginHandler)
		r.Get("/callback", handlers.DiscordAuthCallbackHandler)
	})

	// Apply auth middleware only to specific routes
	api.Group(func(r chi.Router) {
		r.Use(auth.JWTAuthMiddleware)
		
		r.Route("/delete", func(r chi.Router) {
			r.Delete("/submissions", handlers.DeleteAllSubmissionsHandler)
			r.Delete("/feedback", handlers.DeleteAllFeedbackHandler)
			r.Delete("/queues", handlers.DeleteAllQueuesHandler)
			r.Delete("/users", handlers.DeleteAllUsersHandler)
			r.Delete("/data", handlers.DeleteAllDataHandler)
		})
		
		r.Route("/test", func(r chi.Router) {
			r.Get("/", handlers.TestingHandler)
		})
	})

	// Mount the API router under "/api" prefix
	r.Mount("/api", api)

	log.Println("Server starting on port 3001...")
	if err := http.ListenAndServe(":3001", r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
