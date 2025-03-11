package main

// Unit Test: TempIsSafe

import "testing"

func TestSensorTempIsSafe(t *testing.T) {
	si := SensorInfo{
		Temps:       make([]int, 6), //Sensors count
		MinSafeTemp: 125,            // Set same values as config.ini
		MaxSafeTemp: 350,
	}

	// Create the struct "in-line" good for small test.
	tests := []struct {
		sensorNumber int
		temp         int
		expected     bool
	}{
		{0, 100, false},
		{1, 125, true},
		{2, 200, true},
		{3, 350, true},
		{4, 381, false},
		{5, 400, false},
	}

    // Set the temp value and Check if safe
	for _, tt := range tests {
        si.Temps[tt.sensorNumber] = tt.temp
		result := si.SensorTempIsSafe(tt.sensorNumber)
		if result != tt.expected {
			t.Errorf("Temp %d, expected %v, but got %v", tt.temp, tt.expected, result)
		}
	}
}
