package controller

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type ResourcePathComputer struct {
}

type CompiledResourcePathComputer struct {
	tmpl *template.Template
}

func (rpc *ResourcePathComputer) Compile(src string) (CompiledResourcePathComputer, error) {
	crpc := CompiledResourcePathComputer{}
	tmpl, err := template.New("_").Funcs(sprig.TxtFuncMap()).Parse(src)
	if err != nil {
		return crpc, fmt.Errorf("parse a template: %w", err)
	}
	crpc.tmpl = tmpl

	return crpc, nil
}

func (crpc *CompiledResourcePathComputer) Parse(rsc interface{}) (string, error) {
	buf := &bytes.Buffer{}
	if err := crpc.tmpl.Execute(buf, rsc); err != nil {
		return "", fmt.Errorf("render a template with params: %w", err)
	}
	p := buf.String()
	if !hclsyntax.ValidIdentifier(p) {
		return "", fmt.Errorf("invalid resource path: " + p)
	}
	return buf.String(), nil
}
