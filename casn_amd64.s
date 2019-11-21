#include "textflag.h"


TEXT ·Cas(SB),NOSPLIT,$0
	MOVQ		ptr+0(FP), BP
	MOVQ		old+8(FP), AX
	MOVQ		new+16(FP), CX
	LOCK
	CMPXCHGQ 	CX, (BP)
	MOVQ 		AX, ret+24(FP)
	SETEQ		ret+32(FP)
	RET
