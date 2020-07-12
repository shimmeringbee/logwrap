package logwrap

import (
	"runtime"
	"strings"
)

// SourceLocation is a representation of where the log line was called from.
type SourceLocation struct {
	Function string
	File     string
	Line     int
}

const SourceTraceField = "sourceTrace"
const maximumInternalFrameDepth = 32
const entrypointFunctionPrefix = "github.com/shimmeringbee/logwrap.Logger"

// SourceTrace inserts the file and line number of the file that Log was called at. This routine searches for the frame
// pointer before the first call to the Logger object. The search is used rather than static in case the option is used
// in a post option filter, and to allow utility wrappers such as LogInfo instead of `Log(...,...,Level(Info))`.
func SourceTrace(message *Message) {
	programCounters := make([]uintptr, maximumInternalFrameDepth)
	n := runtime.Callers(0, programCounters)

	location := SourceLocation{
		Function: "unknown",
		File:     "unknown",
		Line:     0,
	}

	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])
		acceptNew := true

		for more, frameIndex := true, 0; more && frameIndex <= maximumInternalFrameDepth; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()

			if strings.HasPrefix(frameCandidate.Function, entrypointFunctionPrefix) {
				acceptNew = true
			} else if acceptNew {
				location.Function = frameCandidate.Function
				location.File = frameCandidate.File
				location.Line = frameCandidate.Line
				acceptNew = false
			}
		}
	}

	message.Fields[SourceTraceField] = location
}
