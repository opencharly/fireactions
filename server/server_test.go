package server

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/containerd/containerd"
	"github.com/stretchr/testify/assert"
)

// TestVerifySnapshotterServedAgainstLiveContainerd crosses the REAL containerd boundary and
// asserts both directions of the startup probe: a snapshotter the daemon serves is accepted, and
// a name it cannot serve is reported naming the configured value.
//
// It never mocks containerd — a mocked boundary proves nothing about the real one. Where no
// socket is present (a plain CI runner, a host without containerd) or the daemon is not reachable
// for this user, the test SKIPS VISIBLY with the reason, which is the honest outcome for a
// boundary it cannot cross.
func TestVerifySnapshotterServedAgainstLiveContainerd(t *testing.T) {
	const address = "/run/containerd/containerd.sock"

	_, statErr := os.Stat(address)
	if statErr != nil {
		t.Skipf("skipping: no containerd socket at %s (%v)", address, statErr)
	}

	// The client needs a default namespace: containerd answers a snapshot request without one
	// with "namespace is required: failed precondition", which the probe would otherwise report
	// as an unserved snapshotter. The server always sets one (config.Containerd.Namespace).
	client, err := containerd.New(address,
		containerd.WithTimeout(5*time.Second),
		containerd.WithDefaultNamespace("default"))
	if err != nil {
		t.Skipf("skipping: containerd at %s is not reachable for this user (%v)", address, err)
	}

	defer func() {
		_ = client.Close()
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// The native snapshotter is built into containerd, so it is always served.
	assert.NoError(t, verifySnapshotterServed(ctx, client, "native"),
		"a snapshotter containerd serves must pass the probe")

	err = verifySnapshotterServed(ctx, client, "fireactions-no-such-snapshotter")
	assert.Error(t, err, "a snapshotter containerd cannot serve must be reported, not swallowed")
	assert.ErrorContains(t, err, "not served", "the error must say the snapshotter is not served")
	assert.ErrorContains(t, err, "fireactions-no-such-snapshotter",
		"the error must name the configured value")
}
