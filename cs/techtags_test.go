//go:build !wasi && !wasm

package cs

import "testing"

func TestTechTags_hasTags(t *testing.T) {
	type args struct {
		tagsToInclude []TechTag
		tagsToExclude []TechTag
	}
	tests := []struct {
		name string
		tt   TechTags
		args args
		want bool
	}{
		{"no tags", TechTags{}, args{tagsToInclude: nil, tagsToExclude: nil}, false},
		{"has tag", newTechTags(TechTagArmor), args{tagsToInclude: []TechTag{TechTagArmor}, tagsToExclude: nil}, true},
		{"no shield tag", newTechTags(TechTagArmor), args{tagsToInclude: []TechTag{TechTagShield}, tagsToExclude: nil}, false},
		{"shield blacklisted", newTechTags(TechTagArmor, TechTagShield), args{tagsToInclude: []TechTag{}, tagsToExclude: []TechTag{TechTagShield}}, false},
		{"shield whitelisted", newTechTags(TechTagArmor, TechTagShield), args{tagsToInclude: []TechTag{TechTagShield}, tagsToExclude: []TechTag{TechTagShield}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tt.hasTags(tt.args.tagsToInclude, tt.args.tagsToExclude...); got != tt.want {
				t.Errorf("TechTags.hasTags() = %v, want %v", got, tt.want)
			}
		})
	}
}
