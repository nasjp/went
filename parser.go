package main

type NodeKind int

const (
	NDAdd    NodeKind = iota // +
	NDSub                    // -
	NDMul                    // *
	NDDiv                    // /
	NDEq                     // ==
	NDNe                     // !=
	NDLt                     // <
	NDLe                     // <=
	NDNum                    // 123
	NDAssign                 // =
	NDLocalV                 // ローカル変数
	NDReturn
)

type Node struct {
	Kind   NodeKind
	Left   *Node
	Right  *Node
	Val    int
	Offset int
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

func NewNodeIdent(offset int) *Node {
	node := NewNode(NDLocalV, nil, nil)
	node.Offset = offset

	return node
}

func program() ([]*Node, error) {
	code := make([]*Node, 0, 100)

	for !token.AtEOF() {
		node, err := stmt()
		if err != nil {
			return nil, err
		}

		code = append(code, node)
	}

	return code, nil
}

func stmt() (*Node, error) {
	var (
		node *Node
		err  error
	)

	if token.Consume([]rune("return")...) {
		proceedToken()

		left, err := expr()
		if err != nil {
			return nil, err
		}

		node = NewNode(NDReturn, left, nil)
	} else {
		node, err = expr()
		if err != nil {
			return nil, err
		}
	}

	if err := token.Expect(';'); err != nil {
		return nil, err
	}

	proceedToken()

	return node, nil
}

func expr() (*Node, error) {
	return assign()
}

func assign() (*Node, error) {
	node, err := equality()
	if err != nil {
		return nil, err
	}

	if token.Consume('=') {
		proceedToken()

		right, err := equality()
		if err != nil {
			return nil, err
		}

		return NewNode(NDAssign, node, right), nil
	}

	return node, nil
}

func equality() (*Node, error) {
	node, err := relational()
	if err != nil {
		return nil, err
	}

	for {
		if token.Consume([]rune("==")...) {
			proceedToken()

			right, err := relational()
			if err != nil {
				return nil, err
			}

			node = NewNode(NDEq, node, right)

			continue
		}

		if token.Consume([]rune("!=")...) {
			proceedToken()

			right, err := relational()
			if err != nil {
				return nil, err
			}

			node = NewNode(NDNe, node, right)

			continue
		}

		return node, nil
	}
}

func relational() (*Node, error) {
	node, err := add()
	if err != nil {
		return nil, err
	}

	for {
		if token.Consume([]rune("<")...) {
			proceedToken()

			right, err := add()
			if err != nil {
				return nil, err
			}

			node = NewNode(NDLt, node, right)

			continue
		}

		if token.Consume([]rune("<=")...) {
			proceedToken()

			right, err := add()
			if err != nil {
				return nil, err
			}

			node = NewNode(NDLe, node, right)

			continue
		}

		if token.Consume([]rune(">")...) {
			proceedToken()

			left, err := add()
			if err != nil {
				return nil, err
			}

			node = NewNode(NDLt, left, node)

			continue
		}

		if token.Consume([]rune(">=")...) {
			proceedToken()

			left, err := add()
			if err != nil {
				return nil, err
			}

			node = NewNode(NDLe, left, node)

			continue
		}

		return node, nil
	}
}

func add() (*Node, error) {
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

		return unary()
	}

	if token.Consume('-') {
		proceedToken()

		node, err := unary()
		if err != nil {
			return nil, err
		}

		return NewNode(NDSub, NewNodeNum(0), node), nil
	}

	return primary()
}

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

	if token.ConsumeIdent() {
		var offset int

		if lval := findLocalValue(token); lval != nil {
			offset = lval.Offset
		} else {
			var localValueOffset int
			if localValue != nil {
				localValueOffset = localValue.Offset
			}

			localValue = NewLocalValue(token, localValueOffset)
			offset = localValue.Offset
		}

		proceedToken()

		return NewNodeIdent(offset), nil
	}

	n, err := token.ExpectNum()
	if err != nil {
		return nil, err
	}

	proceedToken()

	return NewNodeNum(n), nil
}
