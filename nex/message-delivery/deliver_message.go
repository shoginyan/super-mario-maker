package nex_message_delivery

import (
	nex "github.com/PretendoNetwork/nex-go/v2"
	message_delivery "github.com/PretendoNetwork/nex-protocols-go/v2/message-delivery"
	"github.com/PretendoNetwork/super-mario-maker-secure/globals"
)

func DeliverMessage(err error, packet nex.PacketInterface, callID uint32, oUserMessage *nex.DataHolder) uint32 {
	if err != nil {
		globals.Logger.Error(err.Error())
		return nex.ResultCodes.DataStore.Unknown
	}

	client := packet.Sender()

	// TODO - See what this does

	rmcResponse := nex.NewRMCResponse(message_delivery.ProtocolID, callID)
	rmcResponse.SetSuccess(message_delivery.MethodDeliverMessage, nil)

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
