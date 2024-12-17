package middleware

import (
	"errors"
	"meals-app/exception"
	"meals-app/model/entity"
	"os"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func Protected(db *gorm.DB) fiber.Handler{
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: []byte(os.Getenv("JWT_TOKEN_SECRET")),
		},
		ErrorHandler: jwtError,
		SuccessHandler: func(c *fiber.Ctx) error {
			return jwtSuccess(c, db)
		},
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return exception.ErrorHandler(400, "BAD REQUEST", errors.New("missing or malformed JWT"))(c)
	}
	return exception.ErrorHandler(401, "UNAUTHORIZE", errors.New("unauthorize, login instead"))(c)
}

func jwtSuccess(c *fiber.Ctx, db *gorm.DB) error {
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	userId := int(claims["user_id"].(float64))
	tokenVersion := int(claims["token_version"].(float64))

	user := entity.User{}
	err := db.Take(&user, "id = ?", userId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exception.ErrorHandler(401, "UNAUTHORIZE", errors.New("user not found"))(c)
		}
		return exception.ErrorHandler(500, "INTERNAL SERVER ERROR", err)(c)
	}

	if tokenVersion != user.TokenVersion {
		return exception.ErrorHandler(401, "UNAUTHORIZE", errors.New("token revoked"))(c)
	}

	c.Locals("currentUser", user)

	return c.Next()
}
