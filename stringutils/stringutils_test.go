package stringutils

import (
	"testing"
)

func Test_Stripchars_strips_properly(t *testing.T) {
	actual := Stripchars("qwerty.asdfgh.", ".")
	const expected = "qwertyasdfgh"
	if actual != expected {
		t.Errorf("Expected: '%s', got '%s'", expected, actual)
	}
}

func Test_Stripchars_doesnt_strip_unspecified_characters(t *testing.T) {
	actual := Stripchars("abc", ".qwe")
	const expected = "abc"
	if actual != expected {
		t.Errorf("Expected: '%s', got '%s'", expected, actual)
	}
}

func Test_Stripchars_stripping_all_characters_gives_empty_string(t *testing.T) {
	actual := Stripchars("123456", "123456")
	const expected = ""
	if actual != expected {
		t.Errorf("Expected: '%s', got '%s'", expected, actual)
	}
}

func Test_LeftPad1(t *testing.T) {
	actual := LeftPad("aa", "0", 3)
	const expected = "000aa"
	if actual != expected {
		t.Errorf("Expected: '%s', got '%s'", expected, actual)
	}
}

func Test_LeftPad2(t *testing.T) {
	actual := LeftPad("abc", "g", 3)
	const expected = "gggabc"
	if actual != expected {
		t.Errorf("Expected: '%s', got '%s'", expected, actual)
	}
}

func Test_LeftPad3(t *testing.T) {
	actual := LeftPad("1", "0", 0)
	const expected = "1"
	if actual != expected {
		t.Errorf("Expected: '%s', got '%s'", expected, actual)
	}
}

func Test_RightPad1(t *testing.T) {
	actual := RightPad("abc", "g", 3)
	const expected = "abcggg"
	if actual != expected {
		t.Errorf("Expected: '%s', got '%s'", expected, actual)
	}
}

func Test_RightPad2(t *testing.T) {
	actual := RightPad("1", "0", 5)
	const expected = "100000"
	if actual != expected {
		t.Errorf("Expected: '%s', got '%s'", expected, actual)
	}
}

func Test_RightPad3(t *testing.T) {
	actual := RightPad("1", "0", 0)
	const expected = "1"
	if actual != expected {
		t.Errorf("Expected: '%s', got '%s'", expected, actual)
	}
}

func Test_LeftPad2Len1(t *testing.T) {
	actual := LeftPad2Len("1", "0", 10)
	const expected = "0000000001"
	if actual != expected {
		t.Errorf("Expected: '%s', got '%s'", expected, actual)
	}
}

func Test_LeftPad2Len2(t *testing.T) {
	actual := LeftPad2Len("1", "a", 5)
	const expected = "aaaa1"
	if actual != expected {
		t.Errorf("Expected: '%s', got '%s'", expected, actual)
	}
}

func Test_RightPad2Len1(t *testing.T) {
	actual := RightPad2Len("1", "0", 10)
	const expected = "1000000000"
	if actual != expected {
		t.Errorf("Expected: '%s', got '%s'", expected, actual)
	}
}

func Test_RightPad2Len2(t *testing.T) {
	actual := RightPad2Len("abc", "1", 5)
	const expected = "abc11"
	if actual != expected {
		t.Errorf("Expected: '%s', got '%s'", expected, actual)
	}
}
