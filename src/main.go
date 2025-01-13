package main

import (
	"encoding/json"
	"goscraper/src/handlers"
	"goscraper/src/utils"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	app := fiber.New(fiber.Config{
		Prefork:      true,
		ServerHeader: "GoScraper",
		AppName:      "GoScraper v1.0",
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return utils.HandleError(c, err)
		},
	})

	// Middleware stack
	app.Use(recover.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))
	app.Use(etag.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://class-pro.vercel.app,http://localhost:3024",
		AllowMethods:     "GET,POST,DELETE",
		AllowHeaders:     "Origin,Content-Type,Accept,X-CSRF-Token",
		ExposeHeaders:    "Content-Length",
		AllowCredentials: true,
		MaxAge:           int((12 * time.Hour).Seconds()),
	}))

	app.Use(limiter.New(limiter.Config{
		Max:        25,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			token := c.Get("X-CSRF-Token")
			if token != "" {
				return utils.Encode(token)
			}
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "🔨 SHUT UP! Rate limit exceeded. Please try again later.",
			})
		},
		SkipFailedRequests: false,
		LimiterMiddleware:  limiter.SlidingWindow{},
	}))

	// CSRF middleware
	app.Use(func(c *fiber.Ctx) error {
		if c.Path() == "/login" {
			return c.Next()
		}

		token := c.Get("X-CSRF-Token")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing X-CSRF-Token header",
			})
		}
		return c.Next()
	})

	// Add cache middleware configuration
	cacheConfig := cache.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Method() != "GET"
		},
		Expiration: 2 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Path() + "_" + c.Get("X-CSRF-Token")
		},
	}

	// Route group for authenticated endpoints
	api := app.Group("/", func(c *fiber.Ctx) error {
		token := c.Get("X-CSRF-Token")
		if token == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Missing X-CSRF-Token header",
			})
		}
		return c.Next()
	})

	// Routes
	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Hello, World!"})
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		var creds struct {
			Username string `json:"account"`
			Password string `json:"password"`
		}

		if err := c.BodyParser(&creds); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid JSON body",
			})
		}

		if creds.Username == "" || creds.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Missing account or password",
			})
		}

		lf := &handlers.LoginFetcher{}
		session, err := lf.CampusLogin(creds.Username, creds.Password)
		if err != nil {
			return err
		}

		return c.JSON(session)
	})

	api.Delete("/logout", func(c *fiber.Ctx) error {
		lf := &handlers.LoginFetcher{}
		session, err := lf.Logout(c.Get("X-CSRF-Token"))
		if err != nil {
			return err
		}
		return c.JSON(session)
	})

	// Apply cache middleware to GET routes
	api.Get("/attendance", cache.New(cacheConfig), func(c *fiber.Ctx) error {
		attendance, err := handlers.GetAttendance(c.Get("X-CSRF-Token"))
		if err != nil {
			return err
		}
		return c.JSON(attendance)
	})

	api.Get("/marks", cache.New(cacheConfig), func(c *fiber.Ctx) error {
		marks, err := handlers.GetMarks(c.Get("X-CSRF-Token"))
		if err != nil {
			return err
		}
		return c.JSON(marks)
	})

	api.Get("/courses", cache.New(cacheConfig), func(c *fiber.Ctx) error {
		courses, err := handlers.GetCourses(c.Get("X-CSRF-Token"))
		if err != nil {
			return err
		}
		return c.JSON(courses)
	})

	api.Get("/user", cache.New(cacheConfig), func(c *fiber.Ctx) error {
		user, err := handlers.GetUser(c.Get("X-CSRF-Token"))
		if err != nil {
			return err
		}
		return c.JSON(user)
	})

	api.Get("/calendar", cache.New(cacheConfig), func(c *fiber.Ctx) error {
		cal, err := handlers.GetCalendar(c.Get("X-CSRF-Token"))
		if err != nil {
			return err
		}
		return c.JSON(cal)
	})

	api.Get("/timetable", cache.New(cacheConfig), func(c *fiber.Ctx) error {
		tt, err := handlers.GetTimetable(c.Get("X-CSRF-Token"))
		if err != nil {
			return err
		}
		return c.JSON(tt)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %s...", port)
	if err := app.Listen("0.0.0.0:" + port); err != nil {
		log.Printf("Server error: %+v", err)
		log.Fatal(err)
	}
}
