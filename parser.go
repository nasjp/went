package main

type NodeKind int

const (
	NDAdd NodeKind = iota
	NDSub
	NDMul
	NDDiv
	NDNum
)

type Node struct {
	Kind  NodeKind
	Left  *Node
	Right *Node
	Val   int
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

// expr    = mul ("+" mul | "-" mul)*
func expr() (*Node, error) {
	node, err := mul()
	if err != nil {
		return nil, err
	}

	for {
		if token.Consume('+') {
			proceedToken()

			right, err := mul()
			if err != nil {
				return nil, err
			}

			node = NewNode(NDAdd, node, right)

			continue
		}

		if token.Consume('-') {
			proceedToken()

			right, err := mul()
			if err != nil {
				return nil, err
			}

			node = NewNode(NDSub, node, right)

			continue
		}

		return node, nil
	}
}

// mul     = unary ("*" unary | "/" unary)*
func mul() (*Node, error) {
	node, err := unary()
	if err != nil {
		return nil, err
	}

	for {
		if token.Consume('*') {
			proceedToken()

			right, err := unary()
			if err != nil {
				return nil, err
			}

			node = NewNode(NDMul, node, right)

			continue
		}

		if token.Consume('/') {
			proceedToken()

			right, err := unary()
			if err != nil {
				return nil, err
			}

			node = NewNode(NDDiv, node, right)

			continue
		}

		return node, nil
	}
}

func unary() (*Node, error) {
	if token.Consume('+') {
		proceedToken()

		return primary()
	}

	if token.Consume('-') {
		proceedToken()

		node, err := primary()
		if err != nil {
			return nil, err
		}

		return NewNode(NDSub, NewNodeNum(0), node), nil
	}

	return primary()
}

// primary = num | "(" expr ")"
func primary() (*Node, error) {
	if token.Consume('(') {
		proceedToken()

		node, err := expr()
		if err != nil {
			return nil, err
		}

		if err := token.Expect(')'); err != nil {
			return nil, err
		}

		proceedToken()

		return node, nil
	}

	n, err := token.ExpectNum()
	if err != nil {
		return nil, err
	}

	proceedToken()

	return NewNodeNum(n), nil
}
