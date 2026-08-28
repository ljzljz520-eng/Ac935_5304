package hospitaldesk

import (
	"hospitaldesk/model"
	"hospitaldesk/storage"
	"testing"
)

func TestStorageSearch(t *testing.T) {
	store, err := storage.Open(t.TempDir() + "/search.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	if err := store.SavePolicy(model.PolicyDocument{ID: "p1", Title: "Medication", Department: "x", Content: "safe", Status: model.PolicyPublished}); err != nil {
		t.Fatal(err)
	}
	items, err := store.SearchPolicies(model.SearchFilter{Query: "medication", Department: "x"})
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatalf("len=%d", len(items))
	}
}
