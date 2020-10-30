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
	NDReturn                 // return
	NDIf                     // if
	NDFor                    // if
	NDBlock                  // {}
	NDFunc                   // if
)

type Node struct {
	Kind     NodeKind
	Left     *Node
	Right    *Node
	Cond     *Node
	Then     *Node
	Else     *Node
	Init     *Node
	Inc      *Node
	Body     *Node
	Next     *Node
	Args     *Node
	Val      int
	Offset   int
	FuncName []rune
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

func NewNodeFunc(name []rune, args *Node) *Node {
	node := NewNode(NDFunc, nil, nil)
	node.FuncName = name
	node.Args = args

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
	switch {
	case token.Consume(TKReturn):
		proceedToken()

		return stmtReturn()
	case token.Consume(TKIf):
		proceedToken()

		return stmtIf()
	case token.Consume(TKFor):
		proceedToken()

		return stmtFor()
	case token.Consume(TKReserved, '{'):
		proceedToken()

		node := NewNode(NDBlock, nil, nil)

		var stm *Node

		for {
			now, err := stmt()
			if err != nil {
				return nil, err
			}

			if stm == nil {
				stm = now
				node.Body = now
			} else {
				stm.Next = now
				stm = now
			}

			if token.Consume(TKReserved, '}') {
				proceedToken()

				break
			}
		}

		return node, nil
	default:
		node, err := expr()
		if err != nil {
			return nil, err
		}

		if err := token.Expect(TKReserved, ';'); err != nil {
			return nil, err
		}

		proceedToken()

		return node, nil
	}
}

func stmtIf() (*Node, error) {
	// ブロックスコープができたら削る
	if err := token.Expect(TKReserved, '('); err != nil {
		return nil, err
	}

	proceedToken()

	cond, err := expr()
	if err != nil {
		return nil, err
	}

	// ブロックスコープができたら削る
	if err := token.Expect(TKReserved, ')'); err != nil {
		return nil, err
	}

	proceedToken()

	stm, err := stmt()
	if err != nil {
		return nil, err
	}

	var els *Node

	if token.Consume(TKElse) {
		proceedToken()

		var err error
		if els, err = stmt(); err != nil {
			return nil, err
		}
	}

	node := NewNodeIf(cond, stm, els)

	return node, nil
}

func stmtReturn() (*Node, error) {
	left, err := expr()
	if err != nil {
		return nil, err
	}

	node := NewNode(NDReturn, left, nil)

	if err := token.Expect(TKReserved, ';'); err != nil {
		return nil, err
	}

	proceedToken()

	return node, nil
}

func stmtFor() (*Node, error) {
	// ブロックスコープができたら削る
	if err := token.Expect(TKReserved, '('); err != nil {
		return nil, err
	}

	proceedToken()

	var (
		ini  *Node
		cond *Node
		inc  *Node
	)

	// firstNext := token.Skip()

	if token.Consume(TKReserved, ';') {
		proceedToken()
	} else {
		var err error
		if ini, err = expr(); err != nil {
			return nil, err
		}

		if err := token.Expect(TKReserved, ';'); err != nil {
			return nil, err
		}

		proceedToken()
	}

	if token.Consume(TKReserved, ';') {
		proceedToken()
	} else {
		var err error
		if cond, err = expr(); err != nil {
			return nil, err
		}

		if err := token.Expect(TKReserved, ';'); err != nil {
			return nil, err
		}

		proceedToken()
	}

	if token.Consume(TKReserved, ')') {
		proceedToken()
	} else {
		var err error
		if inc, err = expr(); err != nil {
			return nil, err
		}

		if err := token.Expect(TKReserved, ')'); err != nil {
			return nil, err
		}

		proceedToken()
	}

	then, err := stmt()
	if err != nil {
		return nil, err
	}

	return NewNodeFor(ini, cond, inc, then), nil
}

func expr() (*Node, error) {
	return assign()
}

func assign() (*Node, error) {
	node, err := equality()
	if err != nil {
		return nil, err
	}

	if token.Consume(TKReserved, '=') {
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
		if token.Consume(TKReserved, []rune("==")...) {
			proceedToken()

			right, err := relational()
			if err != nil {
				return nil, err
			}

			node = NewNode(NDEq, node, right)

			continue
		}

		if token.Consume(TKReserved, []rune("!=")...) {
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
		if token.Consume(TKReserved, []rune("<")...) {
			proceedToken()

			right, err := add()
			if err != nil {
				return nil, err
			}

			node = NewNode(NDLt, node, right)

			continue
		}

		if token.Consume(TKReserved, []rune("<=")...) {
			proceedToken()

			right, err := add()
			if err != nil {
				return nil, err
			}

			node = NewNode(NDLe, node, right)

			continue
		}

		if token.Consume(TKReserved, []rune(">")...) {
			proceedToken()

			left, err := add()
			if err != nil {
				return nil, err
			}

			node = NewNode(NDLt, left, node)

			continue
		}

		if token.Consume(TKReserved, []rune(">=")...) {
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
		if token.Consume(TKReserved, '+') {
			proceedToken()

			right, err := mul()
			if err != nil {
				return nil, err
			}

			node = NewNode(NDAdd, node, right)

			continue
		}

		if token.Consume(TKReserved, '-') {
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
		if token.Consume(TKReserved, '*') {
			proceedToken()

			right, err := unary()
			if err != nil {
				return nil, err
			}

			node = NewNode(NDMul, node, right)

			continue
		}

		if token.Consume(TKReserved, '/') {
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
	if token.Consume(TKReserved, '+') {
		proceedToken()

		return unary()
	}

	if token.Consume(TKReserved, '-') {
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
	if token.Consume(TKReserved, '(') {
		proceedToken()

		node, err := expr()
		if err != nil {
			return nil, err
		}

		if err := token.Expect(TKReserved, ')'); err != nil {
			return nil, err
		}

		proceedToken()

		return node, nil
	}

	if token.ConsumeIdent() {
		return ident()
	}

	n, err := token.ExpectNum()
	if err != nil {
		return nil, err
	}

	proceedToken()

	return NewNodeNum(n), nil
}

func ident() (*Node, error) {
	if token.Skip().Consume(TKReserved, '(') {
		return identFunc()
	}

	node, err := identVal()

	proceedToken()

	return node, err
}

func identVal() (*Node, error) {
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

	return NewNodeIdent(offset), nil
}

func identFunc() (*Node, error) {
	funcName := token.Str

	proceedToken()

	if err := token.Expect(TKReserved, '('); err != nil {
		return nil, err
	}

	proceedToken()

	args, err := funcArgs()
	if err != nil {
		return nil, err
	}

	return NewNodeFunc(funcName, args), nil
}

func funcArgs() (*Node, error) {
	if token.Consume(TKReserved, ')') {
		proceedToken()

		return nil, nil
	}

	head, err := assign()
	if err != nil {
		return nil, err
	}

	cur := head

	for token.Consume(TKReserved, ',') {
		proceedToken()

		node, err := assign()
		if err != nil {
			return nil, err
		}

		cur.Next = node

		cur = cur.Next
	}

	if err := token.Expect(TKReserved, ')'); err != nil {
		return nil, err
	}

	proceedToken()

	return head, nil
}
