package state

func (self *luaState) GetTop() int  {
	return self.stack.top
}

func (self *luaState) AbsIndex(idx int) int  {
	return self.stack.absIndex(idx)
}

//check 确保 可以插入 n个值而不会溢出

func (self *luaState) CheckStack(n int) bool  {
	self.stack.check(n)
	return true
}

func (self *luaState) Pop(n int)  {
	for i := 0; i < n; i++ {
		self.stack.pop()
	}
}

func (self *luaState) Copy(fromIdx, toIdx int)  {
	val := self.stack.get(fromIdx)
	self.stack.set(toIdx, val)
}

func (self *luaState) PushValue(idx int)  {
	val := self.stack.get(idx)
	self.stack.push(val)
}

func (self *luaState) Replace(idx int)  {
	val := self.stack.pop()
	self.stack.set(idx, val)
}

func (self *luaState) Insert(idx int)  {
	self.Rotate(idx, 1)
}

func (self *luaState) Remove(idx int) {
	self.Rotate(idx, -1)
	self.Pop(1)
}

//使用三次翻转进行栈方向旋转
func (self *luaState) Rotate(idx, n int)  {
	t := self.stack.top - 1
	p := self.stack.absIndex(idx) - 1
	var m int
	if n >= 0 {
		m = t - n
	} else {
		m = p - n - 1
	}
	self.stack.reverse(p, m)
	self.stack.reverse(m+1, t)
	self.stack.reverse(p, t)
}

func (self *luaStack) reverse(from, to int)  {
	slots := self.slots
	for from < to {
		slots[from], slots[to] = slots[to], slots[from]
		from++
		to--
	}
}

func (self *luaState) SetTop(idx int)  {
	newTop := self.stack.absIndex(idx)
	if newTop < 0 {
		panic("stack underflow!")
	}

	n := self.stack.top - newTop
	if n > 0 {
		for i := 0; i < n; i++ {
			self.stack.pop()
		}
	} else if n < 0 {
		for i := 0; i < n; i-- {
			self.stack.push(nil)
		}
	}
}

func (self *luaState) Pop(n int) {
	self.SetTop(-n-1)
}