package discourse

import (
	"encoding/json"
	"testing"
)

// sampleTagDataBool provides a minimal JSON with unpinned as a boolean value
const sampleTagDataBool = `{
  "users": [],
  "primary_groups": [],
  "topic_list": {
    "can_create_topic": true,
    "draft": "",
    "draft_key": "test",
    "draft_sequence": 0,
    "per_page": 30,
    "tags": [],
    "topics": [
      {
        "id": 1,
        "title": "Test",
        "fancy_title": "Test",
        "slug": "test",
        "posts_count": 1,
        "reply_count": 0,
        "highest_post_number": 1,
        "image_url": null,
        "created_at": "2025-01-01T00:00:00Z",
        "last_posted_at": "2025-01-01T00:00:00Z",
        "bumped": false,
        "bumped_at": "2025-01-01T00:00:00Z",
        "archetype": "regular",
        "unseen": false,
        "pinned": false,
        "unpinned": false,
        "excerpt": "hello",
        "visible": true,
        "closed": false,
        "archived": false,
        "bookmarked": false,
        "liked": false,
        "tags": [],
        "tags_descriptions": {},
        "like_count": 0,
        "views": 0,
        "category_id": 1,
        "featured_link": null,
        "posters": []
      }
    ]
  }
}`

func TestUnpinnedBool(t *testing.T) {
	var data TagData
	if err := json.Unmarshal([]byte(sampleTagDataBool), &data); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}
	if len(data.TopicList.Topics) != 1 {
		t.Fatalf("expected 1 topic, got %d", len(data.TopicList.Topics))
	}
	if data.TopicList.Topics[0].Unpinned != false {
		t.Fatalf("expected unpinned=false, got %v", data.TopicList.Topics[0].Unpinned)
	}
}
