package logwrap

import (
	"context"
	"sync/atomic"
)

const SegmentField = "segment"
const SegmentIdField = "segmentId"
const ParentSegmentIdField = "parentSegmentId"

const SegmentStartValue = "start"
const SegmentEndValue = "end"

const contextKeySegmentId = "_ShimmeringBeeLogSegmentId"

// Segment is used to wrap a section of a program, this can be used to demonstrate a group of logs are related. Segments
// can be nested and the `segmentId` and `parentSegmentId` fields can be used to reconstruct nested logs into a hierarchy.
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
// This code would product log likes aproximately like:
// * [INFO] api submission {"segment": "start", "segmentId": 1}
// * [INFO] perpare api submission {"segment": "start", "segmentId": 2, "parentSegmentId": 1}
// * [INFO] preparation results {"segmentId": 2, "parentSegmentId": 1, "request": <request object>}
// * [INFO] perpare api submission {"segment": "end", "segmentId": 2, "parentSegmentId": 1}
// * [INFO] api submission {"segment": "end", "segmentId": 1}
func (l Logger) Segment(pctx context.Context, message string, options ...Option) (context.Context, func()) {
	if parentSegmentId, present := l.getSegmentIdFromContext(pctx); present {
		options = append(options, Field(ParentSegmentIdField, parentSegmentId))
	}

	segmentId := atomic.AddUint64(l.segmentId, 1)
	options = append(options, Field(SegmentIdField, segmentId))

	ctx := l.AddOptionsToContext(pctx, options...)
	ctx = context.WithValue(ctx, l.contextKey(contextKeySegmentId), segmentId)

	l.Log(ctx, message, Field(SegmentField, SegmentStartValue))

	return ctx, func() {
		l.Log(ctx, message, Field(SegmentField, SegmentEndValue))
	}
}

func (l Logger) getSegmentIdFromContext(ctx context.Context) (uint64, bool) {
	if uncast := ctx.Value(l.contextKey(contextKeySegmentId)); uncast != nil {
		if segmentId, ok := uncast.(uint64); ok {
			return segmentId, true
		}
	}

	return 0, false
}
