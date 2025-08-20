package main

import (
	"fmt"
	"os"
	"time"

	"github.com/NCCloud/mayfly/pkg/common"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/pkg/browser"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	config        *common.Config
	mgrClient     client.Client
	pageTitle     = "Mayfly Benchmark"
	benchmarkHtml = "mayfly_benchmark.html"
)

func init() {
	var (
		scheme       = runtime.NewScheme()
		mgrClientErr error
	)

	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	config = common.NewConfig()

	mgrClient, mgrClientErr = client.New(ctrl.GetConfigOrDie(), client.Options{Scheme: scheme})
	if mgrClientErr != nil {
		panic(mgrClientErr)
	}
}

func main() {
	itemCount := 1000
	benchmark := NewBenchmark(mgrClient, config, itemCount).Delay(0)
	benchmark.Start()

	benchmarkFile, _ := os.Create(benchmarkHtml)

	_, writeStringErr := benchmarkFile.WriteString(" <meta http-equiv=\"refresh\" content=\"6\" />\nWaiting for Mayfly benchmark to start...")
	if writeStringErr != nil {
		panic(writeStringErr)
	}

	openFilErr := browser.OpenFile(benchmarkFile.Name())
	if openFilErr != nil {
		panic(openFilErr)
	}

	for {
		var (
			waitTimeSecond = 5
			result         = benchmark.GetResult()
			durations      []string
			data           = map[string][]opts.LineData{}
		)

		time.Sleep(time.Duration(waitTimeSecond) * time.Second)

		for _, point := range result.Points {
			durations = append(durations, fmt.Sprintf("%.0fs", point.time.Sub(result.StartedAt).Seconds()))

			for _, resource := range config.Resources {
				data[resource] = append(data[resource], opts.LineData{Name: resource, Value: point.kind[resource]})
			}
		}

		Render(CreateChart(durations, data))
	}
}

func Render(chart *charts.Line) {
	benchFile, _ := os.Create(benchmarkHtml)

	_, writeStringErr := benchFile.WriteString(" <meta http-equiv=\"refresh\" content=\"6\" />")
	if writeStringErr != nil {
		panic(writeStringErr)
	}

	renderErr := chart.Render(benchFile)
	if renderErr != nil {
		panic(renderErr)
	}
}

func CreateChart(xAxis []string, yzAxis map[string][]opts.LineData) *charts.Line {
	var chartAreaOpacity float32 = 0.2

	line := charts.NewLine()
	line.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title: pageTitle,
	}), charts.WithLegendOpts(opts.Legend{
		Show: opts.Bool(true),
	}), charts.WithInitializationOpts(opts.Initialization{
		PageTitle: pageTitle,
		Width:     "1600px",
		Height:    "800px",
		Theme:     "infographic",
	}))

	line.SetXAxis(xAxis)

	for kind, value := range yzAxis {
		line.AddSeries(kind, value).SetSeriesOptions(
			charts.WithAreaStyleOpts(opts.AreaStyle{
				Opacity: opts.Float(chartAreaOpacity),
			}),
			charts.WithLineChartOpts(opts.LineChart{
				Smooth: opts.Bool(false),
			}),
			charts.WithMarkPointStyleOpts(
				opts.MarkPointStyle{Label: &opts.Label{Show: opts.Bool(true)}}),
		)
	}

	return line
}
