package gee

import (
	"strings"
)

type Node struct {
	pattern  string  //待匹配的路由
	part     string  //路由中的一部分
	children []*Node //子节点
	isWild   bool    //是否模糊匹配
}

//匹配一个节点
func (n *Node) matchChild(part string) *Node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

//匹配多个节点
func (n *Node) matchChildren(part string) []*Node {
	nodes := make([]*Node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *Node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &Node{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*',
		}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

func (n *Node) search(parts []string, height int) *Node {
	if height == len(parts) || strings.HasPrefix(n.part, "*") { //搜寻到最后一个part的节点
		if n.pattern == "" { //说明没有这个路由
			return nil
		}
		return n //说明匹配到了路由节点
	}

	part := parts[height]
	children := n.matchChildren(part)
	for _, child := range children {
		result := child.search(parts, height+1) //深度优先搜索
		if result != nil {
			return result
		}
	}

	return nil
}
