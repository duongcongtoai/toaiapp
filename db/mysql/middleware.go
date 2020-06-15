package mysql

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
)

func MysqlMiddleware(dbKey string, db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(dbKey, db)
			return next(c)
		}
	}
}
