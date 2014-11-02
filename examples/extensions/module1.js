implements(__module__, "ProgramObserver")

function ProgramStarted() {
	println(__module__ + " got started")
}

function ProgramFinished() {
	println(__module__ + " got finished")
}
