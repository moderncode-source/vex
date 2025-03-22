module github.com/moderncode-source/vex-svc

go 1.23.3

require github.com/moderncode-source/vex-svc/cmd v0.0.1

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/moderncode-source/vex-svc/vex v0.0.1 // indirect
	github.com/rs/zerolog v1.34.0 // indirect
	github.com/spf13/cobra v1.9.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	golang.org/x/net v0.37.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
)

replace github.com/moderncode-source/vex-svc/cmd v0.0.1 => ./cmd

replace github.com/moderncode-source/vex-svc/vex v0.0.1 => ./vex
