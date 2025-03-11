package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

// For use color with String()
var (
	colorHot  = color.New(color.FgRed).SprintFunc()  // Red
	colorCold = color.New(color.FgBlue).SprintFunc() // Blue
)

type SensorInfo struct {
	Temps           []int // Each element has a temp. sensor reading
	SensorsCount    int   // How many sensors?
	MinReadableTemp int
	MaxReadableTemp int
	MinSafeTemp     int
	MaxSafeTemp     int
}

// ReadTemps calls the mock server that returns random numbers.
// sets those number in Temps
func (si *SensorInfo) ReadTemps() error {
	url := fmt.Sprintf(
		"https://www.random.org/integers/?num=%d&min=%d&max=%d&col=1&base=10&format=plain&rnd=new",
		si.SensorsCount, si.MinReadableTemp, si.MaxReadableTemp)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	// The response is in plain text, one line by temp. Trim the lines and split by "\n"
	lines := strings.Split(strings.TrimSpace(string(body)), "\n")
	for i, line := range lines {
		number, err := strconv.Atoi(line)
		if err != nil {
			return err
		}
		si.Temps[i] = number
	}
	return nil
}

// SensorTempIsSafe returns true if the temp of the indicated sensor in safe range
func (si *SensorInfo) SensorTempIsSafe(sensorNumber int) bool {
	return si.Temps[sensorNumber] >= si.MinSafeTemp && si.Temps[sensorNumber] <= si.MaxSafeTemp
}

// String Implements Print. Diplay sensors and they temps.
// Temps out of safe range are market with "*"
func (si SensorInfo) String() string { // Only works if SensorInfo is passed by value, not by reference
	result := "Sensors: "
	for i, temp := range si.Temps {

		tempString := fmt.Sprintf("%4d", temp) // For color use must be converted to string.
		// Just a more polite output. Blue for cold temp. and red for hot.
		if !si.SensorTempIsSafe(i) {
			if temp < si.MinSafeTemp {
				tempString = colorCold(fmt.Sprintf("%4d", temp))
			} else if temp > si.MaxSafeTemp {
				tempString = colorHot(fmt.Sprintf("%4d", temp))
			}
		}
		result += fmt.Sprintf("[%2d=%s] ", i, tempString)
	}
	return result
}
