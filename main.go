package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"gitlab.zhixuehd.com/linguofeng/webhook/routes/wechat"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)

	// 处理 Gitlab hook
	e.POST("/wechat/:key/gitlab", wechat.GitlabHandler, middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:X-Gitlab-Token",
		Validator: func(key string, _ echo.Context) (bool, error) {
			return key == os.Getenv("GITLAB_TOKEN"), nil
		},
	}))

	// 处理 Jenkins hook
	e.POST("/wechat/:key/jenkins", wechat.JenkinsHandler)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
