package mumble

import (
	"testing"

	"github.com/dchote/go-mumble-server/internal/channel"
)

// A channel an operator marks require_tunnel must force voice over the TLS
// tunnel — channelRequiresTunnel returns true — even when no client is present
// and the crypto is not mixed, so SendAudio and the UDP ping echo take the
// tunnel path. Toggling the flag takes effect after RefreshChannelForceTunnel,
// which the REST/mutation path calls, without a restart.
func TestChannelRequiresTunnel_ForceFlag(t *testing.T) {
	srv := newACLTestServer(t)
	root := srv.chans.GetTree()[0]

	ch := srv.chans.Create(root.ID, "vault", "", 0, false, 0)
	if ch == nil {
		t.Fatal("create channel")
	}
	srv.RefreshChannelForceTunnel()
	if srv.channelRequiresTunnel(ch.ID) {
		t.Fatal("a channel with no force flag and no users must not require the tunnel")
	}

	yes := true
	if !srv.chans.Update(ch.ID, channel.UpdateOpts{RequireTunnel: &yes}) {
		t.Fatal("set require_tunnel")
	}
	srv.RefreshChannelForceTunnel()
	if !srv.channelRequiresTunnel(ch.ID) {
		t.Fatal("a force-tunnel channel must require the tunnel regardless of crypto tier")
	}

	no := false
	if !srv.chans.Update(ch.ID, channel.UpdateOpts{RequireTunnel: &no}) {
		t.Fatal("clear require_tunnel")
	}
	srv.RefreshChannelForceTunnel()
	if srv.channelRequiresTunnel(ch.ID) {
		t.Fatal("clearing require_tunnel must drop the tunnel requirement after refresh")
	}
}

// require_tunnel persists across a manager reload (it is a stored channel
// attribute, not runtime-only state), so a server restart keeps the policy.
func TestChannelRequireTunnel_PersistsAcrossReload(t *testing.T) {
	srv := newACLTestServer(t)
	root := srv.chans.GetTree()[0]
	ch := srv.chans.Create(root.ID, "sensitive", "", 0, false, 0)
	if ch == nil {
		t.Fatal("create channel")
	}
	yes := true
	if !srv.chans.Update(ch.ID, channel.UpdateOpts{RequireTunnel: &yes}) {
		t.Fatal("set require_tunnel")
	}

	// Reload rebuilds the in-memory tree from the database.
	srv.chans.Reload()
	reloaded, ok := srv.chans.GetChannel(ch.ID)
	if !ok {
		t.Fatal("channel missing after reload")
	}
	if !reloaded.RequireTunnel {
		t.Fatal("require_tunnel must survive a manager reload (persisted attribute)")
	}
	srv.RefreshChannelForceTunnel()
	if !srv.channelRequiresTunnel(ch.ID) {
		t.Fatal("force-tunnel must still be enforced after reload + refresh")
	}
}
