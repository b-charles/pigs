module github.com/b-charles/pigs/log

go 1.18

require (
	github.com/b-charles/pigs/ioc v0.3.0
	github.com/b-charles/pigs/json v0.0.0-20220819174236-91e620daa867
	github.com/b-charles/pigs/smartconfig v0.0.0-00010101000000-000000000000
)

require (
	github.com/b-charles/pigs/config v0.3.0 // indirect
	github.com/benbjohnson/clock v1.3.0 // indirect
	github.com/spf13/afero v1.9.3 // indirect
	golang.org/x/sys v0.0.0-20220829200755-d48e67d00261 // indirect
	golang.org/x/text v0.3.7 // indirect
)

replace github.com/b-charles/pigs/ioc => ../ioc

replace github.com/b-charles/pigs/json => ../json

replace github.com/b-charles/pigs/config => ../config

replace github.com/b-charles/pigs/smartconfig => ../smartconfig
