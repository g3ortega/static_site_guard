package authentication

import (
	"log"
	"os"
	"strings"

	"github.com/g3ortega/static_site_guard/encrypt_tools"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/shareed2k/goth_fiber"
)

const (
	loginPath         = "/login"
	loginGithubPath   = "/login/github"
	notAuthorizedPath = "/not_authorized"
	authCallbackPath  = "/auth/callback/github"
	logoutPath        = "/logout"
)

func Callback(ctx *fiber.Ctx, store *session.Store) error {
	user, _ := goth_fiber.CompleteUserAuth(ctx)
	sess, err := store.Get(ctx)
	if err != nil {
		panic(err)
	}

	if err != nil {
		return ctx.Status(400).JSON(map[string]string{"message": err.Error()})
	} else {
		encryptedUsername := encrypt_tools.Encrypt(os.Getenv("SECRET_KEY"), user.NickName)
		sess.Set("userName", encryptedUsername)
		err := sess.Save()

		if err != nil {
			return ctx.Status(400).JSON(map[string]string{"message": err.Error()})
		}
	}

	if eligibleUser(user.NickName) {
		return ctx.Redirect("/")
	} else {
		return ctx.Redirect("/logout")
	}
}

func Logout(ctx *fiber.Ctx, store *session.Store) error {
	if err := goth_fiber.Logout(ctx); err != nil {
		log.Fatal(err)
	}

	sess, err := store.Get(ctx)
	if err != nil {
		panic(err)
	}

	userName := sess.Get("userName")

	if userName != nil {
		sess.Delete("userName")
		errSessionDestroy := sess.Destroy()
		if errSessionDestroy != nil {
			log.Println(errSessionDestroy)
		}

		errSessionSave := sess.Save()
		if errSessionSave != nil {
			log.Println(errSessionSave)
		}
	}

	return ctx.Redirect("/login")
}

func eligibleUser(userName string) bool {
	eligibleUsernames := os.Getenv("ELIGIBLE_USERNAMES")
	if eligibleUsernames == "" {
		log.Fatal("ELIGIBLE_USERNAMES is not set")
	}

	set := make(map[string]bool)
	for _, v := range strings.Split(eligibleUsernames, ",") {
		set[v] = true
	}

	return set[userName]
}

func SessionHandler(ctx *fiber.Ctx, store *session.Store) (err error) {
	sess, err := store.Get(ctx)
	if err != nil {
		panic(err)
	}
	encryptedUsername := sess.Get("userName")
	userName := ""

	if encryptedUsername != nil {
		userName = encrypt_tools.Decrypt(os.Getenv("SECRET_KEY"), encryptedUsername.(string))
	}

	switch ctx.Path() {
	case loginPath, loginGithubPath, notAuthorizedPath, authCallbackPath, logoutPath:
		return ctx.Next()
	default:
		if eligibleUser(userName) {
			return ctx.Next()
		} else {
			return ctx.Redirect(notAuthorizedPath)
		}
	}
}
