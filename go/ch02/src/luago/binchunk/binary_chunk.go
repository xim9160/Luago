package binchunk

type binchunk struct {
	header				// 头部
	sizeUpvales byte 	// 主函数 upvalue 数量
	mainfunc *Prototype // 主函数原型
}

type header struct {
	signature 		[4]byte 	//Magic Number 快速识别文件格式
	version 		byte 		//lua版本号 Major Minor Release 值为 Major * 16 + Minor 
	format			byte 		//格式号, 和预期不符合的话 会拒绝加载
	luacData		[6]byte 	//进一步校验
	cintSize		byte 		//cint 在 chunk 里占用的字节数
	sizeSize		byte		//size_t 在 chunk 里占用的字节数
	instructionSize	byte		//lua 虚拟机指令 在 chunk 里占用的字节数
	luaIntegerSize	byte		//lua 整数 在 chunk 里占用的字节数
	luaNumberSize	byte		//lua 浮点数 在 chunk 里占用的字节数
	luacInt			byte 		//接下来n个字节存放 Lua 整数值 0x5678 目的是检验大小端
	luacNum			byte		// 头部的后n个字节, 存放 lua 浮点数 370.5
}

const (
	LUA_SIGNATURE 		= "\x1bLua"
	LUAC_VERSION		= 0x53
	LUAC_FORMAT			= 0
	LUAC_DATA 			= "\x19\x93\r\n\x1a\n"
	CINT_SIZE			= 4
	CSIZET_SIZE			= 8
	INSTRUCTION_SIZE	= 4
	LUA_INTEGET_SIZE	= 8
	LUA_NUMBER_SIZE		= 8
	LUAC_INT			= 0x5678
	LUAC_NUM			= 370.5	
)

type Prototype struct {
	Source			string			// 存放 源文件名, 只有主函数原型里面有值
	LineDefined 	uint32			// 起始行号
	LastLineDefined	uint32			// 终止行号
	NumParams		byte			// 固定参数个数
	IsVararg		byte			// 是否有变长参数
	MaxStackSize	byte			// 寄存器数量
	Code			[]uint32		// 指令表
	Constants		[]interface{}	// 常量表 存放字面量 nil, bool, 整数, 浮点数, 字符串等 5 个
	Upvalues		[]Upvalue		//
	Protos			[]*Prototype	// 子函数原型
	LineInfo		[]uint32		// 行号表, 记录每条指令在源码中对应的行号
	LocVars			[]LocVar		// 局部变量表, 记录局部变量名
	UpvalueNames	[]string		// Upvalue 名列表 比如 _ENV
}

const (
	TAG_NIL 		= 0x00
	TAG_BOOLEN		= 0x01
	TAG_NUMBER		= 0x03
	TAG_INTEGER		= 0x13
	TAG_SHORT_STR	= 0x04
	TAG_LONG_STR	= 0x14
)

type Upvalue struct {
	Instack byte
	Idx		byte
}

type LocVar struct {
	VarName string
	StartPC uint32
	EndPC	uint32
}

func Undump(data []byte) *Prototype {
	reader := &reader{data}
	reader.checkHeader()		//头部校验
	reader.readByte()			//跳过校验 Upvalue 数量
	return reader.readProto("") //读取函数原型
}