module github.com/0dayfall/generationk/cmd

go 1.15

replace github.com/0dayfall/generationk/strategies => ../strategies

replace github.com/0dayfall/generationk => ../

require (
	github.com/0dayfall/generationk v0.0.0-00010101000000-000000000000
	github.com/0dayfall/generationk/strategies v0.0.0-00010101000000-000000000000
)
