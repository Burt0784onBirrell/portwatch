// Package jitter provides randomised interval spreading for the portwatch
// scan loop.
//
// When several portwatch processes run concurrently (e.g. in a Kubernetes
// DaemonSet) they would otherwise all wake up at the same moment and hammer
// the kernel's netstat interface together.  Adding a small random offset to
// each sleep interval distributes the load evenly over time.
//
// Usage:
//
//	j := jitter.New(0.2)          // up to 20 % spread
//	timer := j.Reset(5 * time.Second)
//	defer timer.Stop()
//	<-timer.C
package jitter
