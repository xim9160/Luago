package main

import "fmt"
import . "luago/api"
import "luago/state"

func main()  {
	ls := state.New()
	ls.PushBoolen(true); 		printStack(ls)
	ls.PushInteger(10); 		printStack(ls)
	ls.PushNil();				printStack(ls)
	ls.PushString("hello");		printStack(ls)
	ls.PushValue(-4);			printStack(ls)
	ls.Replace(3);				printStack(ls)
	ls.SetTop(6);				printStack(ls)
	ls.Remove(-3);				printStack(ls)
	ls.SetTop(-5);				printStack(ls)
}

func list(f *binchunk.Prototype) {
	printHeader(f)
	printCode(f)
	printDetail(f)
	for _, p := range f.Protos {
		list(p)
	}
}

func printHeader(f *binchunk.Prototype) {
	funcType := "main"
	if f.LineDefined > 0 { funcType = "function" }

	varargFlag := ""
	if f.IsVararg > 0 { varargFlag = "+" }

	fmt.Printf("\n%s <%s:%d, %d> (%d instructions)\n", funcType, f.Source, f.LineDefined, f.LastLineDefined, len(f.Code))

	fmt.Printf("%d%s params, %d slots, %d upvalues, ", f.NumParams, varargFlag, f.MaxStackSize, len(f.Upvalues))

	fmt.Printf("%d locals, %d constants, %d functions\n", len(f.LocVars), len(f.Constants), len(f.Protos))
}

/* ch02 
func printCode(f *binchunk.Prototype) {
	for pc, c := range f.Code {
		line := "-"
		if len(f.LineInfo) > 0 {
			line = fmt.Sprint("%d", f.LineInfo[pc])
		}
		fmt.Printf("\t%d\t[%s]\t0x%08X\n", pc + 1, line, c)
	}
}*/

func printDetail(f *binchunk.Prototype) {
	fmt.Printf("constants (%d):\n", len(f.Constants))
	for i, k := range f.Constants {
		fmt.Printf("\t%d\t%s\n", i+1, constantToString(k))
	}

	fmt.Printf("locals (%d):\n", len(f.LocVars))
	for i, locVar := range f.LocVars {
		fmt.Printf("\t%d\t%s\t%d\t%d\n", i, locVar.VarName, locVar.StartPC + 1, locVar.EndPC + 1)
	}

	fmt.Printf("upvalues (%d):\n", len(f.Upvalues))
	for i, upval := range f.Upvalues {
		fmt.Printf("\t%d\t%s\t%d\t%d\n", i, upvalName(f, i), upval.Instack, upval.Idx)
	}
}

func constantToString(k interface{}) string {
	switch k.(type) {
	case nil:		return "nil"
	case bool:		return fmt.Sprintf("%t", k)
	case float64:	return fmt.Sprintf("%g", k)
	case int64:		return fmt.Sprintf("%d", k)
	case string:	return fmt.Sprintf("%q", k)
	default:		return "?"
	}
}

func upvalName(f *binchunk.Prototype, idx int) string {
	if len(f.UpvalueNames) > 0 {
		return f.UpvalueNames[idx]
	}
	return "-"
}

func printCode(f *binchunk.Prototype)  {
	for pc, c := range f.Code {
		line := "-"
		if len(f.LineInfo) > 0 {
			line = fmt.Sprintf("%d", f.LineInfo[pc])
		}
		i := Instruction(c)
		fmt.Printf("\t%d\t[%s]\t%s \t", pc+1, line, i.OpName())
		printOperands(i)
		fmt.Printf("\n")
	}
}

func printOperands(i Instruction)  {
	switch i.OpMode() {
	case IABC:
		a, b, c := i.ABC()
		fmt.Printf("%d", a) // A 是一定使用的, 但是 BC 可能不使用
		if i.BMode() != OpArgN {
			// 如果操作数最高位时1 就认为表示表索引, 按照负数输出
			if b > 0xFF {
				fmt.Printf(" %d", -1-b&0xFF)
			} else {
				fmt.Printf(" %d", b)
			}
		}
		if i.CMode() != OpArgN {
			if c > 0xFF {
				fmt.Printf(" %d", -1-c&0xFF)
			} else {
				fmt.Printf(" %d", c)
			}
		}
	case IABx:
		a, bx := i.ABx()
		fmt.Printf("%d", a)
		// 如果BMode = K 就认为表示表索引, 按照负数输出
		if i.BMode() == OpArgK {
			fmt.Printf(" %d", -1-bx)
		} else if i.BMode() == OpArgU {
			fmt.Printf(" %d", bx)
		}
	case IAsBx:
		a, sbx := i.AsBx()
		// 直接打印
		fmt.Printf(" %d %d", a, sbx)
	case IAx:
		ax := i.Ax()
		// 只有 a 直接打印
		fmt.Printf("%d", -1-ax)
	}
}

func printStack(ls LuaState)  {
	top := ls.GetTop()
	for i := 0; i < top; i++ {
		t := ls.Type(i)
		switch t {
		case LUA_TBOOLEN: 	fmt.Print("[%t]", ls.ToBoolean(i))
		case LUA_TNUMBER: 	fmt.Print("[%g]", ls.ToNumber(i))
		case LUA_TSTRING: 	fmt.Print("[%q]", ls.ToString(i))
		default:			fmt.Print("[%s]", ls.TypeName(t))
		}
		fmt.Println()
	}
}