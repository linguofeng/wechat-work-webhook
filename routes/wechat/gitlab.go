package wechat

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type user struct {
	Username string `json:"username"`
}

type project struct {
	ID   int32  `json:"id"`
	Name string `json:"path_with_namespace"`
	URL  string `json:"web_url"`
}

type assignee struct {
	Username string `json:"username"`
}

type objectAttributes struct {
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Action      string `json:"action"`
}

type payload struct {
	User             user             `json:"user"`
	Project          project          `json:"project"`
	ObjectAttributes objectAttributes `json:"object_attributes"`
	Assignee         assignee         `json:"assignee"`
}

type markdown struct {
	Content string `json:"content"`
}

// 处理合并请求 hook
/**
```markdown
# candyabc/ios 有新的合并请求
标题: 微信登录
描述: 无
提交: @张三
审核: @林国锋
操作: [查看]
`
*/
func handleMergeRequestHook(c echo.Context) error {
	key := c.Param("key")

	payload := new(payload)
	if err := c.Bind(payload); err != nil {
		return err
	}

	// 只处理打开的
	if payload.ObjectAttributes.Action != "open" {
		return c.String(http.StatusOK, "OK")
	}

	description := payload.ObjectAttributes.Description
	if description == "" {
		description = "无"
	}
	content := fmt.Sprint(
		"### [", payload.Project.Name, "](", payload.Project.URL, ") 有新的合并请求\n",
		"> 标题: ", payload.ObjectAttributes.Title, "\n",
		"> 描述: ", description, "\n",
		"> 提交: @", payload.User.Username, "\n",
		"> 审核: @", payload.Assignee.Username, "\n",
		"> 操作: [[查看](", payload.ObjectAttributes.URL, ")]",
	)

	err := send(key, content)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	return c.String(http.StatusOK, "OK")
}

// GitlabHandler 处理 Gitlab
// https://docs.gitlab.com/ee/user/project/integrations/webhooks.html#merge-request-events
func GitlabHandler(c echo.Context) error {
	event := c.Request().Header.Get("X-Gitlab-Event")
	if event == "Merge Request Hook" {
		return handleMergeRequestHook(c)
	}
	return c.String(http.StatusOK, "OK")
}
