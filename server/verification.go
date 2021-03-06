package server

import (
	"errors"
)

// verifyMessage makes sure all required fields exist.
func (pi *packetInfo) verifyMessage() error {
	player := pi.tracker.playerByConnection(pi.conn)
	if player == nil {
		return errors.New("player disconnected")
	}

	for _, pkt := range pi.message.Packet {
		packet := pkt.GetPacket()

		if string(packet.GetSession()) != string(player.sessionID) {
			return errors.New("invalid session")
		}

		if packet.GetFromKey().String() != player.verificationKey {
			return errors.New("invalid verification key")
		}

		if pi.tracker.bannedByPool(player, true) {
			return errors.New("banned player")
		}

		if packet.GetNumber() != player.number {
			return errors.New("invalid user number")
		}

		to := packet.GetToKey()
		if to != nil {
			if pi.tracker.playerByVerificationKey(to.String()) == nil {
				return errors.New("invalid destination")
			}
		}
	}

	return nil
}
