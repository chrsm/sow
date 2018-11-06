package encoding

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestEncode(t *testing.T) {
	m := map[string]interface{}{
		//"x":   time.Now(),
		"y":   "y value",
		"num": 100,
		"payment_crap_lel": map[string]interface{}{
			"hmm": "hmm_value",
			"k": []interface{}{
				"k_value",
			},
		},
	}

	e := NewEncoder()
	x := e.Encode(m)

	spew.Dump(x)
}
