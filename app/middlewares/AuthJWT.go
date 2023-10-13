package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"os"
)

func AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Printf("running authJWT middleware \n")
		fmt.Printf("%v\n", c.Request.URL)
		fmt.Printf("%v\n", c.Request.UserAgent())

		const BEARER_SCHEMA = "JWT "
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && len(BEARER_SCHEMA) < len(authHeader) {
			tokenString := authHeader[len(BEARER_SCHEMA):]

			key, err := jwt.ParseEdPublicKeyFromPEM([]byte(os.Getenv("JWT_KEY")))
			if err != nil {
				fmt.Printf("%v\n", err)
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			parse, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, isvalid := token.Method.(*jwt.SigningMethodEd25519); !isvalid {
					return nil, fmt.Errorf("Invalid token", token.Header["alg"])
				}

				return key, err
			})
			if err != nil {
				print(err.Error())
				c.AbortWithStatus(403)
				return
			}
			if err == nil && parse.Valid {
				claims := parse.Claims.(jwt.MapClaims)
				c.Set("user_id", int64(claims["user_id"].(float64)))
				c.Set("is_superuser", claims["is_superuser"].(bool))
			} else {
				fmt.Println(err)
				c.AbortWithStatus(http.StatusUnauthorized)
			}
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
