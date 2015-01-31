// Package test contains utilities to test topo.Server
// implementations. If you are testing your implementation, you will
// want to call CheckAll in your test method. For an example, look at
// the tests in github.com/youtube/vitess/go/vt/zktopo.
package test

import (
	"testing"

	"github.com/youtube/vitess/go/vt/topo"
)

// CheckShardReplication tests ShardReplication objects
func CheckShardReplication(t *testing.T, ts topo.Server) {
	cell := getLocalCell(t, ts)
	if _, err := ts.GetShardReplication(cell, "test_keyspace", "-10"); err != topo.ErrNoNode {
		t.Errorf("GetShardReplication(not there): %v", err)
	}

	sr := &topo.ShardReplication{
		ReplicationLinks: []topo.ReplicationLink{
			topo.ReplicationLink{
				TabletAlias: topo.TabletAlias{
					Cell: "c1",
					Uid:  1,
				},
			},
		},
	}
	if err := ts.UpdateShardReplicationFields(cell, "test_keyspace", "-10", func(oldSr *topo.ShardReplication) error {
		*oldSr = *sr
		return nil
	}); err != nil {
		t.Fatalf("UpdateShardReplicationFields() failed: %v", err)
	}

	if sri, err := ts.GetShardReplication(cell, "test_keyspace", "-10"); err != nil {
		t.Errorf("GetShardReplication(new guy) failed: %v", err)
	} else {
		if len(sri.ReplicationLinks) != 1 ||
			sri.ReplicationLinks[0].TabletAlias.Cell != "c1" ||
			sri.ReplicationLinks[0].TabletAlias.Uid != 1 {
			t.Errorf("GetShardReplication(new guy) returned wrong value: %v", *sri)
		}
	}

	if err := ts.UpdateShardReplicationFields(cell, "test_keyspace", "-10", func(sr *topo.ShardReplication) error {
		sr.ReplicationLinks = append(sr.ReplicationLinks, topo.ReplicationLink{
			TabletAlias: topo.TabletAlias{
				Cell: "c3",
				Uid:  3,
			},
		})
		return nil
	}); err != nil {
		t.Errorf("UpdateShardReplicationFields() failed: %v", err)
	}

	if sri, err := ts.GetShardReplication(cell, "test_keyspace", "-10"); err != nil {
		t.Errorf("GetShardReplication(after append) failed: %v", err)
	} else {
		if len(sri.ReplicationLinks) != 2 ||
			sri.ReplicationLinks[0].TabletAlias.Cell != "c1" ||
			sri.ReplicationLinks[0].TabletAlias.Uid != 1 ||
			sri.ReplicationLinks[1].TabletAlias.Cell != "c3" ||
			sri.ReplicationLinks[1].TabletAlias.Uid != 3 {
			t.Errorf("GetShardReplication(new guy) returned wrong value: %v", *sri)
		}
	}

	if err := ts.DeleteShardReplication(cell, "test_keyspace", "-10"); err != nil {
		t.Errorf("DeleteShardReplication(existing) failed: %v", err)
	}
	if err := ts.DeleteShardReplication(cell, "test_keyspace", "-10"); err != topo.ErrNoNode {
		t.Errorf("DeleteShardReplication(again) returned: %v", err)
	}
}
