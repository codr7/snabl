package snabl

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
)

const (
	VERSION = 4
)

type Pc = int
type Tag = int

type Vm struct {
	Debug bool
	Trace bool
	
	Stdin io.Reader
	Stdout io.Writer

	AbcLib AbcLib
	TimeLib TimeLib
	
	Tags []V
	
	Code []Op

	Env Env
	Stack *Stack
	
	path string
	fun *Fun
	call *Call
}

func (self *Vm) Init() {
	self.Stdin = os.Stdin
	self.Stdout = os.Stdout
	self.AbcLib.Init(self)
	self.TimeLib.Init(self)
	self.Env = NewEnv(nil)
	self.Stack = new(Stack)
}

func (self *Vm) Path(in string) string {
	if filepath.IsAbs(in) {
		return in
	} 

	return filepath.Join(self.path, in)
}

func (self *Vm) LoadForms(path string, out *Forms) error {
	p := self.Path(path)
	f, err := os.Open(p)

	if err != nil {
		return err
	}

	defer f.Close()
	pos := NewPos(p, 1, 1)
	
	if err := self.ReadForms(pos, bufio.NewReader(f), out); err != nil {
		return err
	}

	return nil
}

func (self *Vm) Load(path string, eval bool) error {
	var forms Forms

	if err := self.LoadForms(path, &forms); err != nil {
		return err
	}
	
	pc := self.EmitPc()

	if err := forms.Emit(self, self.Env); err != nil {
		return err
	}
	
	self.Fuse(pc)

	if !eval {
		return nil
	}

	self.Code[self.Emit()] = StopOp()
	prevPath := self.path

	defer func() {
		self.path = prevPath
	}()
	
	self.path = filepath.Dir(self.Path(path))
	
	if err := self.Eval(&pc); err != nil {
		return err
	}

	return nil
}

func (self *Vm) Tag(t Type, d any) Tag {
	i := len(self.Tags)
	self.Tags = append(self.Tags, V{t: t, d: d})
	return i
}

func (self *Vm) E(pos *Pos, spec string, args...interface{}) *E {
	err := NewE(pos, spec, args...)
	
	if self.Debug {
		panic(err.Error())
	}

	return err
}

func (self *Vm) EmitNoTrace() Pc {
	pc := self.EmitPc()
	self.Code = append(self.Code, 0)
	return pc
}

func (self *Vm) Emit() Pc {
	if self.Trace {
		self.Code[self.EmitNoTrace()] = TraceOp()
	}
	
	return self.EmitNoTrace()
}

func (self *Vm) EmitPos(pos Pos) {
	tag := self.Tag(&self.AbcLib.PosType, pos)
	self.Code[self.EmitNoTrace()] = PosOp(tag)
}


func (self *Vm) EmitTag(t Type, d any) {
	self.Code[self.Emit()] = PushOp(self.Tag(t, d))
}

func (self *Vm) EmitPc() Pc {
	return Pc(len(self.Code))
}
