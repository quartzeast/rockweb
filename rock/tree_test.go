package rock

import (
	"reflect"
	"testing"
)

func TestParsePath(t *testing.T) {
	tests := []struct {
		name string
		path string
		want []string
	}{
		{
			name: "simple path",
			path: "/user/profile",
			want: []string{"user", "profile"},
		},
		{
			name: "path with parameter",
			path: "/user/:id",
			want: []string{"user", ":id"},
		},
		{
			name: "path with catch-all",
			path: "/static/*filepath",
			want: []string{"static", "*filepath"},
		},
		{
			name: "root path",
			path: "/",
			want: []string{},
		},
		{
			name: "path with trailing slash",
			path: "/user/profile/",
			want: []string{"user", "profile"},
		},
		{
			name: "complex path",
			path: "/api/v1/user/:id/posts",
			want: []string{"api", "v1", "user", ":id", "posts"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parsePath(tt.path)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parsePath(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestNodeMatchChild(t *testing.T) {
	root := &node{
		children: []*node{
			{part: "user"},
			{part: "post"},
			{part: ":id", isWild: true},
		},
	}

	tests := []struct {
		name     string
		part     string
		wantPart string
		wantNil  bool
	}{
		{
			name:     "match existing child",
			part:     "user",
			wantPart: "user",
		},
		{
			name:     "match another child",
			part:     "post",
			wantPart: "post",
		},
		{
			name:     "match wildcard child",
			part:     ":id",
			wantPart: ":id",
		},
		{
			name:    "no match",
			part:    "comment",
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := root.matchChild(tt.part)
			if tt.wantNil {
				if got != nil {
					t.Errorf("matchChild(%q) = %v, want nil", tt.part, got)
				}
				return
			}
			if got == nil {
				t.Errorf("matchChild(%q) = nil, want %q", tt.part, tt.wantPart)
				return
			}
			if got.part != tt.wantPart {
				t.Errorf("matchChild(%q).part = %q, want %q", tt.part, got.part, tt.wantPart)
			}
		})
	}
}

func TestNodeMatchChildren(t *testing.T) {
	root := &node{
		children: []*node{
			{part: "user"},
			{part: "post"},
			{part: ":id", isWild: true},
		},
	}

	tests := []struct {
		name      string
		part      string
		wantParts []string
	}{
		{
			name:      "match static and wildcard",
			part:      "user",
			wantParts: []string{"user", ":id"},
		},
		{
			name:      "match only wildcard",
			part:      "123",
			wantParts: []string{":id"},
		},
		{
			name:      "match static post and wildcard",
			part:      "post",
			wantParts: []string{"post", ":id"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := root.matchChildren(tt.part)
			if len(got) != len(tt.wantParts) {
				t.Errorf("matchChildren(%q) returned %d nodes, want %d", tt.part, len(got), len(tt.wantParts))
				return
			}
			for i, n := range got {
				if n.part != tt.wantParts[i] {
					t.Errorf("matchChildren(%q)[%d].part = %q, want %q", tt.part, i, n.part, tt.wantParts[i])
				}
			}
		})
	}
}

func TestNodeInsertAndSearch(t *testing.T) {
	root := &node{}

	// Insert routes
	routes := []string{
		"/",
		"/user",
		"/user/profile",
		"/user/:id",
		"/user/:id/posts",
		"/post/:postId/comments",
		"/static/*filepath",
	}

	for _, route := range routes {
		parts := parsePath(route)
		root.insert(route, parts, 0)
	}

	tests := []struct {
		name     string
		path     string
		wantPart string
		wantNil  bool
	}{
		{
			name:     "match static route /user",
			path:     "/user",
			wantPart: "user",
		},
		{
			name:     "match static route /user/profile",
			path:     "/user/profile",
			wantPart: "profile",
		},
		{
			name:     "match param route /user/:id",
			path:     "/user/123",
			wantPart: ":id",
		},
		{
			name:     "match param route /user/:id/posts",
			path:     "/user/456/posts",
			wantPart: "posts",
		},
		{
			name:     "match another param route",
			path:     "/post/789/comments",
			wantPart: "comments",
		},
		{
			name:     "match catch-all route",
			path:     "/static/css/style.css",
			wantPart: "*filepath",
		},
		{
			name:     "match catch-all with deep path",
			path:     "/static/js/lib/jquery.min.js",
			wantPart: "*filepath",
		},
		{
			name:    "no match",
			path:    "/unknown/path",
			wantNil: true,
		},
		{
			name:    "partial match should not match",
			path:    "/user/123/unknown",
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parts := parsePath(tt.path)
			got := root.search(parts, 0)
			if tt.wantNil {
				if got != nil {
					t.Errorf("search(%q) = %v, want nil", tt.path, got)
				}
				return
			}
			if got == nil {
				t.Errorf("search(%q) = nil, want node with part %q", tt.path, tt.wantPart)
				return
			}
			if got.part != tt.wantPart {
				t.Errorf("search(%q).part = %q, want %q", tt.path, got.part, tt.wantPart)
			}
		})
	}
}

func TestNodeInsertDuplicates(t *testing.T) {
	root := &node{}

	// Insert same route twice
	parts := parsePath("/user/profile")
	root.insert("/user/profile", parts, 0)
	root.insert("/user/profile", parts, 0)

	// Should still only have one path
	if len(root.children) != 1 {
		t.Errorf("expected 1 child after duplicate insert, got %d", len(root.children))
	}
	if root.children[0].part != "user" {
		t.Errorf("expected child part 'user', got %q", root.children[0].part)
	}
}

func TestRootPathSearch(t *testing.T) {
	root := &node{isEnd: true}

	parts := parsePath("/")
	got := root.search(parts, 0)

	if got == nil {
		t.Error("search('/') = nil, want root node")
		return
	}
	if got != root {
		t.Error("search('/') did not return root node")
	}
}
