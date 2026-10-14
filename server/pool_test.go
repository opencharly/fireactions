package server

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/containerd/containerd/mount"
	"github.com/stretchr/testify/assert"
)

// TestValidateRootDriveMounts exercises every branch of the guard that decides whether a
// prepared snapshot can be attached to a microVM. Firecracker attaches a block device or a
// FILE (devmapper yields the former, blockfile the latter); a directory can never work, and a
// mount that has gone away must be reported by path rather than booting an unusable guest.
func TestValidateRootDriveMounts(t *testing.T) {
	dir := t.TempDir()
	rootfs := filepath.Join(dir, "rootfs.ext4")
	assert.NoError(t, os.WriteFile(rootfs, []byte("not really ext4"), 0o600))
	gone := filepath.Join(dir, "gone.ext4")

	tests := []struct {
		name    string
		mounts  []mount.Mount
		wantErr string
	}{
		{
			name:    "no mounts at all is refused",
			mounts:  nil,
			wantErr: "returned no mounts",
		},
		{
			name:    "a directory mount is refused",
			mounts:  []mount.Mount{{Source: dir}},
			wantErr: "cannot back a microVM root drive",
		},
		{
			name:    "a mount that has gone away is reported by path",
			mounts:  []mount.Mount{{Source: gone}},
			wantErr: "stat mount",
		},
		{
			name:   "a file mount is accepted (the blockfile shape)",
			mounts: []mount.Mount{{Source: rootfs}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRootDriveMounts(SnapshotterBlockfile, "snapshot-id", tt.mounts)

			if tt.wantErr == "" {
				assert.NoError(t, err)

				return
			}

			assert.ErrorContains(t, err, tt.wantErr)
			assert.ErrorContains(t, err, "snapshot-id", "the error must name the snapshot it is about")
		})
	}
}
