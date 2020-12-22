module github.com/0dayfall/generationk

go 1.15

require (
	github.com/go-echarts/go-echarts/v2 v2.2.3
	github.com/0dayfall/generationk/indicators v0.0.0-00010101000000-000000000000
	github.com/shiena/ansicolor v0.0.0-20200904210342-c7312218db18
	github.com/sirupsen/logrus v1.7.0
)

replace github.com/0dayfall/generationk/indicators => ./indicators
