package stringutils

import (
	"testing"
)

type stripCharsTestData struct {
	str      string
	chr      string
	expected string
}

var stripCharsTests = []stripCharsTestData{
	{"qwerty.asdfgh.", ".", "qwertyasdfgh"},
	{"qwerty.asdfgh.", ",", "qwerty.asdfgh."},
	{"qwerty.asdfgh.", " ", "qwerty.asdfgh."},
	{"abc", ".qwe", "abc"},
	{"123456", "123456", ""},
	{"1234567", "123456", "7"},
	{"1234567", "123456.", "7"},
}

func Test_Stripchars(t *testing.T) {
	for _, tt := range stripCharsTests {
		actual := Stripchars(tt.str, tt.chr)
		if actual != tt.expected {
			t.Errorf("Expected: '%s', got '%s'", tt.expected, actual)
		}
	}
}

type padTestData struct {
	s        string
	padStr   string
	pLen     int
	expected string
}

var leftPadTests = []padTestData{
	{"aa", "0", 3, "000aa"},
	{"abc", "g", 3, "gggabc"},
	{"1", "0", 0, "1"},
	{"TEST", "0", 5, "00000TEST"},
	{"TEST", "xx", 3, "xxxxxxTEST"},
	{"TEST", "xx", 25, "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxTEST"},
}

func Test_LeftPad(t *testing.T) {
	for _, tt := range leftPadTests {
		actual := LeftPad(tt.s, tt.padStr, tt.pLen)
		if actual != tt.expected {
			t.Errorf("Expected: '%s', got '%s'", tt.expected, actual)
		}
	}
}

var rightPadTests = []padTestData{
	{"aa", "0", 3, "aa000"},
	{"abc", "g", 3, "abcggg"},
	{"1", "0", 0, "1"},
	{"TEST", "0", 5, "TEST00000"},
	{"TEST", "xx", 3, "TESTxxxxxx"},
	{"TEST", "xx", 20, "TESTxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"},
}

func Test_RightPad(t *testing.T) {
	for _, tt := range rightPadTests {
		actual := RightPad(tt.s, tt.padStr, tt.pLen)
		if actual != tt.expected {
			t.Errorf("Expected: '%s', got '%s'", tt.expected, actual)
		}
	}
}

var leftPad2LenTests = []padTestData{
	{"aa", "0", 3, "0aa"},
	{"abc", "g", 5, "ggabc"},
	{"1", "0", 1, "1"},
	{"1", "0", 0, ""},
	{"TEST", "0", 5, "0TEST"},
	{"TEST", "0", 50, "0000000000000000000000000000000000000000000000TEST"},
}

func Test_LeftPad2Len(t *testing.T) {
	for _, tt := range leftPad2LenTests {
		actual := LeftPad2Len(tt.s, tt.padStr, tt.pLen)
		if actual != tt.expected {
			t.Errorf("Expected: '%s', got '%s'", tt.expected, actual)
		}
	}
}

var rightPad2LenTests = []padTestData{
	{"aa", "0", 3, "aa0"},
	{"abc", "g", 5, "abcgg"},
	{"1", "0", 1, "1"},
	{"1", "0", 0, ""},
	{"TEST", "0", 5, "TEST0"},
	{"TEST", "0", 50, "TEST0000000000000000000000000000000000000000000000"},
}

func Test_RightPad2Len(t *testing.T) {
	for _, tt := range rightPad2LenTests {
		actual := RightPad2Len(tt.s, tt.padStr, tt.pLen)
		if actual != tt.expected {
			t.Errorf("Expected: '%s', got '%s'", tt.expected, actual)
		}
	}
}
