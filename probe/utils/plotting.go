// Copyright © 2023 SUSE LLC
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"errors"

	"github.com/montanaflynn/stats"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func GenerateRawDataPlot(data RestartRelatedData,
	mark string,
	timeUnit string,
	fileName string,
	genTS string) (string, error) {

	if restartEvents, hit := data[mark]; hit {
		p := plot.New()
		p.Add(plotter.NewGrid())

		p.Title.Text = "Raw Data Measures: " + mark
		p.X.Label.Text = "Restart ID"
		p.Y.Label.Text = "Duration: " + timeUnit

		_,
			evtSeriesMainData,
			evtSeriesFrontedUpData,
			evtSeriesFUpMainDelta := GetSplitDataForSingleRestartRelatedData(restartEvents, StrTimeUnit2TimeUnit[timeUnit])

		//main

		{
			mainPts := make(plotter.XYs, len(restartEvents))
			for i := range mainPts {
				mainPts[i].X = float64(i)
				mainPts[i].Y = evtSeriesMainData[i]
			}

			lpLineMain, lpPointsMain, err := plotter.NewLinePoints(mainPts)
			if err != nil {
				Logger.Error("GeneratePlot-to-main: NewLinePoints", err.Error())
				return "", err
			}
			lpLineMain.Color = plotutil.Color(0)
			lpPointsMain.Shape = draw.PyramidGlyph{}
			lpPointsMain.Color = plotutil.Color(0)

			p.Add(lpLineMain, lpPointsMain)
			p.Legend.Add("to-main", lpLineMain, lpPointsMain)
		}

		//fronted-up

		{
			frontendUpPts := make(plotter.XYs, len(restartEvents))
			for i := range frontendUpPts {
				frontendUpPts[i].X = float64(i)
				frontendUpPts[i].Y = evtSeriesFrontedUpData[i]
			}

			lpLineFUp, lpPointsFUp, err := plotter.NewLinePoints(frontendUpPts)
			if err != nil {
				Logger.Error("GeneratePlot-to-fronted-up: NewLinePoints", err.Error())
				return "", err
			}
			lpLineFUp.Color = plotutil.Color(1)
			lpPointsFUp.Shape = draw.PyramidGlyph{}
			lpPointsFUp.Color = plotutil.Color(1)

			p.Add(lpLineFUp, lpPointsFUp)
			p.Legend.Add("to-fronted-up", lpLineFUp, lpPointsFUp)
		}

		//delta

		{
			fUpMainDeltaPts := make(plotter.XYs, len(restartEvents))
			for i := range fUpMainDeltaPts {
				fUpMainDeltaPts[i].X = float64(i)
				fUpMainDeltaPts[i].Y = evtSeriesFUpMainDelta[i]
			}

			lpLineFUpD, lpPointsFUpD, err := plotter.NewLinePoints(fUpMainDeltaPts)
			if err != nil {
				Logger.Error("GeneratePlot-fronted-up-main-delta: NewLinePoints", err.Error())
				return "", err
			}
			lpLineFUpD.Color = plotutil.Color(2)
			lpPointsFUpD.Shape = draw.PyramidGlyph{}
			lpPointsFUpD.Color = plotutil.Color(2)

			p.Add(lpLineFUpD, lpPointsFUpD)
			p.Legend.Add("fronted-up-main-delta", lpLineFUpD, lpPointsFUpD)
		}

		fName := fileName + "_" + "raw" + "_" + genTS + ".svg"

		if err := p.Save(30*vg.Centimeter, 20*vg.Centimeter, fName); err != nil {
			Logger.Error("GeneratePlot: Saving plot:", fName)
			return "", err
		}

		return fName, nil

	} else {
		Logger.Error("GeneratePlot: no series with mark:", mark)
		return "", errors.New("no series with mark")
	}
}

func GeneratePercentilesPlot(data RestartRelatedData,
	mark string,
	timeUnit string,
	fileName string,
	genTS string) (string, string, string, error) {
	if restartEvents, hit := data[mark]; hit {
		_,
			evtSeriesMainData,
			evtSeriesFrontedUpData,
			evtSeriesFUpMainDelta := GetSplitDataForSingleRestartRelatedData(restartEvents, StrTimeUnit2TimeUnit[timeUnit])

		//main

		fNameMain := fileName + "_percentiles_to_main_" + genTS + ".svg"

		{
			p := plot.New()
			p.Add(plotter.NewGrid())

			p.Title.Text = "Percentiles - Main (Nearest Rank): " + mark
			p.X.Label.Text = "Percentile"
			p.Y.Label.Text = "Duration: " + timeUnit

			pts := make(plotter.XYs, 100)
			for i := 1; i <= 100; i++ {
				if val, err := stats.PercentileNearestRank(evtSeriesMainData, float64(i)); err == nil {
					pts[i-1].X = float64(i)
					pts[i-1].Y = val
				}
			}

			lpLine, lpPoints, err := plotter.NewLinePoints(pts)
			if err != nil {
				Logger.Error("GeneratePercentilesPlot-to-main: NewLinePoints", err.Error())
				return "", "", "", err
			}
			lpLine.Color = plotutil.Color(0)
			lpPoints.Shape = draw.PyramidGlyph{}
			lpPoints.Color = plotutil.Color(0)

			p.Add(lpLine, lpPoints)
			p.Legend.Add("to-main", lpLine, lpPoints)

			if err := p.Save(30*vg.Centimeter, 20*vg.Centimeter, fNameMain); err != nil {
				Logger.Error("GeneratePlot: Saving plot:", fNameMain)
				return "", "", "", err
			}
		}

		//fronted-up

		fNameFUp := fileName + "_percentiles_to_fup_" + genTS + ".svg"

		{
			p := plot.New()
			p.Add(plotter.NewGrid())

			p.Title.Text = "Percentiles - FrontEndUp (Nearest Rank): " + mark
			p.X.Label.Text = "Percentile"
			p.Y.Label.Text = "Duration: " + timeUnit

			pts := make(plotter.XYs, 100)
			for i := 1; i <= 100; i++ {
				if val, err := stats.PercentileNearestRank(evtSeriesFrontedUpData, float64(i)); err == nil {
					pts[i-1].X = float64(i)
					pts[i-1].Y = val
				}
			}

			lpLine, lpPoints, err := plotter.NewLinePoints(pts)
			if err != nil {
				Logger.Error("GeneratePercentilesPlot-to-fronted-up: NewLinePoints", err.Error())
				return "", "", "", err
			}
			lpLine.Color = plotutil.Color(1)
			lpPoints.Shape = draw.PyramidGlyph{}
			lpPoints.Color = plotutil.Color(1)

			p.Add(lpLine, lpPoints)
			p.Legend.Add("to-fronted-up", lpLine, lpPoints)

			if err := p.Save(30*vg.Centimeter, 20*vg.Centimeter, fNameFUp); err != nil {
				Logger.Error("GeneratePlot: Saving plot:", fNameFUp)
				return "", "", "", err
			}
		}

		//fronted-up-main-delta

		fNameFUpD := fileName + "_percentiles_fup_main_delta_" + genTS + ".svg"

		{
			p := plot.New()
			p.Add(plotter.NewGrid())

			p.Title.Text = "Percentiles - FrontEndUp-Main-Delta (Nearest Rank): " + mark
			p.X.Label.Text = "Percentile"
			p.Y.Label.Text = "Duration: " + timeUnit

			pts := make(plotter.XYs, 100)
			for i := 1; i <= 100; i++ {
				if val, err := stats.PercentileNearestRank(evtSeriesFUpMainDelta, float64(i)); err == nil {
					pts[i-1].X = float64(i)
					pts[i-1].Y = val
				}
			}

			lpLine, lpPoints, err := plotter.NewLinePoints(pts)
			if err != nil {
				Logger.Error("GeneratePercentilesPlot-to-fronted-up: NewLinePoints", err.Error())
				return "", "", "", err
			}
			lpLine.Color = plotutil.Color(2)
			lpPoints.Shape = draw.PyramidGlyph{}
			lpPoints.Color = plotutil.Color(2)

			p.Add(lpLine, lpPoints)
			p.Legend.Add("to-fronted-up", lpLine, lpPoints)

			if err := p.Save(30*vg.Centimeter, 20*vg.Centimeter, fNameFUpD); err != nil {
				Logger.Error("GeneratePlot: Saving plot:", fNameFUpD)
				return "", "", "", err
			}
		}

		return fNameMain, fNameFUp, fNameFUpD, nil

	} else {
		Logger.Error("GeneratePercentilesPlot: no series with mark:", mark)
		return "", "", "", errors.New("no series with mark")
	}
}
