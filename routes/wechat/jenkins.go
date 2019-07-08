package wechat

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

// JenkinsHandler 处理 Jenkins
func JenkinsHandler(c echo.Context) error {
	key := c.Param("key")
	tokey := c.QueryParam("token")
	if tokey != os.Getenv("JENKINS_TOKEN") {
		return nil
	}

	type scm struct {
		Branch string `json:"branch"`
		Commit string `json:"commit"`
	}

	type build struct {
		URL       string `json:"full_url"`
		Number    int32  `json:"number"`
		Timestamp int64  `json:"timestamp"`
		Status    string `json:"status"`
		Scm       scm    `json:"scm"`
	}

	payload := new(struct {
		Name  string `json:"name"`
		Build build  `json:"build"`
	})
	if err := c.Bind(payload); err != nil {
		return err
	}

	status := "构建成功"
	if payload.Build.Status == "FAILURE" {
		status = "构建失败"
	}
	t := time.Unix(payload.Build.Timestamp/1000, 0)
	content := fmt.Sprint(
		"### 【", payload.Build.Status, "】", payload.Name, "\n",
		"> 状态: ", status, "\n",
		"> 时间: ", t.Format("2006-01-02 15:04:05"), "\n",
		"> 分支: ", payload.Build.Scm.Branch, "(", payload.Build.Scm.Commit, ")\n",
		"> 操作: [[查看](", payload.Build.URL, ")]\n",
	)

	err := send(key, content)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	return c.String(http.StatusOK, "OK")
}
