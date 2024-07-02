package scanner

import "fmt"

func nesto() {
	var syncCmdParams = struct {
		SourceStorePath string
		TargetStorePath string
		SyncJobPath     string
	}{}

	fmt.Println(syncCmdParams.SyncJobPath)
}
