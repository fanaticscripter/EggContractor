package solo

import (
	"google.golang.org/protobuf/proto"

	"github.com/fanaticscripter/EggContractor/solo/pb"
)

func (c *SoloContract) Marshal() ([]byte, error) {
	return proto.Marshal(c.ToPBSoloContract())
}

func UnmarshalSoloContract(b []byte) (*SoloContract, error) {
	c := &pb.SoloContract{}
	err := proto.Unmarshal(b, c)
	if err != nil {
		return nil, err
	}
	return &SoloContract{c}, nil
}
