package main

import "fmt"

func prequel() {
	fmt.Println(".intel_syntax noprefix")
	fmt.Println(".globl main")
	fmt.Println("main:")
}

func sequel() {
	fmt.Println("  pop rax")
	fmt.Println("  ret")
}

func generate(node *Node) {
	if node.Kind == NDNum {
		fmt.Printf("  push %d\n", node.Val)

		return
	}

	generate(node.Left)
	generate(node.Right)

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
	case NDNum:
	}

	fmt.Println("  push rax")
}
