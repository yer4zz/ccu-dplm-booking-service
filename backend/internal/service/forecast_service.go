package service

import (
	"math"
	"time"
)

type ForecastPoint struct {
	Day        string  `json:"day"`
	Revenue    float64 `json:"revenue"`
	IsForecast bool    `json:"is_forecast"`
	Lower      float64 `json:"lower"`
	Upper      float64 `json:"upper"`
}

func HoltWinters(
	data    []float64,
	days    []time.Time,
	periods int,
) []ForecastPoint {
	n := len(data)
	if n < 4 {
		return linearFallback(data, days, periods)
	}

	alpha := 0.3
	beta  := 0.1
	gamma := 0.2

	season := 7

	seasonCoeffs := make([]float64, season)
	for i := 0; i < season && i < n; i++ {
		seasonCoeffs[i] = 1.0
	}

	
	if n >= season {
		weekCount := n / season
		weekAvgs  := make([]float64, weekCount)
		for w := 0; w < weekCount; w++ {
			sum := 0.0
			for d := 0; d < season; d++ {
				sum += data[w*season+d]
			}
			weekAvgs[w] = sum / float64(season)
		}
		for d := 0; d < season; d++ {
			sum := 0.0
			cnt := 0
			for w := 0; w < weekCount; w++ {
				if weekAvgs[w] > 0 {
					sum += data[w*season+d] / weekAvgs[w]
					cnt++
				}
			}
			if cnt > 0 {
				seasonCoeffs[d] = sum / float64(cnt)
			} else {
				seasonCoeffs[d] = 1.0
			}
		}
	}

	level := data[0]
	trend := 0.0
	if n > 1 {
		trend = (data[n-1] - data[0]) / float64(n-1)
	}

	smoothed := make([]float64, n)
	for i := 0; i < n; i++ {
		si  := i % season
		obs := data[i]

		prevLevel := level
		sc        := seasonCoeffs[si]
		if sc == 0 { sc = 1.0 }


		level = alpha*(obs/sc) + (1-alpha)*(prevLevel+trend)

		trend = beta*(level-prevLevel) + (1-beta)*trend
		if level > 0 {
			seasonCoeffs[si] = gamma*(obs/level) + (1-gamma)*sc
		}

		smoothed[i] = (level + trend) * seasonCoeffs[si]
	}

	var sumSq float64
	for i, v := range smoothed {
		diff := v - data[i]
		sumSq += diff * diff
	}
	stdDev := math.Sqrt(sumSq / float64(n))

	result := make([]ForecastPoint, 0, n+periods)
	for i, d := range days {
		result = append(result, ForecastPoint{
			Day:        d.Format("2006-01-02"),
			Revenue:    math.Round(data[i]),
			IsForecast: false,
			Lower:      math.Max(0, math.Round(smoothed[i]-stdDev)),
			Upper:      math.Round(smoothed[i] + stdDev),
		})
	}

	lastDay := days[n-1]
	for i := 1; i <= periods; i++ {
		si  := (n + i - 1) % season
		sc  := seasonCoeffs[si]
		if sc == 0 { sc = 1.0 }

		val := (level + float64(i)*trend) * sc
		if val < 0 { val = 0 }
		val = math.Round(val)

		uncertainty := stdDev * math.Sqrt(float64(i))

		result = append(result, ForecastPoint{
			Day:        lastDay.AddDate(0, 0, i).Format("2006-01-02"),
			Revenue:    val,
			IsForecast: true,
			Lower:      math.Max(0, math.Round(val-uncertainty)),
			Upper:      math.Round(val + uncertainty),
		})
	}

	return result
}

func linearFallback(data []float64, days []time.Time, periods int) []ForecastPoint {
	n := len(data)
	result := make([]ForecastPoint, 0, n+periods)

	for i, d := range days {
		result = append(result, ForecastPoint{
			Day:        d.Format("2006-01-02"),
			Revenue:    math.Round(data[i]),
			IsForecast: false,
		})
	}

	if n == 0 {
		return result
	}

	slope := 0.0
	if n > 1 {
		slope = (data[n-1] - data[0]) / float64(n-1)
	}
	lastVal := data[n-1]
	lastDay := days[n-1]

	for i := 1; i <= periods; i++ {
		val := math.Max(0, math.Round(lastVal+slope*float64(i)))
		result = append(result, ForecastPoint{
			Day:        lastDay.AddDate(0, 0, i).Format("2006-01-02"),
			Revenue:    val,
			IsForecast: true,
		})
	}
	return result
}