package converter

import (
	"testing"

	"github.com/sirgwain/craig-stars/cs"
	v1 "github.com/sirgwain/craig-stars/proto/gen/craig_stars/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func Test_message_newSpecsRoundTrip(t *testing.T) {
	for messageType := cs.PlayerMessageFleetColonizeInvalidNotPlanet; messageType <= cs.PlayerMessagePlanetRemoteTerraform; messageType++ {
		message := cs.PlayerMessage{Type: messageType, Spec: cs.PlayerMessageSpec{
			Amount: -1, Amount2: 0, HabType: cs.TerraformHabTypeGrav,
			Distance: 123.45, HasMassDriver: true, PlanetEmptied: true, Error: "details",
			RouteTarget: &cs.MapObjectTarget{TargetName: "Destination", TargetNum: 2},
		}}
		pb := C.ConvertCSPlayer(&cs.Player{Messages: []cs.PlayerMessage{message}}).Messages[0]
		data, err := proto.Marshal(pb)
		require.NoError(t, err)
		decoded := &v1.PlayerMessage{}
		require.NoError(t, proto.Unmarshal(data, decoded))
		require.Equal(t, message, C.ConvertPlayer(&v1.Player{Race: &v1.Race{}, Messages: []*v1.PlayerMessage{decoded}}).Messages[0])
		// The WASM protobuf implementation must preserve the same new fields and enum values.
		data, err = pb.MarshalVT()
		require.NoError(t, err)
		decoded = &v1.PlayerMessage{}
		require.NoError(t, decoded.UnmarshalVT(data))
		require.Equal(t, message, C.ConvertPlayer(&v1.Player{Race: &v1.Race{}, Messages: []*v1.PlayerMessage{decoded}}).Messages[0])
	}
}
