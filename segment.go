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
//      logger.Log(subCtx, "preparation results", Field("request": request))
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
		options = append(options, Field(ParentSegmentIDField, parentSegmentID))
	}

	segmentID := atomic.AddUint64(l.segmentID, 1)
	options = append(options, Field(SegmentIDField, segmentID))

	ctx := l.AddOptionsToContext(pctx, options...)
	ctx = context.WithValue(ctx, l.contextKey(contextKeySegmentID), segmentID)

	l.Log(ctx, message, Field(SegmentField, SegmentStartValue))

	return ctx, func() {
		l.Log(ctx, message, Field(SegmentField, SegmentEndValue))
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
