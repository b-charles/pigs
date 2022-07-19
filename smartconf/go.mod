module github.com/b-charles/pigs/smartconf

go 1.18

require (
	github.com/b-charles/pigs/config v0.3.0
	github.com/b-charles/pigs/ioc v0.3.0
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.19.0
)

require (
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	golang.org/x/net v0.0.0-20220225172249-27dd8689420f // indirect
	golang.org/x/sys v0.0.0-20220319134239-a9b59b0215f8 // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/b-charles/pigs/ioc => ../ioc

replace github.com/b-charles/pigs/config => ../config
