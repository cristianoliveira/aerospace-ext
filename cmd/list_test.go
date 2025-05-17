package cmd

import (
	"testing"

	"go.uber.org/mock/gomock"
)

func TestOutputSetter(t *testing.T) {
	t.Run("SetOutputPosition", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		out, err := cmdExecute("ls")
		if err != nil {
			t.Fatal(err)
		}

		t.Fatalf("Output: %s", out)
	})
}
