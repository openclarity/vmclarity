## Pluggable scanner - Golang example

You only have to implement the files under `scanner` dir. 
CIS Docker benchmark was taken as an example.

Scanner main module is not done as it's importing Golang REST server code rather than running REST server (which is okay for now and provides the same experience).
Ideally, we want to change that to client stub and replace the `server.go` to only use the client itself with reporting mechanism.
The reason for the current approach is that I wanted to implement the job orchestrator/executor so that we can have a simple interface
and be able to implement new scanners quickly.
