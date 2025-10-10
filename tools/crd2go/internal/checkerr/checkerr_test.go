package checkerr

import (
	"bytes"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckErr(t *testing.T) {
	var buf bytes.Buffer
	originalOutput := log.Writer()
	originalFlags := log.Flags()

	log.SetOutput(&buf)
	log.SetFlags(0)

	defer func() {
		log.SetOutput(originalOutput)
		log.SetFlags(originalFlags)
	}()

	tests := map[string]struct {
		msg         string
		f           funcErrs
		expectedLog string
	}{
		"no error": {
			msg: "Operation",
			f: funcErrs(func() error {
				return nil
			}),
			expectedLog: "",
		},
		"with error": {
			msg: "Operation",
			f: funcErrs(func() error {
				return assert.AnError
			}),
			expectedLog: "Operation failed: assert.AnError general error for testing\n",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			CheckErr(tt.msg, tt.f)
			logged := buf.String()
			buf.Reset()
			assert.Equal(t, tt.expectedLog, logged)
		})
	}
}
