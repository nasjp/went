package main

func parse() (*Node, error) {
	head := NewNode(NDUndefined, nil, nil)
	cur := head

	for !currentToken.AtEOF() {
		node, err := function()
		if err != nil {
			return nil, err
		}

		cur.Next = node
		cur = node
	}

	return head.Next, nil
}

func function() (*Node, error) {
	if err := currentToken.Expect(TKIdent); err != nil {
		return nil, err
	}

	funcName := currentToken.Str

	proceedToken()

	if err := currentToken.Expect(TKReserved, '('); err != nil {
		return nil, err
	}

	proceedToken()

	params, err := funcParams()
	if err != nil {
		return nil, err
	}

	if err := currentToken.Expect(TKReserved, '{'); err != nil {
		return nil, err
	}

	proceedToken()

	head := NewNode(NDLocalV, nil, nil)
	localValue = head
	localValue.Root = head

	body, err := block()
	if err != nil {
		return nil, err
	}

	var localNum int
	for local := head.Next; local != nil; local = local.Next {
		localNum++
	}

	return NewNodeFuncDef(funcName, params, body, head.Next, localNum*offsetSize), nil
}

func funcParams() (*Node, error) {
	if currentToken.Consume(TKReserved, ')') {
		proceedToken()

		return nil, nil
	}

	head, err := assign()
	if err != nil {
		return nil, err
	}

	cur := head

	for currentToken.Consume(TKReserved, ',') {
		proceedToken()

		node, err := assign()
		if err != nil {
			return nil, err
		}

		cur.Next = node

		cur = cur.Next
	}

	if err := currentToken.Expect(TKReserved, ')'); err != nil {
		return nil, err
	}

	proceedToken()

	return head, nil
}

func block() (*Node, error) {
	node := NewNode(NDBlock, nil, nil)

	head := NewNode(NDBlock, nil, nil)

	cur := head

	for {
		node, err := stmt()
		if err != nil {
			return nil, err
		}

		cur.Next = node
		cur = node

		if currentToken.Consume(TKReserved, '}') {
			proceedToken()

			break
		}
	}

	node.Body = head.Next

	return node, nil
}

func stmt() (*Node, error) {
	switch {
	case currentToken.Consume(TKReturn):
		proceedToken()

		return stmtReturn()
	case currentToken.Consume(TKIf):
		proceedToken()

		return stmtIf()
	case currentToken.Consume(TKFor):
		proceedToken()

		return stmtFor()
	case currentToken.Consume(TKReserved, '{'):
		proceedToken()

		return block()
	default:
		node, err := expr()
		if err != nil {
			return nil, err
		}

		if err := currentToken.Expect(TKReserved, ';'); err != nil {
			return nil, err
		}

		proceedToken()

		return node, nil
	}
}

func stmtIf() (*Node, error) {
	// ブロックスコープができたら削る
	if err := currentToken.Expect(TKReserved, '('); err != nil {
		return nil, err
	}

	proceedToken()

	cond, err := expr()
	if err != nil {
		return nil, err
	}

	// ブロックスコープができたら削る
	if err := currentToken.Expect(TKReserved, ')'); err != nil {
		return nil, err
	}

	proceedToken()

	stm, err := stmt()
	if err != nil {
		return nil, err
	}

	var els *Node

	if currentToken.Consume(TKElse) {
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

	if err := currentToken.Expect(TKReserved, ';'); err != nil {
		return nil, err
	}

	proceedToken()

	return node, nil
}

func stmtFor() (*Node, error) {
	// ブロックスコープができたら削る
	if err := currentToken.Expect(TKReserved, '('); err != nil {
		return nil, err
	}

	proceedToken()

	var (
		ini  *Node
		cond *Node
		inc  *Node
	)

	if currentToken.Consume(TKReserved, ';') {
		proceedToken()
	} else {
		var err error
		if ini, err = expr(); err != nil {
			return nil, err
		}

		if err := currentToken.Expect(TKReserved, ';'); err != nil {
			return nil, err
		}

		proceedToken()
	}

	if currentToken.Consume(TKReserved, ';') {
		proceedToken()
	} else {
		var err error
		if cond, err = expr(); err != nil {
			return nil, err
		}

		if err := currentToken.Expect(TKReserved, ';'); err != nil {
			return nil, err
		}

		proceedToken()
	}

	if currentToken.Consume(TKReserved, ')') {
		proceedToken()
	} else {
		var err error
		if inc, err = expr(); err != nil {
			return nil, err
		}

		if err := currentToken.Expect(TKReserved, ')'); err != nil {
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

	if currentToken.Consume(TKReserved, '=') {
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
		if currentToken.Consume(TKReserved, []rune("==")...) {
			proceedToken()

			right, err := relational()
			if err != nil {
				return nil, err
			}

			node = NewNode(NDEq, node, right)

			continue
		}

		if currentToken.Consume(TKReserved, []rune("!=")...) {
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
		if currentToken.Consume(TKReserved, []rune("<")...) {
			proceedToken()

			right, err := add()
			if err != nil {
				return nil, err
			}

			node = NewNode(NDLt, node, right)

			continue
		}

		if currentToken.Consume(TKReserved, []rune("<=")...) {
			proceedToken()

			right, err := add()
			if err != nil {
				return nil, err
			}

			node = NewNode(NDLe, node, right)

			continue
		}

		if currentToken.Consume(TKReserved, []rune(">")...) {
			proceedToken()

			left, err := add()
			if err != nil {
				return nil, err
			}

			node = NewNode(NDLt, left, node)

			continue
		}

		if currentToken.Consume(TKReserved, []rune(">=")...) {
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
		if currentToken.Consume(TKReserved, '+') {
			proceedToken()

			right, err := mul()
			if err != nil {
				return nil, err
			}

			node = NewNode(NDAdd, node, right)

			continue
		}

		if currentToken.Consume(TKReserved, '-') {
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
		if currentToken.Consume(TKReserved, '*') {
			proceedToken()

			right, err := unary()
			if err != nil {
				return nil, err
			}

			node = NewNode(NDMul, node, right)

			continue
		}

		if currentToken.Consume(TKReserved, '/') {
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
	if currentToken.Consume(TKReserved, '+') {
		proceedToken()

		return unary()
	}

	if currentToken.Consume(TKReserved, '-') {
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
	if currentToken.Consume(TKReserved, '(') {
		proceedToken()

		node, err := expr()
		if err != nil {
			return nil, err
		}

		if err := currentToken.Expect(TKReserved, ')'); err != nil {
			return nil, err
		}

		proceedToken()

		return node, nil
	}

	if currentToken.ConsumeIdent() {
		return ident()
	}

	n, err := currentToken.ExpectNum()
	if err != nil {
		return nil, err
	}

	proceedToken()

	return NewNodeNum(n), nil
}

func ident() (*Node, error) {
	if currentToken.Skip().Consume(TKReserved, '(') {
		return identFuncCall()
	}

	node, err := identVal()

	proceedToken()

	return node, err
}

func identVal() (*Node, error) {
	if lv := localValue.FindValue(currentToken.Str); lv != nil {
		return lv, nil
	}

	node := NewNodeLocalValue(currentToken.Str, localValue.Offset+offsetSize)

	localValue.Next = node
	node.Root = localValue.Root

	localValue = node

	return node, nil
}

func identFuncCall() (*Node, error) {
	funcName := currentToken.Str

	proceedToken()

	if err := currentToken.Expect(TKReserved, '('); err != nil {
		return nil, err
	}

	proceedToken()

	args, err := funcCallArgs()
	if err != nil {
		return nil, err
	}

	return NewNodeFuncCall(funcName, args), nil
}

func funcCallArgs() (*Node, error) {
	if currentToken.Consume(TKReserved, ')') {
		proceedToken()

		return nil, nil
	}

	head, err := assign()
	if err != nil {
		return nil, err
	}

	cur := head

	for currentToken.Consume(TKReserved, ',') {
		proceedToken()

		node, err := assign()
		if err != nil {
			return nil, err
		}

		cur.Next = node

		cur = cur.Next
	}

	if err := currentToken.Expect(TKReserved, ')'); err != nil {
		return nil, err
	}

	proceedToken()

	return head, nil
}
