package snabl

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
)

const (
	VERSION = 3
)

type Pc = int
type Tag = int

type Vm struct {
	Debug bool
	Trace bool
	
	Path string
	Stdin io.Reader
	Stdout io.Writer

	AbcLib AbcLib
	
	Tags []V
	
	Code []Op
	Stack Slice[V]
	Calls Slice[Call]
	
	env Env
	fun *Fun
}

func (self *Vm) Init() {
	self.Stdin = os.Stdin
	self.Stdout = os.Stdout
	self.AbcLib.Init(self)
	self.env = NewEnv()
}

func (self *Vm) Load(path string, eval bool) error {
	var p string

	if filepath.IsAbs(path) {
		p = path
	} else {
		p = filepath.Join(self.Path, path)
	}
	
	f, err := os.Open(p)

	if err != nil {
		return err
	}

	defer f.Close()
	var forms Forms

	pos := NewPos(p, 1, 1)
	
	if err := self.ReadForms(pos, bufio.NewReader(f), &forms); err != nil {
		return err
	}
	
	pc := self.EmitPc()

	if err := forms.Emit(self, self.env); err != nil {
		return err
	}
	
	if !eval {
		return nil
	}

	self.Code[self.Emit()] = StopOp()
	prevPath := self.Path

	defer func() {
		self.Path = prevPath
	}()
	
	self.Path = filepath.Dir(p)
	self.Fuse(pc)
	
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

func (self *Vm) Env() Env {
	return self.env
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

func (self *Vm) EmitString(str string) {
	tag := self.Tag(&self.AbcLib.StringType, str)
	self.Code[self.Emit()] = PushOp(tag) 
}

func (self *Vm) EmitVal(t Type, d any) {
	tag := self.Tag(t, d)
	self.Code[self.Emit()] = PushOp(tag)
}

func (self *Vm) EmitPc() Pc {
	return Pc(len(self.Code))
}
