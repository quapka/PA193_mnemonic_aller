// Project: PA193_mnemonic_aller
// Mainteners UCO: 408788 497391 497577
// Description: Mnemonic API unit testing

package mnemonic

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

// FIXME rewrite table tests to have inplace struct definition
// FIXME testing for errors is not done nicely - problem with calling .Error() on 'nil'
func gotExp(got, exp string) string {
	return fmt.Sprintf("\nGot: %s\nExp: %s\n", got, exp)
}

func TestEntropyIsNotEmpty(t *testing.T) {
	entropy := ""
	phrase, seed, err := EntropyToPhraseAndSeed(entropy, "english.txt")

	expectedPhrase := ""
	expectedSeed := ""
	expectedErr := newEntropyIsEmptyError()

	if err.Error() != expectedErr.Error() {
		t.Error(gotExp(err.Error(), expectedErr.Error()))
	}
	if phrase != expectedPhrase {
		t.Error(gotExp(phrase, expectedPhrase))
	}
	if seed != expectedSeed {
		t.Error(gotExp(seed, expectedSeed))
	}
}

func TestENTIsInRange(t *testing.T) {
	// FIXME smaller, in range, higher, catch off by one errors
	// fairly small input
	entropy := "FF"
	expectedErr := newENTNotInRangeError()
	// FIXME
	_, _, err := EntropyToPhraseAndSeed(entropy, "english.txt")
	if err.Error() != expectedErr.Error() {
		t.Error(gotExp(err.Error(), expectedErr.Error()))
	}
	// smaller by 1
	entropy = "7FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF" // 2**127 - 1
	_, _, err = EntropyToPhraseAndSeed(entropy, "english.txt")
	if err.Error() != expectedErr.Error() {
		t.Error(gotExp(err.Error(), expectedErr.Error()))
	}
	// bigger by 1
	// FIXME different output this and Go Playground version, difference between 01
	// or 1 in the beginning
	entropy = "010000000000000000000000000000000000000000000000000000000000000000" // 2 ** 256
	_, _, err = EntropyToPhraseAndSeed(entropy, "english.txt")
	if err.Error() != expectedErr.Error() {
		t.Error(gotExp(err.Error(), expectedErr.Error()))
	}
}

func TestENTIsMultipleOf32x(t *testing.T) {
	entropy := "0800000000000000000000000000000000000000000000"
	_, _, err := EntropyToPhraseAndSeed(entropy, "english.txt")
	expectedErr := newEntropyNotDivisibleBy32Error(180)
	if err.Error() != expectedErr.Error() {
		t.Error(gotExp(err.Error(), expectedErr.Error()))
	}
}

func TestEntropyIsHexadecimal(t *testing.T) {
	entropy := "XX"
	phrase, seed, err := EntropyToPhraseAndSeed(entropy, "english.txt")
	expectedPhrase := ""
	expectedSeed := ""
	expectedErr := newEntropyIsNotHexadecimalError()
	// Comparing the actual error string
	if err.Error() != expectedErr.Error() {
		// if errors.Is(err, expectedErr) {
		t.Error(fmt.Sprintf("\nGot: %s\nExp: %s\n", err, expectedErr))
	}
	if phrase != expectedPhrase {
		t.Error(gotExp(phrase, expectedPhrase))
	}

	if seed != expectedSeed {
		t.Error(gotExp(seed, expectedSeed))
	}
}

func TestCannotOpenWordlistFile(t *testing.T) {
	entropy := "B7CB8EE904628CEC2B6779C0FB8B1B91" // 2**127 - 1
	phrase, seed, err := EntropyToPhraseAndSeed(entropy, "does/not/exist")
	expectedPhrase := ""
	expectedSeed := ""
	expectedErr := newOpenWordlistError("does/not/exist")

	if err.Error() != expectedErr.Error() {
		t.Error(gotExp(string(err.Error()), string(expectedErr.Error())))
	}
	if phrase != expectedPhrase {
		t.Error(gotExp(phrase, expectedPhrase))
	}

	if seed != expectedSeed {
		t.Error(gotExp(seed, expectedSeed))
	}
}

func TestFunc_cleanInputEntropy(t *testing.T) {

	testData := []struct {
		in_entropy string
		out_bytes  []byte
		out_error  error
	}{
		{in_entropy: "", out_bytes: nil, out_error: newEntropyIsEmptyError()},
		{in_entropy: "FF", out_bytes: nil, out_error: newENTNotInRangeError()},
		{in_entropy: "010000000000000000000000000000000000000000000000000000000000000000", out_error: newENTNotInRangeError()},
		{in_entropy: "0800000000000000000000000000000000000000000000", out_error: newEntropyNotDivisibleBy32Error(180)},
		{in_entropy: "XX", out_bytes: nil, out_error: newEntropyIsNotHexadecimalError()},
		{in_entropy: "B7CB8EE904628CEC2B6779C0FB8B1B91", out_bytes: []byte{0xB7, 0xCB, 0x8E, 0xE9, 0x04, 0x62, 0x8C, 0xEC, 0x2B, 0x67, 0x79, 0xC0, 0xFB, 0x8B, 0x1B, 0x91}, out_error: nil}, // 2**127 - 1
	}

	for i, td := range testData {
		bytes, err := cleanInputEntropy(td.in_entropy)
		if !reflect.DeepEqual(bytes, td.out_bytes) {
			t.Error(fmt.Sprintf("In %dth table-row", i+1))
			// FIXME give reasonable output!
			t.Error()
		}
		if td.out_error != nil {
			if err.Error() != td.out_error.Error() {
				t.Error(fmt.Sprintf("In %dth table-row", i+1))
				t.Error(gotExp(err.Error(), td.out_error.Error()))
			}
		} else {
			if err != nil {
				t.Error(gotExp(err.Error(), "'nil'"))
			}
		}
	}
}

func TestFunc_getBinaryLength(t *testing.T) {
	type testTemplate struct {
		input       []byte
		expectedLen int
	}
	testData := []testTemplate{
		{input: []byte{0x00}, expectedLen: 0},
		{input: []byte{0x05}, expectedLen: 3},
		{input: []byte{0x05, 0x00}, expectedLen: 11},
		{input: []byte{0x01, 0x00, 0x00, 0x00, 0x00}, expectedLen: 33},
		{input: []byte{0xAC, 0xFB, 0x96, 0x23, 0xE6,
			0x9A, 0x1F, 0xF0, 0xF7, 0xB7,
			0x2E, 0xDE, 0xED, 0x0A, 0x03,
			0xE7, 0xD8, 0x51, 0x3D, 0xE8,
			0xCB, 0x49, 0x73, 0x57, 0x56,
			0xD1, 0x15, 0xE1, 0x85, 0x8B,
			0x7F, 0x36, 0xA5, 0xA6, 0xE7,
			0x41, 0xAE, 0xBD, 0xFE, 0x2B,
			0x01, 0xAC, 0xC8, 0x73, 0x33,
			0x99, 0x19, 0x63, 0x64, 0xEE,
			0xD8, 0x0A, 0x21, 0x0A, 0x3C,
			0xED, 0x98, 0x63, 0xE3, 0x1B,
			0xB3, 0x71, 0x77, 0xCC, 0xAF,
			0x0D, 0xB6, 0x8E, 0x0A, 0x0B,
			0x4C, 0x92, 0x87, 0x10, 0xB4,
			0x2A, 0x7E, 0x9B, 0x87, 0x66,
			0x83, 0x6F, 0xCE, 0x0D, 0x6D,
			0xEA, 0x9F, 0x17, 0x62, 0x6E,
			0x50, 0x52, 0x90, 0x1E, 0x39,
			0x12, 0xCF, 0x49, 0x08, 0x22,
			0x59, 0xA1, 0xC6, 0x80, 0x1B,
			0x6D, 0xFB, 0x99, 0x08, 0x18,
			0xD7, 0x7B, 0x07, 0x5A, 0xFE,
			0x55, 0x69, 0x9C, 0xC9, 0x25,
			0x8A, 0xC5, 0x2F, 0xFE, 0x3B,
			0x63, 0x52, 0xF3, 0x34, 0x11,
			0x1B, 0x4A, 0x56, 0xCA, 0x02,
			0x2C, 0x6A, 0x13, 0x84, 0xD2,
			0xF6, 0xDB, 0xB1, 0x71, 0x62,
			0xF1, 0xB5, 0x20, 0xCC, 0x4B,
			0x76, 0x20, 0x35, 0xFD, 0x4E,
			0xB4, 0x7E, 0xA1, 0xF9, 0x6C,
			0xA3, 0x8F, 0xC1, 0x1C, 0xE4,
			0xCA, 0xAC, 0xB6, 0x28, 0x30,
			0x3E, 0xC2, 0x70, 0xDB, 0xD7,
			0x30, 0x58, 0xD0, 0x36, 0x14,
			0x31, 0x77, 0x4F, 0x2C, 0xA8,
			0x21, 0x50, 0x4F, 0x8F, 0x37,
			0x78, 0x58, 0xEB, 0x53, 0x67,
			0xA9, 0x89, 0xA2, 0x17, 0xE7}, expectedLen: 1600},
	}
	for i, td := range testData {
		got := getBinaryLength(td.input)
		if got != td.expectedLen {
			t.Error(fmt.Sprintf("In %dth table-row", i+1))
			t.Error(gotExp(strconv.Itoa(got), strconv.Itoa(td.expectedLen)))
		}
	}
}

// TODO finish
func TestFunc_calculateCheckSum(t *testing.T) {
	testData := []struct {
		in_bytes     []byte
		out_checkSum string
	}{
		{in_bytes: []byte{0xAE, 0x0F, 0x5C, 0xA4,
			0xDA, 0xB1, 0x59, 0x16,
			0x24, 0xA7, 0x3A, 0x0A,
			0x86, 0x00, 0xF5, 0xB2}, out_checkSum: "0100"},
		{in_bytes: []byte{0xD2, 0x18, 0x6C, 0x87,
			0x83, 0x28, 0xA2, 0x44,
			0xB6, 0x68, 0xF3, 0xF5,
			0xA8, 0x90, 0x55, 0x25,
			0xD7, 0xBF, 0xF7, 0x1C}, out_checkSum: "00001"},
		{in_bytes: []byte{0xB7, 0x52, 0x21, 0x89,
			0x86, 0xDA, 0x1E, 0x61,
			0xCB, 0x14, 0x70, 0x1C,
			0x57, 0x35, 0x78, 0xA1,
			0x8E, 0xC0, 0xFB, 0xBD,
			0x3C, 0xDF, 0x65, 0x42}, out_checkSum: "010000"},
		{in_bytes: []byte{0xB0, 0xE9, 0xE4, 0xDA,
			0x11, 0xE7, 0x84, 0x93,
			0x11, 0x14, 0xE7, 0x4D,
			0xE2, 0x44, 0x18, 0x69,
			0xBB, 0x8F, 0x59, 0xFE,
			0xFF, 0xB5, 0x15, 0x67,
			0x28, 0xC1, 0xAC, 0x01}, out_checkSum: "1011110"},
		{in_bytes: []byte{0xE3, 0x18, 0x71, 0xC8,
			0xFE, 0x1E, 0xC0, 0x01,
			0xBD, 0x10, 0x60, 0xBD,
			0x0C, 0x5D, 0xDC, 0xDF,
			0x54, 0x25, 0x73, 0xF5,
			0x11, 0xAE, 0x55, 0x47,
			0x6E, 0xDC, 0xCC, 0x93,
			0x59, 0x3B, 0x78, 0xD6}, out_checkSum: "11101111"},
	}
	for i, td := range testData {
		checksum := calculateCheckSum(td.in_bytes)
		if checksum != td.out_checkSum {
			t.Error(fmt.Sprintf("In %dth table-row", i+1))
			t.Error(gotExp(checksum, td.out_checkSum))
		}
	}
}

func TestFunc_convertToBinary(t *testing.T) {
	type testTemplate struct {
		input     []byte
		expOutput string
	}
	testData := []testTemplate{
		{input: []byte{0x00}, expOutput: "00000000"},
		{input: []byte{0xFF}, expOutput: "11111111"},
		{input: []byte{0x80, 0xFF}, expOutput: "1000000011111111"},
		{input: []byte{0xAC, 0xFB, 0x96, 0x23, 0xE6,
			0x9A, 0x1F, 0xF0, 0xF7, 0xB7,
			0x2E, 0xDE, 0xED, 0x0A, 0x03,
			0xE7, 0xD8, 0x51, 0x3D, 0xE8,
			0xCB, 0x49, 0x73, 0x57, 0x56,
			0xD1, 0x15, 0xE1, 0x85, 0x8B,
			0x7F, 0x36, 0xA5, 0xA6, 0xE7,
			0x41, 0xAE, 0xBD, 0xFE, 0x2B,
			0x01, 0xAC, 0xC8, 0x73, 0x33,
			0x99, 0x19, 0x63, 0x64, 0xEE,
			0xD8, 0x0A, 0x21, 0x0A, 0x3C,
			0xED, 0x98, 0x63, 0xE3, 0x1B,
			0xB3, 0x71, 0x77, 0xCC, 0xAF,
			0x0D, 0xB6, 0x8E, 0x0A, 0x0B,
			0x4C, 0x92, 0x87, 0x10, 0xB4,
			0x2A, 0x7E, 0x9B, 0x87, 0x66,
			0x83, 0x6F, 0xCE, 0x0D, 0x6D,
			0xEA, 0x9F, 0x17, 0x62, 0x6E,
			0x50, 0x52, 0x90, 0x1E, 0x39,
			0x12, 0xCF, 0x49, 0x08, 0x22,
			0x59, 0xA1, 0xC6, 0x80, 0x1B,
			0x6D, 0xFB, 0x99, 0x08, 0x18,
			0xD7, 0x7B, 0x07, 0x5A, 0xFE,
			0x55, 0x69, 0x9C, 0xC9, 0x25,
			0x8A, 0xC5, 0x2F, 0xFE, 0x3B,
			0x63, 0x52, 0xF3, 0x34, 0x11,
			0x1B, 0x4A, 0x56, 0xCA, 0x02,
			0x2C, 0x6A, 0x13, 0x84, 0xD2,
			0xF6, 0xDB, 0xB1, 0x71, 0x62,
			0xF1, 0xB5, 0x20, 0xCC, 0x4B,
			0x76, 0x20, 0x35, 0xFD, 0x4E,
			0xB4, 0x7E, 0xA1, 0xF9, 0x6C,
			0xA3, 0x8F, 0xC1, 0x1C, 0xE4,
			0xCA, 0xAC, 0xB6, 0x28, 0x30,
			0x3E, 0xC2, 0x70, 0xDB, 0xD7,
			0x30, 0x58, 0xD0, 0x36, 0x14,
			0x31, 0x77, 0x4F, 0x2C, 0xA8,
			0x21, 0x50, 0x4F, 0x8F, 0x37,
			0x78, 0x58, 0xEB, 0x53, 0x67,
			0xA9, 0x89, 0xA2, 0x17, 0xE7}, expOutput: "1010110011111011100101100010001111100110100110100001111111110000111101111011011100101110110111101110110100001010000000111110011111011000010100010011110111101000110010110100100101110011010101110101011011010001000101011110000110000101100010110111111100110110101001011010011011100111010000011010111010111101111111100010101100000001101011001100100001110011001100111001100100011001011000110110010011101110110110000000101000100001000010100011110011101101100110000110001111100011000110111011001101110001011101111100110010101111000011011011011010001110000010100000101101001100100100101000011100010000101101000010101001111110100110111000011101100110100000110110111111001110000011010110110111101010100111110001011101100010011011100101000001010010100100000001111000111001000100101100111101001001000010000010001001011001101000011100011010000000000110110110110111111011100110010000100000011000110101110111101100000111010110101111111001010101011010011001110011001001001001011000101011000101001011111111111000111011011000110101001011110011001101000001000100011011010010100101011011001010000000100010110001101010000100111000010011010010111101101101101110110001011100010110001011110001101101010010000011001100010010110111011000100000001101011111110101001110101101000111111010100001111110010110110010100011100011111100000100011100111001001100101010101100101101100010100000110000001111101100001001110000110110111101011100110000010110001101000000110110000101000011000101110111010011110010110010101000001000010101000001001111100011110011011101111000010110001110101101010011011001111010100110001001101000100001011111100111"},
	}
	for i, td := range testData {
		got := convertToBinary(td.input)
		if got != td.expOutput {
			t.Error(fmt.Sprintf("In %dth table-row", i+1))
			t.Error(gotExp(got, td.expOutput))
		}
	}
}

func TestFunc_createGroups(t *testing.T) {
	// FIXME check for expected error as well!
	type testTemplate struct {
		input     string
		expOutput []string
	}

	testData := []testTemplate{
		{input: "", expOutput: nil},
		{input: "11111111111", expOutput: []string{"11111111111"}},
		{input: "1111111111100000000000", expOutput: []string{"11111111111", "00000000000"}},
		{input: "1010", expOutput: nil},
	}

	for i, td := range testData {
		got, _ := createGroups(td.input)
		// fmt.Println("Got")
		// fmt.Println(got)
		if !reflect.DeepEqual(got, td.expOutput) {
			t.Error(fmt.Sprintf("In %dth table-row", i+1))
			t.Error(gotExp(strings.Join(got, ", "), strings.Join(td.expOutput, ", ")))
		}
	}
}

func TestFunc_createIndices(t *testing.T) {
	type testTemplate struct {
		input     []string
		expOutput []int64
	}

	// FIXME add more tests
	testData := []testTemplate{
		{input: []string{"10101010101", "11111111111"}, expOutput: []int64{1365, 2047}},
	}

	for i, td := range testData {
		got, _ := createIndices(td.input)
		if !reflect.DeepEqual(got, td.expOutput) {
			t.Error(fmt.Sprintf("In %dth table-row", i+1))
			// FIXME give more descriptive errors
		}
	}

}

func TestFunc_createPhraseWords(t *testing.T) {
	type testTemplate struct {
		indices []int64
		words   []string
		phrase  []string
	}
	// FIXME add more testData
	testData := []testTemplate{
		{indices: []int64{1, 1}, words: []string{"hello", "world"}, phrase: []string{"world", "world"}},
	}
	for i, td := range testData {
		got, _ := createPhraseWords(td.indices, td.words)
		if !reflect.DeepEqual(got, td.phrase) {
			t.Error(fmt.Sprintf("In %dth table-row", i+1))
			t.Error(gotExp(strings.Join(got, ", "), strings.Join(td.phrase, ", ")))
		}
	}
}

func TestFunc_cleanLine(t *testing.T) {
	type testTemplate struct {
		input     string
		expOutput string
	}

	testData := []testTemplate{
		{input: "", expOutput: ""},
		{input: "WORD", expOutput: "word"},
		{input: "\n word", expOutput: "word"},
		{input: "word \n", expOutput: "word"},
		{input: "\t word \n", expOutput: "word"},
		{input: "\t WORD \n", expOutput: "word"},
	}

	for i, td := range testData {
		got := cleanLine(td.input)
		if got != td.expOutput {
			t.Error(fmt.Sprintf("In %dth table-row", i+1))
			t.Error(gotExp(got, td.expOutput))
		}
	}
}

func TestFunc_validateWord(t *testing.T) {
	type testTemplate struct {
		input string
		valid bool
	}

	testData := []testTemplate{
		{input: "", valid: false},
		{input: " word", valid: false},
		{input: "word ", valid: false},
		{input: "word word", valid: false},
		{input: "word", valid: true},
		{input: "longword", valid: true},
		{input: "WORD", valid: false},
		{input: "veryveryveryveryveryverylongword", valid: true},
		{input: "1word", valid: false},
		// FIXME we should probably support Unicode
		{input: "š", valid: false},
	}
	for i, td := range testData {
		got := validateWord(td.input)
		if got != td.valid {
			t.Error(fmt.Sprintf("In %dth table-row", i+1))
			t.Error(gotExp(strconv.FormatBool(got), strconv.FormatBool(td.valid)))
		}
	}
}

var wlfile = "../../wordlists/english.txt"

var testVectors = []struct {
	entropy string
	phrase  string
	seed    string
}{
	{
		"00000000000000000000000000000000",
		"abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about",
		"c55257c360c07c72029aebc1b53c05ed0362ada38ead3e3e9efa3708e53495531f09a6987599d18264c1e1c92f2cf141630c7a3c4ab7c81b2f001698e7463b04",
	},
	{
		"7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f",
		"legal winner thank year wave sausage worth useful legal winner thank yellow",
		"2e8905819b8723fe2c1d161860e5ee1830318dbf49a83bd451cfb8440c28bd6fa457fe1296106559a3c80937a1c1069be3a3a5bd381ee6260e8d9739fce1f607",
	},
	{
		"80808080808080808080808080808080",
		"letter advice cage absurd amount doctor acoustic avoid letter advice cage above",
		"d71de856f81a8acc65e6fc851a38d4d7ec216fd0796d0a6827a3ad6ed5511a30fa280f12eb2e47ed2ac03b5c462a0358d18d69fe4f985ec81778c1b370b652a8",
	},
	{
		"ffffffffffffffffffffffffffffffff",
		"zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo wrong",
		"ac27495480225222079d7be181583751e86f571027b0497b5b5d11218e0a8a13332572917f0f8e5a589620c6f15b11c61dee327651a14c34e18231052e48c069",
	},
	{
		"000000000000000000000000000000000000000000000000",
		"abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon agent",
		"035895f2f481b1b0f01fcf8c289c794660b289981a78f8106447707fdd9666ca06da5a9a565181599b79f53b844d8a71dd9f439c52a3d7b3e8a79c906ac845fa",
	},
	{
		"7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f",
		"legal winner thank year wave sausage worth useful legal winner thank year wave sausage worth useful legal will",
		"f2b94508732bcbacbcc020faefecfc89feafa6649a5491b8c952cede496c214a0c7b3c392d168748f2d4a612bada0753b52a1c7ac53c1e93abd5c6320b9e95dd",
	},
	{
		"808080808080808080808080808080808080808080808080",
		"letter advice cage absurd amount doctor acoustic avoid letter advice cage absurd amount doctor acoustic avoid letter always",
		"107d7c02a5aa6f38c58083ff74f04c607c2d2c0ecc55501dadd72d025b751bc27fe913ffb796f841c49b1d33b610cf0e91d3aa239027f5e99fe4ce9e5088cd65",
	},
	{
		"ffffffffffffffffffffffffffffffffffffffffffffffff",
		"zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo when",
		"0cd6e5d827bb62eb8fc1e262254223817fd068a74b5b449cc2f667c3f1f985a76379b43348d952e2265b4cd129090758b3e3c2c49103b5051aac2eaeb890a528",
	},
	{
		"0000000000000000000000000000000000000000000000000000000000000000",
		"abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art",
		"bda85446c68413707090a52022edd26a1c9462295029f2e60cd7c4f2bbd3097170af7a4d73245cafa9c3cca8d561a7c3de6f5d4a10be8ed2a5e608d68f92fcc8",
	},
	{
		"7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f",
		"legal winner thank year wave sausage worth useful legal winner thank year wave sausage worth useful legal winner thank year wave sausage worth title",
		"bc09fca1804f7e69da93c2f2028eb238c227f2e9dda30cd63699232578480a4021b146ad717fbb7e451ce9eb835f43620bf5c514db0f8add49f5d121449d3e87",
	},
	{
		"8080808080808080808080808080808080808080808080808080808080808080",
		"letter advice cage absurd amount doctor acoustic avoid letter advice cage absurd amount doctor acoustic avoid letter advice cage absurd amount doctor acoustic bless",
		"c0c519bd0e91a2ed54357d9d1ebef6f5af218a153624cf4f2da911a0ed8f7a09e2ef61af0aca007096df430022f7a2b6fb91661a9589097069720d015e4e982f",
	},
	{
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
		"zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo vote",
		"dd48c104698c30cfe2b6142103248622fb7bb0ff692eebb00089b32d22484e1613912f0a5b694407be899ffd31ed3992c456cdf60f5d4564b8ba3f05a69890ad",
	},
	{
		"9e885d952ad362caeb4efe34a8e91bd2",
		"ozone drill grab fiber curtain grace pudding thank cruise elder eight picnic",
		"274ddc525802f7c828d8ef7ddbcdc5304e87ac3535913611fbbfa986d0c9e5476c91689f9c8a54fd55bd38606aa6a8595ad213d4c9c9f9aca3fb217069a41028",
	},
	{
		"6610b25967cdcca9d59875f5cb50b0ea75433311869e930b",
		"gravity machine north sort system female filter attitude volume fold club stay feature office ecology stable narrow fog",
		"628c3827a8823298ee685db84f55caa34b5cc195a778e52d45f59bcf75aba68e4d7590e101dc414bc1bbd5737666fbbef35d1f1903953b66624f910feef245ac",
	},
	{
		"68a79eaca2324873eacc50cb9c6eca8cc68ea5d936f98787c60c7ebc74e6ce7c",
		"hamster diagram private dutch cause delay private meat slide toddler razor book happy fancy gospel tennis maple dilemma loan word shrug inflict delay length",
		"64c87cde7e12ecf6704ab95bb1408bef047c22db4cc7491c4271d170a1b213d20b385bc1588d9c7b38f1b39d415665b8a9030c9ec653d75e65f847d8fc1fc440",
	},
	{
		"c0ba5a8e914111210f2bd131f3d5e08d",
		"scheme spot photo card baby mountain device kick cradle pact join borrow",
		"ea725895aaae8d4c1cf682c1bfd2d358d52ed9f0f0591131b559e2724bb234fca05aa9c02c57407e04ee9dc3b454aa63fbff483a8b11de949624b9f1831a9612",
	},
	{
		"6d9be1ee6ebd27a258115aad99b7317b9c8d28b6d76431c3",
		"horn tenant knee talent sponsor spell gate clip pulse soap slush warm silver nephew swap uncle crack brave",
		"fd579828af3da1d32544ce4db5c73d53fc8acc4ddb1e3b251a31179cdb71e853c56d2fcb11aed39898ce6c34b10b5382772db8796e52837b54468aeb312cfc3d",
	},
	{
		"9f6a2878b2520799a44ef18bc7df394e7061a224d2c33cd015b157d746869863",
		"panda eyebrow bullet gorilla call smoke muffin taste mesh discover soft ostrich alcohol speed nation flash devote level hobby quick inner drive ghost inside",
		"72be8e052fc4919d2adf28d5306b5474b0069df35b02303de8c1729c9538dbb6fc2d731d5f832193cd9fb6aeecbc469594a70e3dd50811b5067f3b88b28c3e8d",
	},
	{
		"23db8160a31d3e0dca3688ed941adbf3",
		"cat swing flag economy stadium alone churn speed unique patch report train",
		"deb5f45449e615feff5640f2e49f933ff51895de3b4381832b3139941c57b59205a42480c52175b6efcffaa58a2503887c1e8b363a707256bdd2b587b46541f5",
	},
	{
		"8197a4a47f0425faeaa69deebc05ca29c0a5b5cc76ceacc0",
		"light rule cinnamon wrap drastic word pride squirrel upgrade then income fatal apart sustain crack supply proud access",
		"4cbdff1ca2db800fd61cae72a57475fdc6bab03e441fd63f96dabd1f183ef5b782925f00105f318309a7e9c3ea6967c7801e46c8a58082674c860a37b93eda02",
	},
	{
		"066dca1a2bb7e8a1db2832148ce9933eea0f3ac9548d793112d9a95c9407efad",
		"all hour make first leader extend hole alien behind guard gospel lava path output census museum junior mass reopen famous sing advance salt reform",
		"26e975ec644423f4a4c4f4215ef09b4bd7ef924e85d1d17c4cf3f136c2863cf6df0a475045652c57eb5fb41513ca2a2d67722b77e954b4b3fc11f7590449191d",
	},
	{
		"f30f8c1da665478f49b001d94c5fc452",
		"vessel ladder alter error federal sibling chat ability sun glass valve picture",
		"2aaa9242daafcee6aa9d7269f17d4efe271e1b9a529178d7dc139cd18747090bf9d60295d0ce74309a78852a9caadf0af48aae1c6253839624076224374bc63f",
	},
	{
		"c10ec20dc3cd9f652c7fac2f1230f7a3c828389a14392f05",
		"scissors invite lock maple supreme raw rapid void congress muscle digital elegant little brisk hair mango congress clump",
		"7b4a10be9d98e6cba265566db7f136718e1398c71cb581e1b2f464cac1ceedf4f3e274dc270003c670ad8d02c4558b2f8e39edea2775c9e232c7cb798b069e88",
	},
	{
		"f585c11aec520db57dd353c69554b21a89b20fb0650966fa0a9d6f74fd989d8f",
		"void come effort suffer camp survey warrior heavy shoot primary clutch crush open amazing screen patrol group space point ten exist slush involve unfold",
		"01f5bced59dec48e362f2c45b5de68b9fd6c92c6634f44d6d40aab69056506f0e35524a518034ddc1192e1dacd32c1ed3eaa3c3b131c88ed8e7e54c49a5d0998",
	},
}

func TestEntropyToPhraseAndSeed(t *testing.T) {

	for _, v := range testVectors {
		entropy, e := PhraseToEntropyAndSeed(v.phrase, wlfile)
		if e != nil {
			t.Errorf("Phrase to entropy and seed function failed: %s", e)
		}
		if entropy != v.entropy {
			t.Errorf("Got unexpected entropy. Expected %s, got: %s", v.entropy, entropy)
		}
	}
}
