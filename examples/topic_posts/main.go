package main

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/lvoytek/discourse_client_go/pkg/discourse"
)

// FetchAllPostsForTopic retrieves all posts belonging to a topic.
// Strategy:
//  1. Call GetTopicByID to get initial PostStream (first chunk of posts + full stream of post IDs).
//  2. Build a map of already loaded post IDs.
//  3. For remaining IDs in PostStream.Stream, fetch each post via GetPostByID.
//     (Discourse also supports batched loading, but this client doesn't expose a batch endpoint yet.)
//  4. Return posts sorted by PostNumber.
func FetchAllPostsForTopic(client *discourse.Client, topicID int) ([]discourse.PostData, error) {
	// Initial topic fetch
	topic, err := discourse.GetTopicByID(client, topicID)
	if err != nil {
		return nil, fmt.Errorf("failed to get topic %d: %w", topicID, err)
	}

	loaded := make(map[int]discourse.PostData, len(topic.PostStream.Posts))
	for _, p := range topic.PostStream.Posts {
		loaded[p.ID] = p
	}

	// Load missing posts individually
	for _, postID := range topic.PostStream.Stream {
		if _, exists := loaded[postID]; exists {
			continue
		}
		p, perr := discourse.GetPostByID(client, postID)
		if perr != nil {
			// 404/403 may happen for deleted/hidden posts; skip but log
			log.Printf("warn: could not load post %d in topic %d: %v", postID, topicID, perr)
			continue
		}
		loaded[p.ID] = *p
		// Optional tiny delay to avoid hammering the server (tune as needed)
		time.Sleep(30 * time.Millisecond)
	}

	// Collect & sort by post number
	result := make([]discourse.PostData, 0, len(loaded))
	for _, p := range loaded {
		result = append(result, p)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].PostNumber < result[j].PostNumber })

	return result, nil
}

func main() {
	// Replace with a real topic ID from your target Discourse instance
	const topicID = 22706

	// Anonymous client is enough for public topics
	client := discourse.NewAnonymousClient("https://meta.discourse.org")

	posts, err := FetchAllPostsForTopic(client, topicID)
	if err != nil {
		log.Fatalf("error fetching posts: %v", err)
	}

	fmt.Printf("Total posts fetched for topic %d: %d\n", topicID, len(posts))
	for _, p := range posts {
		// Show basic info; Raw contains the original markdown text, Cooked contains HTML
		fmt.Printf("#%d by %s (id=%d) created %s\n", p.PostNumber, p.Username, p.ID, p.CreatedAt.Format(time.RFC3339))
	}
}
