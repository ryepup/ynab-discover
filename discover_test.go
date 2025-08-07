package discover_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ryepup/ynab-discover"
)

func TestConvertCSV(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input    string
		expected string
		wantErr  bool
	}{
		"basic conversion": {
			input: `Trans. Date,Post Date,Description,Amount,Category
06/23/2023,06/26/2023,"ACME STORE, INC 555-1234567 NY",50.00,"Services"
06/24/2023,06/26/2023,"GROCERY MART #123 ANYTOWN FL",163.73,"Supermarkets"`,
			expected: `Date,Post Date,Description,Amount,Category
06/23/2023,06/26/2023,"ACME STORE, INC 555-1234567 NY",-50.00,Services
06/24/2023,06/26/2023,GROCERY MART #123 ANYTOWN FL,-163.73,Supermarkets
`,
			wantErr: false,
		},
		"negative amounts become positive": {
			input: `Trans. Date,Post Date,Description,Amount,Category
06/23/2023,06/26/2023,"PAYMENT REFUND",-50.00,"Services"`,
			expected: `Date,Post Date,Description,Amount,Category
06/23/2023,06/26/2023,PAYMENT REFUND,50.00,Services
`,
			wantErr: false,
		},
		"missing Trans. Date column": {
			input: `Date,Post Date,Description,Amount,Category
06/23/2023,06/26/2023,"TEST",50.00,"Services"`,
			expected: "",
			wantErr:  true,
		},
		"missing Amount column": {
			input: `Trans. Date,Post Date,Description,Category
06/23/2023,06/26/2023,"TEST","Services"`,
			expected: "",
			wantErr:  true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			reader := strings.NewReader(tt.input)
			var writer bytes.Buffer

			err := discover.ConvertCSV(t.Context(), reader, &writer)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertCSV() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				got := writer.String()
				if got != tt.expected {
					t.Errorf("ConvertCSV() output mismatch:\nGot:\n%s\nExpected:\n%s", got, tt.expected)
				}
			}
		})
	}
}
