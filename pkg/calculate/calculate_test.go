package calculate

import (
	"reflect"
	"testing"
)

func TestCalculatePackages(t *testing.T) {
	packSizes := []int64{250, 500, 1000, 2000, 5000}

	tests := []struct {
		name         string
		orderedItems int64
		want         map[int64]int64
		description  string
	}{
		{
			name:         "1 item",
			orderedItems: 1,
			want:         map[int64]int64{250: 1},
			description:  "Should use smallest pack (250) for minimal overshoot",
		},
		{
			name:         "Exact 250",
			orderedItems: 250,
			want:         map[int64]int64{250: 1},
			description:  "Should use one 250 pack for exact match",
		},
		{
			name:         "251 items",
			orderedItems: 251,
			want:         map[int64]int64{500: 1},
			description:  "Should use 500 pack (not 2x250) - fewer packs",
		},
		{
			name:         "501 items",
			orderedItems: 501,
			want:         map[int64]int64{500: 1, 250: 1},
			description:  "Should use 500+250 (not 1000) - less overshoot",
		},
		{
			name:         "12001 items",
			orderedItems: 12001,
			want:         map[int64]int64{5000: 2, 2000: 1, 250: 1},
			description:  "Should minimize overshoot first, then pack count",
		},
		{
			name:         "751 items",
			orderedItems: 751,
			want:         map[int64]int64{1000: 1},
			description:  "Should use single 1000 pack",
		},
		{
			name:         "1001 items",
			orderedItems: 1001,
			want:         map[int64]int64{1000: 1, 250: 1},
			description:  "Should use 1000+250 for minimal overshoot",
		},
		{
			name:         "2751 items",
			orderedItems: 2751,
			want:         map[int64]int64{2000: 1, 1000: 1},
			description:  "Should use 2000+1000 (not other combinations)",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculatePackages(tt.orderedItems, packSizes)
			if err != nil {
				t.Fatalf("CalculatePackages(%d) unexpected error: %v", tt.orderedItems, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalculatePackages(%d) = %v, want %v\n%s",
					tt.orderedItems, got, tt.want, tt.description)
			}
		})
	}
}

func TestCalculatePackages_EdgeCases(t *testing.T) {
	t.Run("Empty pack sizes", func(t *testing.T) {
		if _, err := CalculatePackages(100, nil); err != ErrEmptySizes {
			t.Fatalf("expected ErrEmptySizes, got %v", err)
		}
	})

	t.Run("Single pack size", func(t *testing.T) {
		got, err := CalculatePackages(751, []int64{500})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := map[int64]int64{500: 2}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("CalculatePackages(751, [500]) = %v, want %v", got, want)
		}
	})

	t.Run("Impossible to fulfill with given pack sizes", func(t *testing.T) {
		got, err := CalculatePackages(3, []int64{2})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := map[int64]int64{2: 2} // 4 items, overshoot of 1
		if !reflect.DeepEqual(got, want) {
			t.Errorf("CalculatePackages(3, [2]) = %v, want %v", got, want)
		}
	})
}
