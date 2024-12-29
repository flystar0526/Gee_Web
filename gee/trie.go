package gee

import "strings"

type node struct {
	path     string  // Route to be matched, ex : /p/:lang
	part     string  // A part of the route, ex : :lang
	children []*node // Child nodes, ex :, [doc, tutorial, intro]
	isWild   bool    // Indicates whether it is an exact match; true if part contains : or *
}

// The first successfully matched node, used for insertion
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// All successfully matched nodes, used for searching
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)

	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}

	return nodes
}

func (n *node) insert(path string, parts []string, height int) {
	if len(parts) == height {
		n.path = path
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

	child.insert(path, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.path == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		node := child.search(parts, height+1)
		if node != nil {
			return node
		}
	}

	return nil
}
