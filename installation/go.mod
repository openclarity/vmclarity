module github.com/openclarity/vmclarity/installation

go 1.22.4

require github.com/openclarity/vmclarity/utils v0.7.2

replace (
	github.com/openclarity/vmclarity/core => ../core
	github.com/openclarity/vmclarity/utils => ../utils
)
