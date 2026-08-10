module uabench

go 1.26.2

require (
	github.com/avct/uasurfer v0.0.0-20260805144106-c7341ba9111e
	github.com/bakhod1r/uax v0.0.0-20260807204757-ea91613811e6
	github.com/mileusna/useragent v1.3.5
	github.com/mssola/user_agent v0.6.0
	github.com/ua-parser/uap-go v0.0.0-20260529044130-17c35e68e58c
)

require (
	github.com/bakhod1r/alx v0.1.0 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/avct/uasurfer => github.com/LumenResearch/uasurfer v0.0.0-20260805144106-c7341ba9111e

replace github.com/bakhod1r/uax => ..
