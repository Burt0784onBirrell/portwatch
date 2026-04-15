package metrics

import (
	"context"
	"errors"
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func makePortSet(ports ...scanner.Port) scanner.PortSet {
	return scanner.PortSetFromSlice(ports)
}

func TestWithMetrics_RecordsScanOnSuccess(t *testing.T) {
	col := New()
	called := false

	wrapped := WithMetrics(func(ctx context.Context) (scanner.PortSet, error) {
		called = true
		return makePortSet(), nil
	}, col)

	_, err := wrapped(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected inner scan function to be called")
	}

	snap := col.Snapshot()
	if snap.TotalScans != 1 {
		t.Errorf("expected TotalScans=1, got %d", snap.TotalScans)
	}
}

func TestWithMetrics_DoesNotRecordOnError(t *testing.T) {
	col := New()

	wrapped := WithMetrics(func(ctx context.Context) (scanner.PortSet, error) {
		return nil, errors.New("scan failed")
	}, col)

	_, err := wrapped(context.Background())
	if err == nil {
		t.Fatal("expected error to be propagated")
	}

	snap := col.Snapshot()
	if snap.TotalScans != 0 {
		t.Errorf("expected TotalScans=0 on error, got %d", snap.TotalScans)
	}
}

func TestWithMetrics_AccumulatesMultipleScans(t *testing.T) {
	col := New()

	wrapped := WithMetrics(func(ctx context.Context) (scanner.PortSet, error) {
		return makePortSet(), nil
	}, col)

	const numScans = 5
	for i := 0; i < numScans; i++ {
		if _, err := wrapped(context.Background()); err != nil {
			t.Fatalf("unexpected error on scan %d: %v", i, err)
		}
	}

	snap := col.Snapshot()
	if snap.TotalScans != numScans {
		t.Errorf("expected TotalScans=%d, got %d", numScans, snap.TotalScans)
	}
}

func TestRecordDiff_UpdatesAlerts(t *testing.T) {
	col := New()
	RecordDiff(col, 3)

	snap := col.Snapshot()
	if snap.TotalAlerts != 3 {
		t.Errorf("expected TotalAlerts=3, got %d", snap.TotalAlerts)
	}
}

func TestRecordDiff_ZeroIsNoop(t *testing.T) {
	col := New()
	RecordDiff(col, 0)

	snap := col.Snapshot()
	if snap.TotalAlerts != 0 {
		t.Errorf("expected TotalAlerts=0 for zero diff, got %d", snap.TotalAlerts)
	}
}
