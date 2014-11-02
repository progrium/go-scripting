implements(__module__, "ProgramObserver")

function ProgramStarted() {
	println("module2 got started")
}

function ProgramFinished() {
	println("module2 got finished")
}
