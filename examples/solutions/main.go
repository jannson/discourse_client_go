package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/lvoytek/discourse_client_go/pkg/discourse"
)

// Example: Fetch the Post marked as the accepted solution for a Topic
// Works when the Discourse instance has the "discourse-solved" plugin enabled.
//
// You can override the site and topic via env:
//
//	SOLVED_EXAMPLE_SITE=https://meta.discourse.org SOLVED_EXAMPLE_TOPIC_ID=30155
func main() {
	// Default to Meta, which has the solved plugin enabled
	site := getenvDefault("SOLVED_EXAMPLE_SITE", "https://meta.discourse.org")
	// Pick a topic id; override with env if desired
	topicID := getenvIntDefault("SOLVED_EXAMPLE_TOPIC_ID", 30155)

	client := discourse.NewAnonymousClient(site)
	topic, err := discourse.GetTopicByID(client, topicID)
	if err != nil {
		fmt.Println("Error fetching topic:", err)
		return
	}

	fmt.Printf("Topic #%d: %s\n", topicID, topic.Title)

	var solutionPost *discourse.PostData
	var solutionPostID int
	var solutionPostNumber int

	// Prefer topic-level accepted_answer summary if available
	if topic.AcceptedAnswer != nil {
		solutionPostNumber = topic.AcceptedAnswer.PostNumber
		// Try to find it in the included posts
		for i := range topic.PostStream.Posts {
			p := &topic.PostStream.Posts[i]
			if p.PostNumber == solutionPostNumber {
				solutionPost = p
				break
			}
		}
		// If not present, map post_number -> id using the stream order
		if solutionPost == nil {
			if solutionPostNumber > 0 && solutionPostNumber <= len(topic.PostStream.Stream) {
				solutionPostID = topic.PostStream.Stream[solutionPostNumber-1]
			}
		}
	} else {
		// Fallback: check post-level flags from the solved plugin
		for i := range topic.PostStream.Posts {
			p := &topic.PostStream.Posts[i]
			if p.AcceptedAnswer {
				solutionPost = p
				solutionPostNumber = p.PostNumber
				break
			}
		}
		if solutionPost == nil && len(topic.PostStream.Posts) > 0 {
			// If we know a topic has any accepted answer, some sites expose topic_accepted_answer
			for i := range topic.PostStream.Posts {
				p := &topic.PostStream.Posts[i]
				if p.TopicAcceptedAnswer && p.PostNumber > 1 { // answers are never post_number 1
					solutionPostNumber = p.PostNumber // best-effort
					break
				}
			}
			if solutionPostNumber > 0 && solutionPostNumber <= len(topic.PostStream.Stream) {
				solutionPostID = topic.PostStream.Stream[solutionPostNumber-1]
			}
		}
	}

	// Fetch the full post if we only have an ID
	if solutionPost == nil && solutionPostID != 0 {
		p, err := discourse.GetPostByID(client, solutionPostID)
		if err == nil {
			solutionPost = p
		}
	}

	if solutionPost == nil {
		fmt.Println("No accepted solution found for this topic (or plugin disabled).")
		return
	}

	fmt.Printf("Accepted solution: post_number=%d id=%d by %s\n", solutionPost.PostNumber, solutionPost.ID, displayName(solutionPost))
	// Build a short preview preferring raw -> excerpt -> cooked
	preview := solutionPost.Raw
	if preview == "" && topic.AcceptedAnswer != nil && topic.AcceptedAnswer.Excerpt != "" {
		preview = topic.AcceptedAnswer.Excerpt
	}
	if preview == "" {
		preview = solutionPost.Cooked // contains HTML
	}
	if len(preview) > 280 {
		preview = preview[:280] + "..."
	}
	fmt.Println("Preview:")
	fmt.Println(preview)
}

func displayName(p *discourse.PostData) string {
	if p.DisplayUsername != "" {
		return p.DisplayUsername
	}
	if p.Name != "" {
		return p.Name
	}
	return p.Username
}

func getenvDefault(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

func getenvIntDefault(k string, d int) int {
	if v := os.Getenv(k); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return d
}
