module github.com/0dayfall/generationk

go 1.15

require (
	github.com/0dayfall/generationk/indicators v0.0.0-00010101000000-000000000000
	github.com/go-delve/delve v1.5.1 // indirect
	github.com/pkg/profile v1.5.0
	github.com/rs/zerolog v1.20.0
	github.com/sirupsen/logrus v1.7.0
)

replace github.com/0dayfall/generationk/indicators => ./indicators
