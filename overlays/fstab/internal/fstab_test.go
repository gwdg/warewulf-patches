package fstab

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/warewulf/warewulf/internal/app/wwctl/overlay/show"
	"github.com/warewulf/warewulf/internal/pkg/config"
	"github.com/warewulf/warewulf/internal/pkg/testenv"
	"github.com/warewulf/warewulf/internal/pkg/wwlog"
)

func Test_fstabOverlay(t *testing.T) {
	env := testenv.New(t)
	defer env.RemoveAll(t)
	env.ImportFile(t, "etc/warewulf/nodes.conf", "nodes.conf")
	env.ImportFile(t, "etc/warewulf/warewulf.conf", "warewulf.conf")
	assert.NoError(t, config.Get().Read(env.GetPath("etc/warewulf/warewulf.conf")))
	env.ImportFile(t, "var/lib/warewulf/overlays/fstab/rootfs/etc/fstab.ww", "../rootfs/etc/fstab.ww")

	tests := []struct {
		name string
		args []string
		log  string
	}{
		{
			name: "/etc/fstab",
			args: []string{"--render", "node1", "fstab", "etc/fstab.ww"},
			log:  fstab,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := show.GetCommand()
			cmd.SetArgs(tt.args)
			stdout := bytes.NewBufferString("")
			stderr := bytes.NewBufferString("")
			logbuf := bytes.NewBufferString("")
			cmd.SetOut(stdout)
			cmd.SetErr(stderr)
			wwlog.SetLogWriter(logbuf)
			err := cmd.Execute()
			assert.NoError(t, err)
			assert.Empty(t, stdout.String())
			assert.Empty(t, stderr.String())
			assert.Equal(t, tt.log, logbuf.String())
		})
	}
}

const fstab string = `backupFile: true
writeFile: true
Filename: etc/fstab
# This file is autogenerated by warewulf
rootfs / tmpfs defaults 0 0
devpts /dev/pts devpts gid=5,mode=620 0 0
tmpfs /run/shm tmpfs defaults 0 0
sysfs /sys sysfs defaults 0 0
proc /proc proc defaults 0 0
# mounts for local file systems created with ignition in nodes.conf
# all with noauto as mounts happens with systemd units
/dev/disk/by-partlabel/scratch /scratch btrfs noauto,defaults 0 0
/dev/disk/by-partlabel/swap swap swap noauto,defaults 0 0
# nfs mounts provided in warewulf.conf
192.168.0.1:/home /home nfs defaults 0 0
`
