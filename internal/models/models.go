package models

// Desc = Description

type Settings struct {
	ProjectName         string
	ProjectDesc         string
	ProjectPath         string
	DocGenPath          string
	IncludeTests        bool
	IncludePrivateFuncs bool
	IncludePrivateVars  bool
	ExcludePackages     []string
}

type Package struct {
	Name  string
	Desc  string
	Usage string
	Files []File
	Types []Type
	Vars  []Var
	Funcs []Func
	Deps  []Dependency
}

type Dependency struct {
	Name string
	Desc string
}

type File struct {
	Path    string
	Name    string
	Desc    string
	Author  string
	Version string
	Date    string
	Funcs   []Func
	Vars    []Var
	Types   []Type
}

type Type struct {
	Name   string
	Desc   string
	Fields []Var
}

type Var struct {
	Name string
	Type string
	Desc string
}

type Func struct {
	Name     string
	Desc     string
	Params   []Var
	Returns  []Return
	Receiver string
}

type Return struct {
	Type string
	Desc string
}
