package handler

import (
	"3d-print-account/config"
	"3d-print-account/database"
	"3d-print-account/enum/sex"
	"3d-print-account/model"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/idtoken"
	"gorm.io/gorm"
)

// CheckPasswordHash compare password with hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func getUserByUsername(u string) (*model.User, error) {
	db := database.DB
	var user model.User
	if err := db.Table("user").Where("user_username = ?", u).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Login get user and password
func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	type UserData struct {
		ID       uuid.UUID `json:"id"`
		Username string    `json:"username"`
	}
	input := new(LoginInput)
	var ud UserData
	fmt.Print(c.Body())
	if err := c.BodyParser(&input); err != nil {
		// return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"code": c.Response().StatusCode(), "message": "Error on login request", "data": err})
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	username := input.Username
	password := input.Password

	user, err := getUserByUsername(username)

	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"code": c.Response().StatusCode(), "message": "Error on username", "data": err})
	}

	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"code": c.Response().StatusCode(), "message": "Invalid username or password", "data": nil})
	}

	if user != nil {
		ud = UserData{
			ID:       user.ID,
			Username: user.Username,
		}
	}

	if !CheckPasswordHash(password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid username or password", "data": nil})
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = ud.Username
	claims["id"] = ud.ID
	claims["exp"] = time.Now().Add(time.Second * 30).Unix()

	t, err := token.SignedString([]byte(config.Config("SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	err = logSession(t)

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t, "message": "Login successfully!", "status": "success", "code": c.Response().StatusCode()})
}

func LoginWithToken(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)

	if err := logSession(token.Raw); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"code": c.Response().StatusCode(), "token": token.Raw, "message": "Login successfully!"})
}

func logSession(token string) error {
	session := model.LoginSession{
		Token: token,
	}

	db := database.DB

	if err := db.Create(&session).Error; err != nil {
		return err
	}

	return nil

}

func verifyGoogleSignIn(c *fiber.Ctx) error {
	db := database.DB
	// Parse the JSON payload from the client-side token
	var jsonPayload struct {
		IDToken string `json:"id_token"`
	}
	if err := c.BodyParser(&jsonPayload); err != nil {
		return err
	}

	// Verify the token and extract the payload
	payload, err := idtoken.Validate(context.Background(), jsonPayload.IDToken, "<your-client-id>")
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Unauthorized", "data": nil})

	}

	// Check if the user is already in your application database
	// If not, retrieve the user's data from the `google-auth-library-go` package to save it to your database
	// ...
	user := model.User{
		Email: payload.Claims["email"].(string),
		Sex:   sex.OTHER,
	}

	if err := db.First("userId = ?", payload.Claims["id"]).Error; err != nil {
		if err := db.Create(&user).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": c.Response().StatusCode(), "message": "Failed Creating User", "data": err})
		}
	}

	// Return a success response
	return c.SendStatus(http.StatusOK)
}
