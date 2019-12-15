/*
Copyright 2019 Jamie Thompson.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package holtwinters provides functionality for using the Holt-Winters exponential smoothing algorithm
// to make predictions and to smooth results for a time series.
// Built using these articles https://grisha.org/blog/2016/01/29/triple-exponential-smoothing-forecasting/
// Thanks to the author, Gregory Trubetskoy
package holtwinters

import "fmt"

// Predict takes in a seasonal historical series of data and produces a prediction of what the data will be in the future using triple
// exponential smoothing. Existing data will also be smoothed alongside predictions. Returns the entire dataset with the predictions
// appended to the end.
// series - Historical seasonal data, must be at least a full season, for optimal results use at least two full seasons,
// the first value should be at the start of a season
// seasonLength - The length of the data's seasons, must be at least 2
// alpha - Exponential smoothing coefficient for level, must be between 0 and 1
// beta - Exponential smoothing coefficient for trend, must be between 0 and 1
// gamma - Exponential smoothing coefficient for seasonality, must be between 0 and 1
// predictionLength - Number of predictions to make, set to 0 to make no predictions and only smooth, can't be negative
func Predict(series []float64, seasonLength int, alpha float64, beta float64, gamma float64, predictionLength int) ([]float64, error) {
	// Parameter validation mainly to avoid out of bounds errors and division by zero
	err := validateParams(series, seasonLength, alpha, beta, gamma, predictionLength)
	if err != nil {
		return nil, err
	}

	// Assumptions at this point, after params have been validated
	// seasonLength >= 2
	// series >= seasonLength
	// alpha, beta, gamma >= 0.0 and <= 1.0

	// Initial setup
	result := []float64{series[0]}
	smooth := series[0]
	trend := initialTrend(series, seasonLength)
	seasonals := initialSeasonalComponents(series, seasonLength)

	// Build prediction and smooth existing values
	for i := 1; i < len(series)+predictionLength; i++ {
		if i >= len(series) {
			// Prediction
			m := float64(i - len(series) + 1)
			result = append(result, (smooth+m*trend)+seasonals[i%seasonLength])
		} else {
			// Smooth existing values
			val := series[i]
			lastSmooth := smooth
			smooth = alpha*(val-seasonals[i%seasonLength]) + (1-alpha)*(smooth+trend)
			trend = beta*(smooth-lastSmooth) + (1-beta)*trend
			seasonals[i%seasonLength] = gamma*(val-smooth) + (1-gamma)*seasonals[i%seasonLength]
			result = append(result, smooth+trend+seasonals[i%seasonLength])
		}
	}
	return result, nil
}

// initialTrend calculates the initial trend based on average trends between the first and second
// seasons, if there is not enough data for two full seasons to be compared, instead the trend is
// calculated by comparing the first and second points of the first season
func initialTrend(series []float64, seasonLength int) float64 {
	// If not enough data to compare two seasons, more rough trend calculated using first two points
	if len(series) < seasonLength*2 {
		return series[1] - series[0]
	}

	// Enough data for two seasons, compare first two and average for trend
	sum := float64(0)
	for i := 0; i < seasonLength; i++ {
		sum += (series[i+seasonLength] - series[i]) / float64(seasonLength)
	}
	return sum / float64(seasonLength)
}

func initialSeasonalComponents(series []float64, seasonLength int) []float64 {
	var seasonals = make([]float64, seasonLength)
	seasonAverages := []float64{}
	nSeasons := len(series) / seasonLength
	for i := 0; i < nSeasons; i++ {
		// Calculate sum of season
		sum := float64(0)
		for j := seasonLength * i; j < seasonLength*i+seasonLength; j++ {
			sum += series[j]
		}
		// Calculate average of season and add to slice
		seasonAverages = append(seasonAverages, sum/float64(seasonLength))
	}
	for i := 0; i < seasonLength; i++ {
		sumOfValuesOverAverage := float64(0)
		for j := 0; j < nSeasons; j++ {
			sumOfValuesOverAverage += series[seasonLength*j+i] - seasonAverages[j]
		}
		seasonals[i] = sumOfValuesOverAverage / float64(nSeasons)
	}
	return seasonals
}

func validateParams(series []float64, seasonLength int, alpha float64, beta float64, gamma float64, predictionLength int) error {
	if seasonLength <= 1 {
		return fmt.Errorf("Invalid parameter for prediction; season length must be at least 2, is %d", seasonLength)
	}
	if predictionLength < 0 {
		return fmt.Errorf("Invalid parameter for prediction; prediction length must be at least 0, cannot be negative, is %d", predictionLength)
	}
	if alpha < 0.0 || alpha > 1.0 {
		return fmt.Errorf("Invalid parameter for prediction; alpha must be between 0 and 1, is %f", alpha)
	}
	if beta < 0.0 || beta > 1.0 {
		return fmt.Errorf("Invalid parameter for prediction; beta must be between 0 and 1, is %f", beta)
	}
	if gamma < 0.0 || gamma > 1.0 {
		return fmt.Errorf("Invalid parameter for prediction; gamma must be between 0 and 1, is %f", gamma)
	}
	if len(series) < seasonLength {
		return fmt.Errorf("Invalid parameter for prediction; must have at least 1 season of data to predict, season length: %d, series length: %d", seasonLength, len(series))
	}
	return nil
}
