package main

import "fmt"

func prequel() {
	fmt.Println(".intel_syntax noprefix")
	fmt.Println(".globl main")
	fmt.Println("main:")

	// 変数領域を確保
	fmt.Println("  push rbp")
	fmt.Println("  mov rbp, rsp")
	fmt.Printf("  sub rsp, %d\n", offsetSize*26)
}

var argreg = []string{"rdi", "rsi", "rdx", "rcx", "r8", "r9"}

func generate(node *Node) error {
	switch node.Kind {
	case NDNum:
		fmt.Printf("  push %d\n", node.Val)

		return nil
	case NDLocalV:
		if err := generateLocalValue(node); err != nil {
			return err
		}

		fmt.Println("  pop rax")
		fmt.Println("  mov rax, [rax]")
		fmt.Println("  push rax")

		return nil
	case NDAssign:
		if err := generateLocalValue(node.Left); err != nil {
			return err
		}

		if err := generate(node.Right); err != nil {
			return err
		}

		fmt.Println("  pop rdi")
		fmt.Println("  pop rax")
		fmt.Println("  mov [rax], rdi")
		fmt.Println("  push rdi")

		return nil
	case NDReturn:
		if err := generate(node.Left); err != nil {
			return err
		}

		fmt.Println("  pop rax")
		fmt.Println("  mov rsp, rbp")
		fmt.Println("  pop rbp")
		fmt.Println("  ret")

		return nil
	case NDIf:
		if err := generate(node.Cond); err != nil {
			return err
		}

		label := uniqueLabel()

		fmt.Println("  pop rax")
		fmt.Println("  cmp rax, 0")

		if node.Else != nil {
			fmt.Printf("  je  .L.else.%s\n", label)

			if err := generate(node.Then); err != nil {
				return err
			}

			fmt.Printf("  jmp .L.end.%s\n", label)
			fmt.Printf(".L.else.%s:\n", label)

			if err := generate(node.Else); err != nil {
				return err
			}

			fmt.Printf(".L.end.%s:\n", label)
		} else {
			fmt.Printf("  je  .L.end.%s\n", label)

			if err := generate(node.Then); err != nil {
				return err
			}

			fmt.Printf(".L.end.%s:\n", label)
		}

		return nil
	case NDFor:
		label := uniqueLabel()

		if node.Init != nil {
			if err := generate(node.Init); err != nil {
				return err
			}
		}

		fmt.Printf(".L.begin.%s:\n", label)

		if node.Cond != nil {
			if err := generate(node.Cond); err != nil {
				return err
			}

			fmt.Println("pop rax")
			fmt.Println("cmp rax, 0")
			fmt.Printf("je .L.end.%s\n", label)
		}

		if err := generate(node.Then); err != nil {
			return err
		}

		if node.Inc != nil {
			if err := generate(node.Inc); err != nil {
				return err
			}
		}

		fmt.Printf("jmp .L.begin.%s\n", label)
		fmt.Printf(".L.end.%s:\n", label)

		return nil

	case NDBlock:
		for cur := node.Body; cur != nil; cur = cur.Next {
			if err := generate(cur); err != nil {
				return err
			}
		}

		return nil
	case NDFunc:
		var nargs int

		for arg := node.Args; arg != nil; arg = arg.Next {
			if err := generate(arg); err != nil {
				return err
			}
			nargs++
		}

		for i := nargs - 1; i >= 0; i-- {
			fmt.Printf("  pop %s\n", argreg[i])
		}

		label := uniqueLabel()

		fmt.Printf("  mov rax, rsp\n")
		fmt.Printf("  and rax, 15\n")
		fmt.Printf("  jnz .L.call.%s\n", label)
		fmt.Printf("  mov rax, 0\n")
		fmt.Printf("  call %s\n", string(node.FuncName))
		fmt.Printf("  jmp .L.end.%s\n", label)
		fmt.Printf(".L.call.%s:\n", label)
		fmt.Printf("  sub rsp, 8\n")
		fmt.Printf("  mov rax, 0\n")
		fmt.Printf("  call %s\n", string(node.FuncName))
		fmt.Printf("  add rsp, 8\n")
		fmt.Printf(".L.end.%s:\n", label)
		fmt.Printf("  push rax\n")

		return nil
	}

	if err := generate(node.Left); err != nil {
		return err
	}

	if err := generate(node.Right); err != nil {
		return err
	}

	fmt.Println("  pop rdi")
	fmt.Println("  pop rax")

	switch node.Kind {
	case NDAdd:
		fmt.Println("  add rax, rdi")
	case NDSub:
		fmt.Println("  sub rax, rdi")
	case NDMul:
		fmt.Println("  imul rax, rdi")
	case NDDiv:
		fmt.Println("  cqo")
		fmt.Println("  idiv rdi")
	case NDEq:
		fmt.Println("  cmp rax, rdi")
		fmt.Println("  sete al")
		fmt.Println("  movzb rax, al")
	case NDNe:
		fmt.Println("  cmp rax, rdi")
		fmt.Println("  setne al")
		fmt.Println("  movzb rax, al")
	case NDLt:
		fmt.Println("  cmp rax, rdi")
		fmt.Println("  setl al")
		fmt.Println("  movzb rax, al")
	case NDLe:
		fmt.Println("  cmp rax, rdi")
		fmt.Println("  setle al")
		fmt.Println("  movzb rax, al")
	}

	fmt.Println("  push rax")

	return nil
}

func generateLocalValue(node *Node) error {
	if node.Kind != NDLocalV {
		return userInput.Err(token.Loc, "変数ではありません")
	}

	fmt.Println("  mov rax, rbp")
	fmt.Printf("  sub rax, %d\n", node.Offset)
	fmt.Println("  push rax")

	return nil
}

func sequel() {
	fmt.Println("  mov rsp, rbp")
	fmt.Println("  pop rbp")
	fmt.Println("  ret")
}

func uniqueLabel() string {
	one := label % 26
	two := label / 26 % 26
	three := label / (26 * 26) % (26 * 26)

	label++

	return fmt.Sprintf("%c%c%c", 'A'+three, 'A'+two, 'A'+one)
}
