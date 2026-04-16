// Package batch implements a size-and-time-bounded event accumulator.
//
// A Batcher collects alert.Event values and signals when they should be
// forwarded downstream. Flushing is triggered by whichever condition is
// met first:
//
//   - The number of buffered events reaches maxSize.
//   - The configured window duration has elapsed since the first event
//     in the current batch was received.
//
// Use Flush to drain the buffer explicitly (e.g. on daemon shutdown).
package batch
