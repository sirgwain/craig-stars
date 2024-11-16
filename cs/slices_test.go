package cs

import (
	"reflect"
	"testing"
)

func TestAppendWithoutDuplicates(t *testing.T) {
	type args struct {
		slice  []string
		slice2 []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"same slice", args{[]string{"1", "0", "100", "a"}, []string{"1", "0", "100", "a"}}, []string{"1", "0", "100", "a"}},
		{"no duplicates", args{[]string{"1", "0", "100", "a"}, []string{"aaa"}}, []string{"1", "0", "100", "a", "aaa"}},
		{"duplicates removed from 2nd list", args{[]string{"1", "0", "100", "a"}, []string{"a", "1", "0", "a", "3"}}, []string{"1", "0", "100", "a", "3"}},
		{"duplicates removed from 1st list", args{[]string{"llllll"}, []string{"llllll", "11", "11", "11", "11", "llllll"}}, []string{"lllll", "11"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AppendWithoutDuplicates(tt.args.slice, tt.args.slice2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppendWithoutDuplicates() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompareSlicesUnordered(t *testing.T) {
	type args struct {
		slice     []string
		other     []string
		identical bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"identical with flag on", args{[]string{"1", "0", "100", "a"}, []string{"1", "0", "a", "100"}, true}, true},
		{"identical tech lists", args{[]string{"Tritanium", "Crobmnium", "Carbonic Armor", "Strobnium", "Organic Armor", "Kelarium", "Neutronium", "Valanium", "Superlatanium",
			"Mole-skin Shield", "Cow-hide Shield", "Wolverine Diffuse Shield", "Croby Sharmor", "Shadow Shield", "Bear Neutrino Barrier", "Langston Shell", "Gorilla Delagator", "Elephant Hide Fortress", "Complete Phase Shield"}, []string{"Tritanium", "Crobmnium", "Carbonic Armor", "Strobnium", "Organic Armor", "Kelarium", "Neutronium", "Valanium", "Superlatanium",
			"Mole-skin Shield", "Cow-hide Shield", "Wolverine Diffuse Shield", "Croby Sharmor", "Shadow Shield", "Bear Neutrino Barrier", "Langston Shell", "Gorilla Delagator", "Elephant Hide Fortress", "Complete Phase Shield"}, true}, true},
		{"buncha techs in different order; flag on", args{[]string{"Mole-skin Shield", "Cow-hide Shield", "Wolverine Diffuse Shield", "Croby Sharmor", "Shadow Shield", "Bear Neutrino Barrier", "Langston Shell", "Gorilla Delagator", "Elephant Hide Fortress", "Complete Phase Shield"},
			[]string{"Mole-skin Shield", "Langston Shell", "Wolverine Diffuse Shield", "Croby Sharmor", "Shadow Shield", "Bear Neutrino Barrier", "Elephant Hide Fortress", "Gorilla Delagator", "Complete Phase Shield", "Cow-hide Shield"}, true}, true},
		{"flag on & proper (non-identical) subset", args{[]string{"1", "100", "0", "a"}, []string{"1", "a", "100"}, true}, false},
		{"flag off & proper subset", args{[]string{"0", "0", "100", "a"}, []string{"0", "0"}, false}, true},
		{"flag off & proper subset in wrong order", args{[]string{"0", "99", "0", "100", "a"}, []string{"a", "99", "0", "b", "0", "100"}, false}, false},
		{"different", args{[]string{"0", "99", "0", "100", "a"}, []string{"banana", "0", "100"}, false}, false},
		{"multiple identical values", args{[]string{"3", "3", "3", "3", "3"}, []string{"3", "3", "3"}, false}, true},
		{"empty sets", args{[]string{}, []string{}, true}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CompareSlicesUnordered(tt.args.slice, tt.args.other, tt.args.identical); got != tt.want {
				t.Errorf("CompareSlicesUnordered() = %v, want %v", got, tt.want)
			}
		})
	}
}
