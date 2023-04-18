// Package Handler
package handler

import (
	"3d-print-account/database"
	"3d-print-account/enum/sex"
	"3d-print-account/model"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func validToken(t *jwt.Token, id string) bool {
	n, err := strconv.Atoi(id)
	if err != nil {
		return false
	}

	claims := t.Claims.(jwt.MapClaims)
	uid := int(claims["user_id"].(float64))

	if uid != n {
		return false
	}

	return true
}

func validUser(username string, password string) bool {
	db := database.DB
	var user model.User
	db.First(&user, username)
	if user.Username == "" {
		return false
	}
	if !CheckPasswordHash(password, user.Password) {
		return false
	}
	return true
}

// Get user by username
func GetUser(c *fiber.Ctx) error {
	username := c.Params("username")
	db := database.DB
	var user model.User
	db.Find(&user, username)
	if user.Username == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "No user found with username", "data": nil})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "user found", "data": user})
}

// Get all users
func GetUsers(c *fiber.Ctx) error {
	db := database.DB
	var users []model.User
	db.Find(&users)
	return c.JSON(fiber.Map{"status": "success", "message": "users found", "data": users})
}

// Register new account
func Register(c *fiber.Ctx) error {
	db := database.DB
	var user model.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"code": c.Response().StatusCode(), "message": "Unprocessable Entity", "data": err})
	}

	response := struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}{
		Username: user.Username,
		Email:    user.Email,
	}

	if govalidator.IsNull(user.Username) || govalidator.IsNull(user.Email) || govalidator.IsNull(user.Password) {
		return c.JSON(fiber.Map{"message": "Empty Fields!", "status": "success", "code": c.Response().StatusCode(), "data": response})
	}

	if !govalidator.IsEmail(user.Email) {
		return c.JSON(fiber.Map{"message": "Invalid Email!", "status": "success", "code": c.Response().StatusCode(), "data": response})
	}

	var result model.User
	if db.First(&result, "user_email = ?", user.Email); result.Username != "" {
		return c.JSON(fiber.Map{"code": c.Response().StatusCode(), "message": "Existed Email", "data": response})
	}

	if db.First(&result, "user_username = ?", user.Username); result.Username != "" {
		return c.JSON(fiber.Map{"code": c.Response().StatusCode(), "message": "Existed Username", "data": response})
	}

	hash, err := hashPassword(user.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": c.Response().StatusCode(), "message": "Failed Hashing Password", "data": err})
	}

	user.Password = hash
	user.Sex = sex.OTHER
	if err := db.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": c.Response().StatusCode(), "message": "Failed Creating User", "data": err})
	}

	return c.JSON(fiber.Map{"code": c.Response().StatusCode(), "message": "Register Successfully", "data": response})

}

// Update account
func UpdateUser(c *fiber.Ctx) error {
	db := database.DB
	id := c.Params("id")
	var user model.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"code": c.Response().StatusCode(), "message": "Unprocessable Entity", "data": err})
	}
	if err := db.Model(&model.User{}).Where("user_id=?", id).Updates(user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": c.Response().StatusCode(), "message": "Failed Updating Account", "data": err})

	}

	// if err := db.Save(&user).Error; err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"code": c.Response().StatusCode(), "message": "Failed Updating Account", "data": err})
	// }
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"code": c.Response().StatusCode(), "message": "success", "data": user})
}
