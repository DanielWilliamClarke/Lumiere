package test

import (
	"testing"

	"dwc.com/lumiere/utils"
)

func Test_GeneratorGeneratesCodeOfLength(t *testing.T) {
	expectedLength := 10
	generator := utils.CodeGenerator{}
	code, err := generator.Generate(expectedLength)
	if err != nil {
		t.Errorf("Did not expect error: %v", err)
	}

	if len(code) != expectedLength {
		t.Errorf("Expected code of length %d", expectedLength)
	}
}
