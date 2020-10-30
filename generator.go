package main

import "fmt"

var argReg = []string{"rdi", "rsi", "rdx", "rcx", "r8", "r9"}

func generate(nodes *Node) error {
	output.L(".intel_syntax noprefix")

	for node := nodes; node != nil; node = node.Next {
		if node.Kind == NDFuncDef {
			if err := genFunction(node); err != nil {
				return err
			}
		}
	}

	return nil
}

func genFunction(node *Node) error {
	funcName := string(node.Name)
	output.F(".global %s\n", funcName)
	output.F("%s:\n", funcName)

	output.L("  push rbp")
	output.L("  mov rbp, rsp")
	output.F("  sub rsp, %d\n", node.Size)

	var i int

	for p := node.Params; p != nil; p = p.Next {
		output.F("  mov [rbp-%d], %s\n", p.Offset, argReg[i])
		i++
	}

	for body := node.Body; body != nil; body = body.Next {
		if err := genStmt(body); err != nil {
			return err
		}
	}

	output.F(".L.return.%s:\n", funcName)
	output.L("  mov rsp, rbp")
	output.L("  pop rbp")
	output.L("  ret")

	return nil
}

func genStmt(node *Node) error {
	switch node.Kind {
	case NDNum:
		output.F("  push %d\n", node.Val)

		return nil
	case NDLocalV:
		if err := generateLocalValue(node); err != nil {
			return err
		}

		output.L("  pop rax")
		output.L("  mov rax, [rax]")
		output.L("  push rax")

		return nil
	case NDAssign:
		if err := generateLocalValue(node.Left); err != nil {
			return err
		}

		if err := genStmt(node.Right); err != nil {
			return err
		}

		output.L("  pop rdi")
		output.L("  pop rax")
		output.L("  mov [rax], rdi")
		output.L("  push rdi")

		return nil
	case NDReturn:
		if err := genStmt(node.Left); err != nil {
			return err
		}

		output.L("  pop rax")
		output.L("  mov rsp, rbp")
		output.L("  pop rbp")
		output.L("  ret")

		return nil
	case NDIf:
		if err := genStmt(node.Cond); err != nil {
			return err
		}

		label := uniqueLabel()

		output.L("  pop rax")
		output.L("  cmp rax, 0")

		if node.Else != nil {
			output.F("  je  .L.else.%s\n", label)

			if err := genStmt(node.Then); err != nil {
				return err
			}

			output.F("  jmp .L.end.%s\n", label)
			output.F(".L.else.%s:\n", label)

			if err := genStmt(node.Else); err != nil {
				return err
			}

			output.F(".L.end.%s:\n", label)
		} else {
			output.F("  je  .L.end.%s\n", label)

			if err := genStmt(node.Then); err != nil {
				return err
			}

			output.F(".L.end.%s:\n", label)
		}

		return nil
	case NDFor:
		label := uniqueLabel()

		if node.Init != nil {
			if err := genStmt(node.Init); err != nil {
				return err
			}
		}

		output.F(".L.begin.%s:\n", label)

		if node.Cond != nil {
			if err := genStmt(node.Cond); err != nil {
				return err
			}

			output.L("pop rax")
			output.L("cmp rax, 0")
			output.F("je .L.end.%s\n", label)
		}

		if err := genStmt(node.Then); err != nil {
			return err
		}

		if node.Inc != nil {
			if err := genStmt(node.Inc); err != nil {
				return err
			}
		}

		output.F("jmp .L.begin.%s\n", label)
		output.F(".L.end.%s:\n", label)

		return nil
	case NDBlock:
		for cur := node.Body; cur != nil; cur = cur.Next {
			if err := genStmt(cur); err != nil {
				return err
			}
		}

		return nil
	case NDFuncCall:
		var nargs int

		for arg := node.Args; arg != nil; arg = arg.Next {
			if err := genStmt(arg); err != nil {
				return err
			}
			nargs++
		}

		for i := nargs - 1; i >= 0; i-- {
			output.F("  pop %s\n", argReg[i])
		}

		label := uniqueLabel()

		output.F("  mov rax, rsp\n")
		output.F("  and rax, 15\n")
		output.F("  jnz .L.call.%s\n", label)
		output.F("  mov rax, 0\n")
		output.F("  call %s\n", string(node.Name))
		output.F("  jmp .L.end.%s\n", label)
		output.F(".L.call.%s:\n", label)
		output.F("  sub rsp, 8\n")
		output.F("  mov rax, 0\n")
		output.F("  call %s\n", string(node.Name))
		output.F("  add rsp, 8\n")
		output.F(".L.end.%s:\n", label)
		output.F("  push rax\n")

		return nil
	}

	if err := genStmt(node.Left); err != nil {
		return err
	}

	if err := genStmt(node.Right); err != nil {
		return err
	}

	output.L("  pop rdi")
	output.L("  pop rax")

	switch node.Kind {
	case NDAdd:
		output.L("  add rax, rdi")
	case NDSub:
		output.L("  sub rax, rdi")
	case NDMul:
		output.L("  imul rax, rdi")
	case NDDiv:
		output.L("  cqo")
		output.L("  idiv rdi")
	case NDEq:
		output.L("  cmp rax, rdi")
		output.L("  sete al")
		output.L("  movzb rax, al")
	case NDNe:
		output.L("  cmp rax, rdi")
		output.L("  setne al")
		output.L("  movzb rax, al")
	case NDLt:
		output.L("  cmp rax, rdi")
		output.L("  setl al")
		output.L("  movzb rax, al")
	case NDLe:
		output.L("  cmp rax, rdi")
		output.L("  setle al")
		output.L("  movzb rax, al")
	}

	output.L("  push rax")

	return nil
}

func generateLocalValue(node *Node) error {
	if node.Kind != NDLocalV {
		return userInput.Err(currentToken.Loc, "変数ではありません")
	}

	output.L("  mov rax, rbp")
	output.F("  sub rax, %d\n", node.Offset)
	output.L("  push rax")

	return nil
}

func uniqueLabel() string {
	one := label % 26
	two := label / 26 % 26
	three := label / (26 * 26) % (26 * 26)

	label++

	return fmt.Sprintf("%c%c%c", 'A'+three, 'A'+two, 'A'+one)
}
