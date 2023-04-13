package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name:  "test",
		SKU:   "dadsa-dsad-dsad",
		Price: 2.0,
	}
	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
