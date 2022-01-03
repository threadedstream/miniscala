
	.text
	.globl miniscala_main
	miniscala_main:
	
	movq $20, %rax
	imulq $30, %rax
	
	ret
	