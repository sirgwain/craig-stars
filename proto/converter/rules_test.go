package converter

import (
	"testing"

	"github.com/sirgwain/craig-stars/cs"
	craig_starsv1 "github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestBattleRulesRoundTrip(t *testing.T) {
	rules := cs.NewRulesWithSeed(1)
	rules.TorpedoSplashDamage = 0.25
	rules.MovementMassVariance = 0.2

	encoded, err := proto.Marshal(C.ConvertCSRules(&rules))
	require.NoError(t, err)
	var received craig_starsv1.Rules
	require.NoError(t, proto.Unmarshal(encoded, &received))

	decoded := C.ConvertRules(&received)
	assert.Equal(t, rules.BattleRules, decoded.BattleRules)

	// The WASM bindings use VT protobuf serialization instead of proto.Marshal.
	encoded, err = C.ConvertCSRules(&rules).MarshalVT()
	require.NoError(t, err)
	var wasmReceived craig_starsv1.Rules
	require.NoError(t, wasmReceived.UnmarshalVT(encoded))
	assert.Equal(t, rules.BattleRules, C.ConvertRules(&wasmReceived).BattleRules)
}
