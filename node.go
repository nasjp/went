package main

const offsetSize = 8

type NodeKind int

const (
	NDUndefined NodeKind = iota
	NDAdd                // +
	NDSub                // -
	NDMul                // *
	NDDiv                // /
	NDEq                 // ==
	NDNe                 // !=
	NDLt                 // <
	NDLe                 // <=
	NDNum                // 123
	NDAssign             // =
	NDLocalV             // ローカル変数
	NDReturn             // return
	NDIf                 // if
	NDFor                // if
	NDBlock              // {}
	NDFuncCall           // 関数呼び出し
	NDFuncDef            // 関数定義
)

type Node struct {
	Kind   NodeKind
	Left   *Node
	Right  *Node
	Cond   *Node
	Then   *Node
	Else   *Node
	Init   *Node
	Inc    *Node
	Body   *Node
	Next   *Node
	Args   *Node
	Params *Node
	Locals *Node
	Root   *Node
	Val    int
	Offset int
	Name   []rune
	Size   int
}

func NewNode(kind NodeKind, left *Node, right *Node) *Node {
	return &Node{
		Kind:  kind,
		Left:  left,
		Right: right,
	}
}

func NewNodeNum(val int) *Node {
	node := NewNode(NDNum, nil, nil)
	node.Val = val

	return node
}

func NewNodeIf(cond *Node, then *Node, els *Node) *Node {
	node := NewNode(NDIf, nil, nil)
	node.Cond = cond
	node.Then = then
	node.Else = els

	return node
}

func NewNodeFor(ini *Node, cond *Node, inc *Node, then *Node) *Node {
	node := NewNode(NDFor, nil, nil)
	node.Init = ini
	node.Cond = cond
	node.Inc = inc
	node.Then = then

	return node
}

func NewNodeFuncCall(name []rune, args *Node) *Node {
	node := NewNode(NDFuncCall, nil, nil)
	node.Name = name
	node.Args = args

	return node
}

func NewNodeFuncDef(name []rune, params *Node, body *Node, locals *Node, size int) *Node {
	node := NewNode(NDFuncDef, nil, nil)
	node.Name = name
	node.Params = params
	node.Body = body
	node.Locals = locals
	node.Size = size

	return node
}

func NewNodeLocalValue(name []rune, offset int) *Node {
	node := NewNode(NDLocalV, nil, nil)
	node.Name = name
	node.Offset = offset

	return node
}

func (node *Node) Len() int {
	return len(node.Name)
}

func (node *Node) FindValue(name []rune) *Node {
	for val := localValue.Root; val != nil; val = val.Next {
		if val.Len() == len(name) && string(name) == string(val.Name) {
			return val
		}
	}

	return nil
}
