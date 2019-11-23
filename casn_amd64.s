#include "textflag.h"

TEXT ·cas(SB), NOSPLIT, $0
	MOVQ		ptr+0(FP), BX
	MOVQ		old+8(FP), AX
	MOVQ		new+16(FP), CX
	LOCK
	CMPXCHGQ 	CX, 0(BX)
	MOVQ 		AX, ret+24(FP)
	RET
