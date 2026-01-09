package service

import "testing"

func TestSlugify(t *testing.T) {

	tests := []struct {
		input    string
		expected string
	}{
		{"Hello World", "hello-world"},
		{"  Trim Spaces  ", "trim-spaces"},
		{"SpecialChars!@#", "specialchars"},
		{"Multiple   Spaces", "multiple-spaces"},
		{"CamelCaseString", "camelcasestring"},
		{"123 Numbers", "123-numbers"},
		{"--Dashes--", "dashes"},
	}

	for _, test := range tests {
		result := Slugify(test.input)
		if result != test.expected {
			t.Errorf("Slugify(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}
