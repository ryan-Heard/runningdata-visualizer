package main

// I could use reflect in this package to simplify the code
// but that would less effiecent as the data points grow

import (
	"net/http"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

var (
	running_rawdata = ReadCSVtoRunData("Activities.csv")
	scatter3DColor  = []string{
		"#313695", "#4575b4", "#74add1", "#abd9e9", "#e0f3f8",
		"#fee090", "#fdae61", "#f46d43", "#d73027", "#a50026",
	}
)

func genScatter3dData() []opts.Chart3DData {
	dataset := make([]opts.Chart3DData, 0)

	for _, data := range running_rawdata {
		dataset = append(dataset, opts.Chart3DData{
			Name: data.Date,
			Value: []interface{}{
				int(data.AvgPace.Minutes()),
				int(data.Distance),
				int(data.AvgHR)},
			Label:     &opts.Label{Show: false, Position: "insidetop", Formatter: "{b}"},
			ItemStyle: &opts.ItemStyle{Color: "green"},
		})
	}
	return dataset
}

func scatter3DBase() *charts.Scatter3D {
	scatter3d := charts.NewScatter3D()
	scatter3d.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Running HR x Distance x Pace"}),
		charts.WithVisualMapOpts(opts.VisualMap{
			Calculable: true,
			Max:        100,
			InRange:    &opts.VisualMapInRange{Color: scatter3DColor},
		}),
		charts.WithXAxis3DOpts(opts.XAxis3D{Name: "Avg Pace (Mins)", Show: true}),
		charts.WithYAxis3DOpts(opts.YAxis3D{Name: "Distance", Show: true}),
		charts.WithZAxis3DOpts(opts.ZAxis3D{Name: "Avg HR", Show: true}),
	)

	scatter3d.AddSeries("scatter3d", genScatter3dData())
	return scatter3d
}

func Scatter3DRender(w http.ResponseWriter, _ *http.Request) {
	page := components.NewPage()
	page.AddCharts(
		scatter3DBase(),
	)

	page.Render(w)
}
