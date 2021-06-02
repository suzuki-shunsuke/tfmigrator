package expr_test

import (
	"testing"

	"github.com/suzuki-shunsuke/tfmigrator/pkg/expr"
	"gopkg.in/yaml.v2"
)

func TestBool_UnmarshalYAML(t *testing.T) {
	t.Parallel()
	data := []struct {
		title string
		yaml  string
		param interface{}
		exp   bool
	}{
		{
			title: "normal",
			yaml:  `name == "foo"`,
			param: map[string]interface{}{
				"name": "foo",
			},
			exp: true,
		},
	}
	for _, d := range data {
		d := d
		t.Run(d.title, func(t *testing.T) {
			t.Parallel()
			b := expr.Bool{}
			if err := yaml.Unmarshal([]byte(d.yaml), &b); err != nil {
				t.Fatal(err)
			}
			if b.Empty() {
				t.Fatal("bool is empty")
			}
			f, err := b.Run(d.param)
			if err != nil {
				t.Fatal(err)
			}
			if f && !d.exp {
				t.Fatal(`got true, wanted false`)
			}
			if !f && d.exp {
				t.Fatal(`got false, wanted true`)
			}
		})
	}
}

func TestBool_Empty(t *testing.T) {
	t.Parallel()
	b := expr.Bool{}
	if !b.Empty() {
		t.Fatal("Bool.Empty() should be true")
	}
}

func TestNewBool(t *testing.T) {
	t.Parallel()
	b, err := expr.NewBool("false")
	if err != nil {
		t.Fatal(err)
	}
	f, err := b.Run(nil)
	if err != nil {
		t.Fatal(err)
	}
	if f {
		t.Fatal("Bool must be false")
	}
}
