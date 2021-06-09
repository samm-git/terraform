package globalref

import (
	"sort"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform/internal/addrs"
)

func TestAnalyzerMetaReferences(t *testing.T) {
	tests := []struct {
		InputModule string
		InputRef    string
		WantModule  string
		WantRefs    []string
	}{
		{
			``,
			`local.a`,
			``,
			nil,
		},
		{
			``,
			`test_thing.single`,
			``,
			[]string{
				"local.a",
				"local.b",
			},
		},
		{
			``,
			`test_thing.single.string`,
			``,
			[]string{
				"local.a",
			},
		},
		{
			``,
			`test_thing.for_each`,
			``,
			[]string{
				"local.a",
				"test_thing.single.string",
			},
		},
		{
			``,
			`test_thing.for_each["whatever"]`,
			``,
			[]string{
				"local.a",
				"test_thing.single.string",
			},
		},
		{
			``,
			`test_thing.for_each["whatever"].single`,
			``,
			[]string{
				"test_thing.single.string",
			},
		},
		{
			``,
			`test_thing.for_each["whatever"].single.z`,
			``,
			[]string{
				"test_thing.single.string",
			},
		},
		{
			``,
			`test_thing.count`,
			``,
			[]string{
				"local.a",
			},
		},
		{
			``,
			`test_thing.count[0]`,
			``,
			[]string{
				"local.a",
			},
		},
		{
			``,
			`module.single.a`,
			`module.single`,
			[]string{
				"test_thing.foo",
				"var.a",
			},
		},
		{
			``,
			`module.for_each["whatever"].a`,
			`module.for_each["whatever"]`,
			[]string{
				"test_thing.foo",
				"var.a",
			},
		},
		{
			``,
			`module.count[0].a`,
			`module.count[0]`,
			[]string{
				"test_thing.foo",
				"var.a",
			},
		},
		{
			`module.single`,
			`var.a`,
			``,
			[]string{
				"test_thing.single",
			},
		},
		{
			`module.single`,
			`test_thing.foo`,
			`module.single`,
			[]string{
				"var.a",
			},
		},
	}

	azr := testAnalyzer(t, "assorted")

	for _, test := range tests {
		name := test.InputRef
		if test.InputModule != "" {
			name = test.InputModule + " " + test.InputRef
		}
		t.Run(name, func(t *testing.T) {
			t.Logf("testing %s", name)
			moduleAddr := addrs.RootModuleInstance
			if test.InputModule != "" {
				moduleAddrTarget, diags := addrs.ParseTargetStr(test.InputModule)
				if diags.HasErrors() {
					t.Fatalf("input module address is invalid: %s", diags.Err())
				}
				var ok bool
				moduleAddr, ok = moduleAddrTarget.Subject.(addrs.ModuleInstance)
				if !ok {
					t.Fatalf("input module address is invalid: must be %T, not %T", moduleAddr, moduleAddrTarget.Subject)
				}
			}

			ref, diags := addrs.ParseRefStr(test.InputRef)
			if diags.HasErrors() {
				t.Fatalf("input reference is invalid: %s", diags.Err())
			}

			modAddr, refs := azr.MetaReferences(moduleAddr, ref)
			if got, want := modAddr.String(), test.WantModule; got != want {
				t.Errorf("wrong module address\ngot:  %s\nwant: %s", got, want)
			}

			want := test.WantRefs
			var got []string
			for _, ref := range refs {
				got = append(got, ref.DisplayString())
			}
			sort.Strings(got)
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("wrong references\n%s", diff)
			}
		})
	}
}
