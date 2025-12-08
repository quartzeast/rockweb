package rock

import "strings"

// node represents a node in the prefix tree for route matching.
type node struct {
	pattern  string  // full pattern of the route (only set for end nodes)
	part     string  // path segment of this node (e.g., "user", ":id", "*")
	children []*node // child nodes
	isWild   bool    // whether this node is a wildcard (starts with ':' or '*')
	isEnd    bool    // whether this node represents a complete route
}

// matchChild finds a child node that exactly matches the given part.
// Used during insertion.
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part {
			return child
		}
	}
	return nil
}

// matchChildren finds all child nodes that could match the given part.
// Used during search.
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// insert adds a route pattern to the tree.
// Pattern examples:
//   - "/user/profile" - static route
//   - "/user/:id" - route with parameter
//   - "/static/*filepath" - route with catch-all
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		n.isEnd = true
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*',
		}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

// search finds a matching node for the given path parts.
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.isEnd {
			return n
		}
		return nil
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}

// parsePath splits a URL path into parts.
// Example: "/user/:id/profile" -> ["user", ":id", "profile"]
func parsePath(path string) []string {
	vs := strings.Split(path, "/")
	parts := make([]string, 0)
	for _, v := range vs {
		if v != "" {
			parts = append(parts, v)
			// Stop at catch-all
			if v[0] == '*' {
				break
			}
		}
	}
	return parts
}
