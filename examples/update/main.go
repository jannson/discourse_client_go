package main

import (
	"fmt"
	"os"
	"time"

	"github.com/lvoytek/discourse_client_go/pkg/discourse"
)

func main() {

	site := os.Getenv("DISCOURSE_SITE")
	apiKey := os.Getenv("DISCOURSE_APIKEY")
	username := os.Getenv("DISCOURSE_USERNAME")
	discourseClient := discourse.NewClient(site, apiKey, username)

	// 示例：更新一个帖子的内容
	// 步骤 1：获取最新帖子列表，选择一个当前用户可编辑的帖子
	latest, err := discourse.GetLatestPosts(discourseClient)
	if err != nil || latest == nil || len(latest.LatestPosts) == 0 {
		fmt.Println("无法获取最新帖子:", err)
		return
	}

	var editable *discourse.PostData
	for i := range latest.LatestPosts {
		if latest.LatestPosts[i].CanEdit { // 只有你能编辑的帖子才允许更新
			editable = &latest.LatestPosts[i]
			break
		}
	}
	if editable == nil {
		fmt.Println("当前列表中没有可编辑的帖子（可能需要使用有权限的用户或创建一个自己的帖子）")
		return
	}

	fmt.Printf("准备更新帖子 ID=%d 原内容摘要: %.40s\n", editable.ID, editable.Raw)

	// 步骤 2：构造更新内容
	newBody := fmt.Sprintf("这是一条在 %s 通过 Go 客户端更新的示例内容。原文前 40 字: %.40s", time.Now().Format(time.RFC3339), editable.Raw)
	updatePayload := &discourse.UpdatePost{Raw: newBody, EditReason: "演示修改帖子"}

	// 步骤 3：调用更新接口
	updated, err := discourse.UpdatePostByID(discourseClient, editable.ID, updatePayload)
	if err != nil {
		fmt.Println("更新帖子失败:")
		fmt.Println(err)
		return
	}

	fmt.Printf("更新成功！帖子 ID=%d 当前版本=%d\n新内容前 60 字: %.60s\n", updated.ID, updated.Version, updated.Raw)
}
