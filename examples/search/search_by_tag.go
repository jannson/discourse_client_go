package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/lvoytek/discourse_client_go/pkg/discourse"
)

func searchByTag() {
	site := os.Getenv("DISCOURSE_SITE")
	apiKey := os.Getenv("DISCOURSE_APIKEY")
	username := os.Getenv("DISCOURSE_USERNAME")

	fmt.Printf("正在连接到: %s\n", site)
	fmt.Printf("用户名: %s\n\n", username)

	discourseClient := discourse.NewClient(site, apiKey, username)

	// 命令行参数
	tagName := flag.String("tag", "知识库", "要搜索的标签名称")
	page := flag.Int("page", 0, "页码（从0开始）")
	showAll := flag.Bool("all", false, "显示所有页面")
	flag.Parse()

	fmt.Printf("=== 搜索标签为 '%s' 的主题 ===\n", *tagName)

	if *showAll {
		// 获取所有页面
		fetchAllPages(discourseClient, *tagName, site)
	} else {
		// 获取指定页面
		fetchSinglePage(discourseClient, *tagName, *page, site)
	}
}

func fetchSinglePage(client *discourse.Client, tagName string, page int, site string) {
	fmt.Printf("\n--- 第 %d 页 ---\n", page)

	var tagData *discourse.TagData
	var err error

	if page == 0 {
		tagData, err = discourse.GetTagByName(client, tagName)
	} else {
		tagData, err = discourse.GetTagByNameWithPage(client, tagName, page)
	}

	if err != nil {
		fmt.Printf("❌ 错误: %v\n", err)
		fmt.Println("\n可能的原因:")
		fmt.Println("  1. 标签不存在")
		fmt.Println("  2. API 密钥权限不足")
		fmt.Println("  3. 页码超出范围")
		return
	}

	topics := tagData.TopicList.Topics
	fmt.Printf("✅ 当前页找到 %d 个主题\n", len(topics))
	fmt.Printf("每页显示: %d 个\n\n", tagData.TopicList.PerPage)

	if len(topics) == 0 {
		fmt.Println("⚠️  当前页没有主题")
		return
	}

	displayTopics(topics, site)
}

func fetchAllPages(client *discourse.Client, tagName string, site string) {
	page := 0
	totalTopics := 0

	for {
		fmt.Printf("\n--- 第 %d 页 ---\n", page)

		var tagData *discourse.TagData
		var err error

		if page == 0 {
			tagData, err = discourse.GetTagByName(client, tagName)
		} else {
			tagData, err = discourse.GetTagByNameWithPage(client, tagName, page)
		}

		if err != nil {
			fmt.Printf("❌ 获取第 %d 页时出错: %v\n", page, err)
			break
		}

		topics := tagData.TopicList.Topics
		if len(topics) == 0 {
			fmt.Println("✓ 已到达最后一页")
			break
		}

		fmt.Printf("找到 %d 个主题\n", len(topics))
		displayTopics(topics, site)

		totalTopics += len(topics)

		// 如果返回的主题数少于每页数量，说明这是最后一页
		if len(topics) < tagData.TopicList.PerPage {
			fmt.Println("✓ 已到达最后一页")
			break
		}

		page++
	}

	fmt.Println("\n========================================")
	fmt.Printf("总计获取了 %d 个主题，共 %d 页\n", totalTopics, page+1)
}

func displayTopics(topics []discourse.SuggestedTopic, site string) {
	fmt.Println("主题列表:")
	fmt.Println("----------------------------------------")
	for i, topic := range topics {
		fmt.Printf("%d. %s\n", i+1, topic.Title)
		fmt.Printf("   ID: %d | 浏览: %d | 回复: %d | 点赞: %d\n",
			topic.ID, topic.Views, topic.PostsCount-1, topic.LikeCount)

		// 获取最后回复者
		lastPoster := "未知"
		if len(topic.Posters) > 0 {
			lastPoster = topic.Posters[len(topic.Posters)-1].User.Username
		}

		fmt.Printf("   最后回复: %s 于 %s\n",
			lastPoster, topic.LastPostedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("   URL: %s/t/%s/%d\n", site, topic.Slug, topic.ID)

		// 显示标签
		if len(topic.Tags) > 0 {
			fmt.Printf("   标签: %v\n", topic.Tags)
		}

		fmt.Println()
	}
}
