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

package holtwinters_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jthomperoo/holtwinters"
)

func TestPredict(t *testing.T) {
	equateErrorMessage := cmp.Comparer(func(x, y error) bool {
		if x == nil || y == nil {
			return x == nil && y == nil
		}
		return x.Error() == y.Error()
	})

	var tests = []struct {
		description      string
		expected         []float64
		expectedErr      error
		series           []float64
		seasonLength     int
		alpha            float64
		beta             float64
		gamma            float64
		predictionLength int
	}{
		{
			"Fail, season length too short",
			nil,
			errors.New(`Invalid parameter for prediction; season length must be at least 2, is 1`),
			[]float64{1, 2, 3, 2, 1},
			1,
			0.9,
			0.9,
			0.9,
			3,
		},
		{
			"Fail, negative prediction length",
			nil,
			errors.New(`Invalid parameter for prediction; prediction length must be at least 0, cannot be negative, is -3`),
			[]float64{1, 2, 3, 2, 1},
			5,
			0.9,
			0.9,
			0.9,
			-3,
		},
		{
			"Fail, alpha too high",
			nil,
			errors.New(`Invalid parameter for prediction; alpha must be between 0 and 1, is 1.500000`),
			[]float64{1, 2, 3, 2, 1},
			5,
			1.5,
			0.9,
			0.9,
			3,
		},
		{
			"Fail, alpha too low",
			nil,
			errors.New(`Invalid parameter for prediction; alpha must be between 0 and 1, is -0.200000`),
			[]float64{1, 2, 3, 2, 1},
			5,
			-0.2,
			0.9,
			0.9,
			3,
		},
		{
			"Fail, beta too high",
			nil,
			errors.New(`Invalid parameter for prediction; beta must be between 0 and 1, is 2.300000`),
			[]float64{1, 2, 3, 2, 1},
			5,
			0.9,
			2.3,
			0.9,
			3,
		},
		{
			"Fail, beta too low",
			nil,
			errors.New(`Invalid parameter for prediction; beta must be between 0 and 1, is -5.000000`),
			[]float64{1, 2, 3, 2, 1},
			5,
			0.9,
			-5,
			0.9,
			3,
		},
		{
			"Fail, gamma too high",
			nil,
			errors.New(`Invalid parameter for prediction; gamma must be between 0 and 1, is 30.000000`),
			[]float64{1, 2, 3, 2, 1},
			5,
			0.9,
			0.9,
			30,
			3,
		},
		{
			"Fail, gamma too low",
			nil,
			errors.New(`Invalid parameter for prediction; gamma must be between 0 and 1, is -20.000000`),
			[]float64{1, 2, 3, 2, 1},
			5,
			0.9,
			0.9,
			-20,
			3,
		},
		{
			"Fail, data provided less than full season",
			nil,
			errors.New(`Invalid parameter for prediction; must have at least 1 season of data to predict, season length: 5, series length: 3`),
			[]float64{1, 2, 3},
			5,
			0.9,
			0.9,
			0.9,
			5,
		},
		{
			"Success, 1 season, no prediction",
			[]float64{1, 2.8400000000000003, 3.1516, 1.959964, 0.9732295599999999},
			nil,
			[]float64{1, 2, 3, 2, 1},
			5,
			0.9,
			0.9,
			0.9,
			0,
		},
		{
			"Success, 1 and a half seasons data",
			[]float64{1, 2.8400000000000003, 3.1516, 1.959964, 0.9732295599999999, 1.1762401724, 1.7801866939959998, 3.2631861240188407, 2.2876973806214185, 1.4767954296844166, 1.6533669041674148, 2.7683539322222135, 3.930203928270833},
			nil,
			[]float64{1, 2, 3, 2, 1, 1.1, 1.9, 3.1},
			5,
			0.9,
			0.9,
			0.9,
			5,
		},
		{
			"Success, less than 2 seasons data",
			[]float64{1, 2.8400000000000003, 3.1516, 1.959964, 0.9732295599999999, 0.971479762, 1.926903744, 2.8411077259999997, 1.8711579079999998, 0.866925488},
			nil,
			[]float64{1, 2, 3, 2, 1},
			5,
			0.9,
			0.9,
			0.9,
			5,
		},
		{
			"Success, 2 seasons data",
			[]float64{
				1, 2.7064000000000004, 3.132456, 1.96677224, 0.9771183496000001, 1.1766870973840002, 1.7830314232813598, 3.2515613630131943,
				2.1199062313456905, 1.0747739825249312, 1.0894589192483668, 2.0086996332729483, 2.991675122285811, 1.967955201522516, 0.9716977015641067,
			},
			nil,
			[]float64{1, 2, 3, 2, 1, 1.1, 1.9, 3.1, 2.1, 1.1},
			5,
			0.9,
			0.9,
			0.9,
			5,
		},
		{
			"Success, more than 2 seasons data",
			[]float64{30, 20.34449316666667, 28.410051892109554, 30.438122252647577, 39.466817731253066, 47.54961891047195, 52.52339682497974, 46.53453460769274,
				36.558407328055765, 38.56283307754578, 30.51864332437879, 28.425963657825292, 16.30247725646635, 8.228588857142476, 19.30036874234319, 23.38657154193773,
				26.323990741396006, 34.356648660113095, 40.36971459184453, 37.44298129818558, 26.469996240541015, 30.51819842804787, 26.580158132275145, 25.556750355604414,
				20.59232938487544, 12.557525846506284, 20.536167580315634, 17.449559582909338, 32.589947392978274, 34.559067611499714, 39.524706984702796, 35.54354494552727,
				21.507741573047714, 23.48782855767762, 20.541994359470845, 19.543228201110367, 16.60700323688017, 13.697607405158983, 16.621224546074888, 18.619564648649416,
				25.57626419227017, 28.544672577127326, 39.62603432821338, 30.578678843303678, 19.58514452366992, 23.614663453052163, 17.606991212001635, 25.767260902774442,
				16.759148937441683, 8.712803906763776, 16.72824428057732, 20.7768592516643, 27.760289930117256, 31.74794281311134, 45.85701109377136, 32.77988806685826,
				22.769367642515853, 27.80450001645962, 21.806956583618057, 26.862261134868607, 17.863888132693965, 7.79136434612686, 16.79511449881349, 20.831653319362697,
				30.885227379775543, 33.87620406969448, 43.8722204956629, 37.93866311702782, 31.017079798498486, 29.952760178336057, 25.95873287479028, 32.01973275816115,
				22.42511411230803, 15.343371755223066, 24.14282581581347, 27.02259921391996, 35.31139046245393, 38.999014669337356, 49.243283875692654, 40.84636009563803,
				31.205180503707012, 32.96259980122959, 28.5164783238384, 32.30616336737171, 22.737583867810464, 15.655841510725496, 24.4552955713159, 27.33506896942239,
				35.62386021795636, 39.31148442483978, 49.55575363119508, 41.15882985114047, 31.517650259209443, 33.275069556732014, 28.82894807934083, 32.618633122874144},
			nil,
			[]float64{30, 21, 29, 31, 40, 48, 53, 47, 37, 39, 31, 29, 17, 9, 20, 24, 27, 35, 41, 38,
				27, 31, 27, 26, 21, 13, 21, 18, 33, 35, 40, 36, 22, 24, 21, 20, 17, 14, 17, 19,
				26, 29, 40, 31, 20, 24, 18, 26, 17, 9, 17, 21, 28, 32, 46, 33, 23, 28, 22, 27,
				18, 8, 17, 21, 31, 34, 44, 38, 31, 30, 26, 32},
			12,
			0.716,
			0.029,
			0.993,
			24,
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			prediction, err := holtwinters.Predict(test.series, test.seasonLength, test.alpha, test.beta, test.gamma, test.predictionLength)

			if !cmp.Equal(&err, &test.expectedErr, equateErrorMessage) {
				t.Errorf("Error mismatch (-want +got):\n%s", cmp.Diff(test.expectedErr, err, equateErrorMessage))
				return
			}

			if !cmp.Equal(test.expected, prediction) {
				t.Errorf("prediction mismatch (-want +got):\n%s", cmp.Diff(test.expected, prediction))
			}
		})
	}

}
