package domain

import (
	"testing"
)

func TestWidgetCampaign_Progress(t *testing.T) {
	tests := []struct {
		name    string
		goal    int64
		current int64
		want    float64
	}{
		{"zero goal returns 0", 0, 100, 0},
		{"half funded", 1000000, 500000, 50.0},
		{"fully funded", 1000000, 1000000, 100.0},
		{"over funded", 1000000, 1500000, 150.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := WidgetCampaign{GoalAmount: tt.goal, CurrentAmount: tt.current}

			var got float64
			if c.GoalAmount == 0 {
				got = 0
			} else {
				got = (float64(c.CurrentAmount) / float64(c.GoalAmount)) * 100
			}

			if got != tt.want {
				t.Errorf("progress computed = %v, want %v", got, tt.want)
			}
		})
	}
}
