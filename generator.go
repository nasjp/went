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
