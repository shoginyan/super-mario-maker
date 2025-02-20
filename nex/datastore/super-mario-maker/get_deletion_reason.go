package nex_datastore_super_mario_maker

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	datastore_super_mario_maker "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/super-mario-maker"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
)

func GetDeletionReason(err error, packet nex.PacketInterface, callID uint32, dataIDLst []uint64) uint32 {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nex.ResultCodes.DataStore.Unknown
	}

	client := packet.Sender()

	// TODO - Complete this
	// * It's not actually known what the
	// * real "deletion reason" values are.
	// * This is stubbed until we figure
	// * that out
	pDeletionReasons := make([]uint32, 0)

	for range dataIDLst {
		// * Every course I've checked has had this
		// * set to 0, even if the course is not
		// * deleted
		pDeletionReasons = append(pDeletionReasons, 0)
	}

	rmcResponseStream := nex.NewStreamOut(globals.SecureServer)

	rmcResponseStream.WriteListUInt32LE(pDeletionReasons)

	rmcResponseBody := rmcResponseStream.Bytes()

	rmcResponse := nex.NewRMCResponse(datastore_super_mario_maker.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore_super_mario_maker.MethodGetDeletionReason, rmcResponseBody)

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV1(client, nil)

	responsePacket.SetVersion(1)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	globals.SecureServer.Send(responsePacket)

	return 0
}
