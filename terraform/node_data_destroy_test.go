package terraform

import (
	"testing"

	"github.com/hashicorp/terraform/addrs"
	"github.com/hashicorp/terraform/states"
)

func TestNodeDataDestroyExecute(t *testing.T) {
	state := states.NewState()
	state.Module(addrs.RootModuleInstance).SetResourceInstanceCurrent(
		addrs.Resource{
			Mode: addrs.DataResourceMode,
			Type: "test_instance",
			Name: "foo",
		}.Instance(addrs.NoKey),
		&states.ResourceInstanceObjectSrc{
			Status:    states.ObjectReady,
			AttrsJSON: []byte(`{"dynamic":{"type":"string","value":"hello"}}`),
		},
		addrs.AbsProviderConfig{
			Provider: addrs.NewDefaultProvider("test"),
			Module:   addrs.RootModule,
		},
	)
	ctx := &MockEvalContext{
		StateState: state.SyncWrapper(),
	}

	node := NodeDestroyableDataResourceInstance{&NodeAbstractResourceInstance{
		Addr: addrs.Resource{
			Mode: addrs.DataResourceMode,
			Type: "test_instance",
			Name: "foo",
		}.Instance(addrs.NoKey).Absolute(addrs.RootModuleInstance),
	}}

	err := node.Execute(ctx, walkApply)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	// verify resource removed from state
	if state.HasResources() {
		t.Fatal("resources still in state after NodeDataDestroy.Execute")
	}
}
