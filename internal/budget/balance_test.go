package budget

import (
	"encoding/json"
	"os"
	"testing"
)

type balanceCase struct {
	Name         string        `json:"name"`
	Transactions []Transaction `json:"transactions"`
	Want         int64         `json:"want"`
}

func TestBalance(t *testing.T) {
	data, err := os.ReadFile("testdata/data.json")
	if err != nil {
		t.Fatal(err)
	}

	var cases []balanceCase
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatalf("bad test file: %v", err)
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			store := NewMemoryStore()
			for _, tx := range tc.Transactions {
				store.Add(tx)
			}

			if got := Balance(store); got != tc.Want {
				t.Errorf("got %d, want %d", got, tc.Want)
			}
		})
	}

}
