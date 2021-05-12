package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/go-ping/ping"
	"gonum.org/v1/gonum/stat"
)

func execOneMeasure(target string, pingSamples, byteNum int, interval time.Duration) (time.Duration, time.Duration, error) {
	pinger, err := ping.NewPinger(target)

	if err != nil {
		return time.Second, time.Second, err
	}

	pinger.Count = pingSamples
	pinger.Size = byteNum
	pinger.Interval = interval
	pinger.SetPrivileged(true)

	err = pinger.Run()

	if err != nil {
		return time.Second, time.Second, err
	}
	stats := pinger.Statistics()

	return stats.MinRtt, stats.StdDevRtt, nil
}

func finalMeasures(target string, maxSize, numberOfSamples, numberOfPingSample int, abs, mes *[]float64) error {

	MinRTT, _, err := execOneMeasure(target, 5, 56, 100*time.Millisecond)

	if err != nil {
		panic(err)
	}

	for i := 0; i < numberOfSamples; i++ {
		bytenum := 32 + (maxSize-32)*i/numberOfSamples
		*abs = append(*abs, float64(8*bytenum))
		tempMinRTT, _, err := execOneMeasure(target, numberOfPingSample, bytenum, 2*MinRTT)

		MinRTT = tempMinRTT

		if err != nil {
			return err
		}
		*mes = append(*mes, float64(MinRTT)/1000000000)
	}

	return nil
}

func findCorrectI(MinRTT, baseMinRTT, baseStdDevRtt time.Duration, target string, numberOfSample, numberOfPingSample int) (int, error) {

	var correctI int

	for i := 6; i < 16; i++ {
		tempMinRTT, _, err := execOneMeasure(target, numberOfPingSample, int(math.Pow(2.0, float64(i))), 2*MinRTT)

		MinRTT = tempMinRTT

		if err != nil {
			return 0, err
		}

		if MinRTT-baseMinRTT > time.Duration(numberOfSample-1)*baseStdDevRtt || i == 15 {
			correctI = i
			break
		}
	}

	return correctI, nil

}

func main() {

	var numberOfSample int
	var numberOfPingSample int
	var correctI int
	var err error

	argsWithoutProg := os.Args[1:]

	target := argsWithoutProg[0]

	numberOfSample = 20
	numberOfPingSample = 20

	if len(argsWithoutProg) > 1 {
		numberOfSample, err = strconv.Atoi(argsWithoutProg[1])

		if err != nil {
			numberOfSample = 20
			fmt.Printf("Invalid value for numberOfSample : %s, replaced by %d\n", argsWithoutProg[1], numberOfSample)
		}
	}

	if len(argsWithoutProg) > 2 {
		numberOfPingSample, err = strconv.Atoi(argsWithoutProg[2])

		if err != nil {
			numberOfPingSample = 20
			fmt.Printf("Invalid value for numberOfPingSample : %s, replaced by %d\n", argsWithoutProg[2], numberOfSample)
		}
	}

	MinRTT, _, err := execOneMeasure(target, 5, 56, 100*time.Millisecond)

	if err != nil {
		panic(err)
	}

	baseMinRTT, baseStdDevRTT, err := execOneMeasure(target, numberOfPingSample, 56, MinRTT)

	if err != nil {
		panic(err)
	}

	correctI, err = findCorrectI(MinRTT, baseMinRTT, baseStdDevRTT, target, numberOfSample, numberOfPingSample)

	if err != nil {
		panic(err)
	}

	var abs = new([]float64)
	var mes = new([]float64)

	finalMeasures(target, int(math.Pow(2.0, float64(correctI))), numberOfSample, numberOfPingSample, abs, mes)

	_, a := stat.LinearRegression(*abs, *mes, nil, false)

	fmt.Println((2/a)/1000000, "Mbps")

}
