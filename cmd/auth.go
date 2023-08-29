package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	TIMESTAMP = "TIMESTAMP"
)

type filterHandle func(c *gin.Context) error

var filterMapping = map[string]filterHandle{
	"/v1/example": checkV3Auth,
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if handle, ok := filterMapping[c.Request.URL.Path]; !ok {
			c.AbortWithError(401, fmt.Errorf("%s do not register", c.Request.URL.Path))
			return
		} else {
			if err := handle(c); err != nil {
				Logger.With("module", "authMiddleware", "error", err).Error("failed to check auth")
				c.AbortWithError(401, fmt.Errorf("failed to verity header. err: %v", err))
				return
			}
		}
		c.Next()
	}
}

func noCheck(c *gin.Context) error {
	return nil
}

func checkV3Auth(c *gin.Context) error {
	var key string
	var secret string
	if key = c.GetHeader("http-auth-key"); strings.TrimSpace(key) != Conf.HttpAuth.Key {
		return errors.New("failed to verify key")
	}

	if secret = c.GetHeader("http-auth-secret"); strings.TrimSpace(secret) != Conf.HttpAuth.Secret {
		fmt.Println(strings.TrimSpace(secret), Conf.HttpAuth.Secret)
		return errors.New("failed to verify secret")
	}

	return nil
}
