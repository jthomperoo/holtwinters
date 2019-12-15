[![Build](https://github.com/jthomperoo/holtwinters/workflows/main/badge.svg)](https://github.com/jthomperoo/holtwinters)
[![codecov](https://codecov.io/gh/jthomperoo/holtwinters/branch/master/graph/badge.svg)](https://codecov.io/gh/jthomperoo/holtwinters)
[![GoDoc](https://godoc.org/github.com/jthomperoo/holtwinters?status.svg)](https://godoc.org/github.com/jthomperoo/holtwinters)
[![Go Report Card](https://goreportcard.com/badge/github.com/jthomperoo/holtwinters)](https://goreportcard.com/report/github.com/jthomperoo/holtwinters)
[![License](http://img.shields.io/:license-apache-blue.svg)](http://www.apache.org/licenses/LICENSE-2.0.html)

# Holt Winters Exponential Smoothing

Package holtwinters provides functionality for using the Holt-Winters exponential smoothing algorithm to make predictions and to smooth results for a time series.

Built using these articles [https://grisha.org/blog/2016/01/29/triple-exponential-smoothing-forecasting/](https://grisha.org/blog/2016/01/29/triple-exponential-smoothing-forecasting/).  
Thanks to the author, Gregory Trubetskoy.

## Installation

```
go get -u github.com/jthomperoo/holtwinters
```

## Reference

This package exposes a single function:

```go
Predict(series []float64, seasonLength int, alpha float64, beta float64, gamma float64, predictionLength int) ([]float64, error)
```
Predict takes in a seasonal historical series of data and produces a prediction of what the data will be in the future using triple
exponential smoothing. Existing data will also be smoothed alongside predictions. Returns the entire dataset with the predictions
appended to the end.  
 - **series** - Historical seasonal data, must be at least a full season, for optimal results use at least two full seasons, the first value should be at the start of a season
 - **seasonLength** - The length of the data's seasons, must be at least 2
 - **alpha** - Exponential smoothing coefficient for level, must be between 0 and 1
 - **beta** - Exponential smoothing coefficient for trend, must be between 0 and 1
 - **gamma** - Exponential smoothing coefficient for seasonality, must be between 0 and 1
 - **predictionLength** - Number of predictions to make, set to 0 to make no predictions and only smooth, can't be negative  

Returns the full series that has been smoothed, with predictions appended to the end. The only errors that can be returned are parameter validation errors, such as season length being too short, or alpha, beta, or gamma values being beyond 0-1.

## Developing

### Environment

Developing this project requires these dependencies:

* `Go 1.13`
* `Golint`

### Pipeline

This project uses GitHub Actions to run a pipeline against every commit, pull request and release that lints and runs the tests, if any warnings/errors are found or the tests fail, the pipeline will fail.

### Commands

This project includes a makefile to help development.

* `make lint` - lints the code, exits with non-zero exit code if errors are found.
* `make test` - runs the tests against the code.