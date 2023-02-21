module github.com/b-charles/pigs/smartconfig

go 1.18

require (
	github.com/b-charles/pigs/config v1.0.0
	github.com/b-charles/pigs/ioc v1.0.0
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.20.0
)

require (
	github.com/b-charles/pigs/json v1.0.0 // indirect
	github.com/benbjohnson/clock v1.3.0 // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/spf13/afero v1.9.3 // indirect
	golang.org/x/net v0.0.0-20220425223048-2871e0cb64e4 // indirect
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f // indirect
	golang.org/x/text v0.7.0 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/b-charles/pigs/ioc => ../ioc

replace github.com/b-charles/pigs/config => ../config
