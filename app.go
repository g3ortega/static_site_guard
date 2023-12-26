package main

import (
	"log"
	"os"

	"github.com/g3ortega/static_site_guard/authentication"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/joho/godotenv"
	"github.com/markbates/goth/providers/github"
	"github.com/shareed2k/goth_fiber"

	"github.com/markbates/goth"

	"github.com/gofiber/storage/postgres/v3"

	"github.com/gofiber/template/html"
)

func main() {
	if os.Getenv("SSG_ENVIRONMENT") == "development" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	goth.UseProviders(
		github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), os.Getenv("GITHUB_CALLBACK")),
	)

	storage := postgres.New(postgres.Config{
		ConnectionURI: os.Getenv("POSTGRES_DATABASE_URL"),
		Table:         "sessions",
		Reset:         true,
	})

	config := session.Config{
		Storage:        storage,
		CookieSameSite: "Lax",
	}

	store := session.New(config)
	goth_fiber.SessionStore = store

	app.Use(logger.New())

	app.Use(func(ctx *fiber.Ctx) error {
		return authentication.SessionHandler(ctx, store)
	})

	app.Static("/", "./content/public", fiber.Static{
		Browse: true,
	})

	app.Get("/login/:provider", goth_fiber.BeginAuthHandler)

	app.Get("/auth/callback/:provider", func(ctx *fiber.Ctx) error {
		return authentication.Callback(ctx, store)
	})

	app.Get("/login", func(ctx *fiber.Ctx) error {
		return ctx.Render("index", fiber.Map{})
	})

	app.Get("/not_authorized", func(ctx *fiber.Ctx) error {
		return ctx.Render("not_authorized", fiber.Map{})
	})

	app.Get("/logout", func(ctx *fiber.Ctx) error {
		return authentication.Logout(ctx, store)
	})

	if err := app.Listen(":8088"); err != nil {
		log.Fatal(err)
	}
}
