package crawler

import (
	"crypto/ecdsa"
	"fmt"
	"net"
	"time"

	"github.com/Fantom-foundation/go-opera/gossip"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/rlpx"
	"github.com/pkg/errors"

	"github.com/ethereum/node-crawler/pkg/common"
)

var (
	lastStatusUpdate time.Time

	protocolCaps []p2p.Cap
)

func init() {
	// init protocolCaps:
	protocolCaps = make([]p2p.Cap, len(gossip.ProtocolVersions))
	for i, version := range gossip.ProtocolVersions {
		protocolCaps[i] = p2p.Cap{Name: gossip.ProtocolName, Version: version}
	}
	// TODO: snap protocol
}

func getClientInfo(opera *OperaStatus, nodeURL string, n *enode.Node) (*common.ClientInfo, error) {
	var info common.ClientInfo

	conn, sk, err := dial(n)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	if err = conn.SetDeadline(time.Now().Add(5 * time.Second)); err != nil {
		return nil, errors.Wrap(err, "cannot set conn deadline")
	}
	if err = writeHello(conn, sk); err != nil {
		return nil, err
	}
	if err = readHello(conn, &info); err != nil {
		return nil, err
	}
	// If node provides no eth version, we can skip it.
	if conn.negotiatedProtoVersion == 0 {
		return &info, nil
	}

	if err = conn.SetDeadline(time.Now().Add(15 * time.Second)); err != nil {
		return nil, errors.Wrap(err, "cannot set conn deadline")
	}
	if err = writeHandshake(conn, opera); err != nil {
		return nil, err
	}
	if err = readHandshake(conn, opera); err != nil {
		return nil, err
	}

	if err = conn.SetDeadline(time.Now().Add(15 * time.Second)); err != nil {
		return nil, errors.Wrap(err, "cannot set conn deadline")
	}
	if err = writeProgress(conn, opera); err != nil {
		return nil, err
	}
	if err = readProgress(conn, opera, &info); err != nil {
		return nil, err
	}

	// Disconnect from client
	_ = conn.Write(Disconnect{Reason: p2p.DiscQuitting})

	return &info, nil
}

// dial attempts to dial the given node and perform a handshake,
func dial(n *enode.Node) (*Conn, *ecdsa.PrivateKey, error) {
	var conn Conn

	// dial
	dialer := net.Dialer{Timeout: 10 * time.Second}
	fd, err := dialer.Dial("tcp", fmt.Sprintf("%v:%d", n.IP(), n.TCP()))
	if err != nil {
		return nil, nil, err
	}

	conn.Conn = rlpx.NewConn(fd, n.Pubkey())

	if err = conn.SetDeadline(time.Now().Add(15 * time.Second)); err != nil {
		return nil, nil, errors.Wrap(err, "cannot set conn deadline")
	}

	// do encHandshake
	ourKey, _ := crypto.GenerateKey()

	_, err = conn.Handshake(ourKey)
	if err != nil {
		return nil, nil, err
	}

	return &conn, ourKey, nil
}

func writeHello(conn *Conn, priv *ecdsa.PrivateKey) error {
	pub0 := crypto.FromECDSAPub(&priv.PublicKey)[1:]

	h := &Hello{
		Version: 5,
		Caps:    protocolCaps,
		ID:      pub0,
	}

	conn.ourHighestProtoVersion = gossip.ProtocolVersion
	conn.ourHighestSnapProtoVersion = 1 // TODO: snap

	return conn.Write(h)
}

func readHello(conn *Conn, info *common.ClientInfo) error {
	switch msg := conn.Read().(type) {
	case *Hello:
		// set snappy if version is at least 5
		if msg.Version >= 5 {
			conn.SetSnappy(true)
		}
		info.Capabilities = msg.Caps
		info.SoftwareVersion = msg.Version
		info.ClientType = msg.Name

		conn.negotiateEthProtocol(info.Capabilities)

		return nil
	case *Disconnect:
		return fmt.Errorf("bad hello handshake disconnect: %v", msg.Reason.Error())
	case *Error:
		return fmt.Errorf("bad hello handshake error: %v", msg.Error())
	default:
		return fmt.Errorf("bad hello handshake code: %v", msg.Code())
	}
}

func writeHandshake(conn *Conn, opera *OperaStatus) error {
	h := Handshake{
		ProtocolVersion: conn.negotiatedProtoVersion,
		NetworkID:       opera.NodeInfo.Network,
		Genesis:         opera.NodeInfo.Genesis,
	}
	return conn.Write(h)
}

func readHandshake(conn *Conn, opera *OperaStatus) error {
	switch msg := conn.Read().(type) {
	case *Handshake:
		if msg.Genesis != opera.NodeInfo.Genesis {
			return fmt.Errorf("Genesis mismatch: got %s, need %s", msg.Genesis.String(), opera.NodeInfo.Genesis.String())
		}
		if msg.NetworkID != opera.NodeInfo.Network {
			return fmt.Errorf("Network mismatch: got %d, need %d", msg.NetworkID, opera.NodeInfo.Network)
		}
		if msg.ProtocolVersion != conn.negotiatedProtoVersion {
			return fmt.Errorf("Protocol version mismatch: got %d, need %d", msg.ProtocolVersion, conn.negotiatedProtoVersion)
		}
	case *Disconnect:
		return fmt.Errorf("bad status handshake disconnect: %v", msg.Reason.Error())
	case *Error:
		return fmt.Errorf("bad status handshake error: %v", msg.Error())
	default:
		return fmt.Errorf("bad status handshake code: %v", msg.Code())
	}
	return nil
}

func writeProgress(conn *Conn, opera *OperaStatus) error {
	p := Progress(opera.Progress)
	return conn.Write(p)
}

func readProgress(conn *Conn, opera *OperaStatus, info *common.ClientInfo) error {
	switch msg := conn.Read().(type) {
	case *Progress:
	// TODO: update peer info here
	case *Disconnect:
		return fmt.Errorf("bad progress handshake disconnect: %v", msg.Reason.Error())
	case *Error:
		return fmt.Errorf("bad progress handshake error: %v", msg.Error())
	default:
		return fmt.Errorf("bad progress handshake code: %v", msg.Code())
	}
	return nil
}
