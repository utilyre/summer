package pipeline

import "testing"

func TestPipeline_Append(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("wanted panic; got no panic")
		}
	}()

	var pl Pipeline
	pl.Append(0, nil)
	pl.Append(-2, nil)
	pl.AppendFunc(0, nil)
	pl.AppendFunc(-2, nil)
}
