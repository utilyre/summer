package pipeline_test

import (
	"testing"

	"github.com/utilyre/summer/pkg/pipeline"
)

func TestPipeline_Append(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("wanted panic; got no panic")
		}
	}()

	var pl pipeline.Pipeline[struct{}]
	pl.Append(0, nil)
	pl.Append(-2, nil)
	pl.AppendFunc(0, nil)
	pl.AppendFunc(-2, nil)
}
