package logwrap

import (
	"context"
	"sync/atomic"
)

// SegmentField is the name of the field which a segment will place the start and end markers in.
const SegmentField = "segment"

// SegmentIDField is the name of the field which contains the segment ID.
const SegmentIDField = "segmentID"

// ParentSegmentIDField is the name of the field which contains this segments parent ID.
const ParentSegmentIDField = "parentSegmentID"

// SegmentStartValue is the value placed in the SegmentField denoting the start of a segment.
const SegmentStartValue = "start"

// SegmentEndValue is the value placed in the SegmentField denoting the end of a segment.
const SegmentEndValue = "end"

const contextKeySegmentID = "_ShimmeringBeeLogSegmentID"

// Segment is used to wrap a section of a program, this can be used to demonstrate a group of logs are related. Segments
// can be nested and the `segmentID` and `parentSegmentId` fields can be used to reconstruct nested logs into a hierarchy.
//
// An expected use of Segment might be as follows:
//  func submitToAPI(pctx context.Context) {
//      ctx, end := logger.Segment(pctx, "api submission")
//      defer end()
//      //
//      subCtx, subEnd := logger.Segment(ctx, "prepare api submission")
//      request, err := // Prepare api submission
//      logger.Log(subCtx, "preparation results", Datum("request": request))
//      subEnd()
//      //
//      err := // Submit to api functional code
//      if err != nil {
//           logger.Log(ctx, "failed to submit to api", Err(err))
//      }
//  }
//
// This code would product log likes approximately like:
// * [INFO] api submission {"segment": "start", "segmentID": 1}
// * [INFO] prepare api submission {"segment": "start", "segmentID": 2, "parentSegmentId": 1}
// * [INFO] preparation results {"segmentID": 2, "parentSegmentId": 1, "request": <request object>}
// * [INFO] prepare api submission {"segment": "end", "segmentID": 2, "parentSegmentId": 1}
// * [INFO] api submission {"segment": "end", "segmentID": 1}
func (l Logger) Segment(pctx context.Context, message string, options ...Option) (context.Context, func()) {
	if parentSegmentID, present := l.getSegmentIDFromContext(pctx); present {
		options = append(options, Datum(ParentSegmentIDField, parentSegmentID))
	}

	segmentID := atomic.AddUint64(l.segmentID, 1)
	options = append(options, Datum(SegmentIDField, segmentID))

	ctx := l.AddOptionsToContext(pctx, options...)
	ctx = context.WithValue(ctx, l.contextKey(contextKeySegmentID), segmentID)

	l.Log(ctx, message, Datum(SegmentField, SegmentStartValue))

	return ctx, func() {
		l.Log(ctx, message, Datum(SegmentField, SegmentEndValue))
	}
}

//SegmentFn works similar to Segment, but returns a function that takes a new function to be called. This can be used
//to wrap calls with a Segment, removing the complexity of handling the end function of Segment. The function returned
//passes the new child context into the wrapped function.
//
// For example, converting the Segment function example:
//  func submitToAPI(pctx context.Context) {
//      ctx, end := logger.Segment(pctx, "api submission")
//      defer end()
//      //
//      if err := logger.SegmentFn(ctx, "prepare api submission")(func(subCtx context.Context) error {
//          request, err := // Prepare api submission
//          logger.Log(subCtx, "preparation results", Datum("request": request))
//      }); err != nil {
//           logger.Log(ctx, "failed to submit to api", Err(err))
//      }
//  }
//
//It is expected that errors in the returned function are actual problems, as it will log the error. It is not expected
//that segments will be used where the error is unimportant.
func (l Logger) SegmentFn(pctx context.Context, message string, options ...Option) func(func(ctx context.Context) error) error {
	return func(f func(ctx context.Context) error) error {
		c, done := l.Segment(pctx, message, options...)

		err := f(c)
		if err != nil {
			l.Error(c, "segment errored", Err(err))
		}

		done()
		return err
	}
}

func (l Logger) getSegmentIDFromContext(ctx context.Context) (uint64, bool) {
	if uncast := ctx.Value(l.contextKey(contextKeySegmentID)); uncast != nil {
		if segmentID, ok := uncast.(uint64); ok {
			return segmentID, true
		}
	}

	return 0, false
}
