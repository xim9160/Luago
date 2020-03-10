package vm

type Instruction uint32

const MAXARG_Bx	 = 1<<18 - 1 // 2^18 - 1 = 262143
const MAXARG_sBx = MAXARG_Bx >> 1 // 262143 / 2 = 131071 

// 从指令中获取操作码
func (self Instruction) Opcode() int  {
	return int(self & 0x3F)
}

// ABC 从 iABC 模式中提取参数
// mark: 重要 因为 iABC 内存分布是 按照 B:9 C:9 A:8 opcode:6 分布, 所以这里容易出错
func (self Instruction) ABC() (a, b, c int)  {
	a = int(self >> 6 & 0xFF)
	c = int(self >> 14 & 0x1FF)
	b = int(self >> 23 & 0x1FF)
	return
}

// ABx 从 iABx 中提取参数
func (self Instruction) ABx() (a, bx int)  {
	a = int(self >> 6 & 0xFF)
	bx = int(self >> 14)
	return
}

// AsBx 从 iAsBx 中提取参数
func (self Instruction) AsBx() (a, sbx int)  {
	a, bx := self.ABx()
	return a, bx - MAXARG_sBx
}

// Ax 方法从 iAx 模式照片那个提取参数
func (self Instruction) Ax() int  {
	return int(self >> 6)
}

func (self Instruction) OpName() string  {
	return opcodes[self.Opcode()].name
}

func (self Instruction) OpMode() byte  {
	return opcodes[self.Opcode()].opMode
}

func (self Instruction) BMode() byte  {
	return opcodes[self.Opcode()].argBMode
}

func (self Instruction) CMode() byte {
	return opcodes[self.Opcode()].argCMode
}