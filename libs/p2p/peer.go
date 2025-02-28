package p2p

import (
	"fmt"
	"net"
	"time"

	cmn "github.com/lianxiangcloud/linkchain/libs/common"
	"github.com/lianxiangcloud/linkchain/libs/crypto"
	"github.com/lianxiangcloud/linkchain/libs/log"

	"github.com/lianxiangcloud/linkchain/config"
	tmconn "github.com/lianxiangcloud/linkchain/libs/p2p/conn"
	"github.com/lianxiangcloud/linkchain/libs/ser"
)

var testIPSuffix uint32

// Peer is an interface representing a peer connected on a reactor.
//----------------------------------------------------------

// peerConn contains the raw connection and its config.
type peerConn struct {
	outbound     bool
	persistent   bool
	config       *config.P2PConfig
	conn         net.Conn    // source connection
	originalAddr *NetAddress // nil for inbound connections
}

// peer implements Peer.
//
// Before using a peer, you will need to perform a handshake on connection.
type peer struct {
	cmn.BaseService

	// raw peerConn and the multiplex connection
	peerConn
	mconn *tmconn.MConnection

	// peer's node info and the channel it knows about
	// channels = nodeInfo.Channels
	// cached to avoid copying nodeInfo in hasChannel
	nodeInfo NodeInfo
	channels []byte

	// User data
	Data *cmn.CMap
}

func newPeer(
	pc peerConn,
	mConfig tmconn.MConnConfig,
	nodeInfo NodeInfo,
	reactorsByCh map[byte]Reactor,
	chDescs []*tmconn.ChannelDescriptor,
	onPeerError func(Peer, interface{}),
) *peer {
	p := &peer{
		peerConn: pc,
		nodeInfo: nodeInfo,
		channels: nodeInfo.Channels,
		Data:     cmn.NewCMap(),
	}

	p.mconn = createMConnection(
		pc.conn,
		p,
		reactorsByCh,
		chDescs,
		onPeerError,
		mConfig,
	)
	p.BaseService = *cmn.NewBaseService(nil, "Peer", p)

	return p
}

func newOutboundPeerConn(
	addr *NetAddress,
	config *config.P2PConfig,
	persistent bool,
	ourNodePrivKey crypto.PrivKey,
) (peerConn, error) {
	conn, err := dial(addr, config)
	if err != nil {
		return peerConn{}, cmn.ErrorWrap(err, "Error creating peer")
	}

	pc, err := newPeerConn(conn, config, true, persistent, ourNodePrivKey, addr)
	if err != nil {
		if cerr := conn.Close(); cerr != nil {
			return peerConn{}, cmn.ErrorWrap(err, cerr.Error())
		}
		return peerConn{}, err
	}

	return pc, nil
}

func newInboundPeerConn(
	conn net.Conn,
	config *config.P2PConfig,
	ourNodePrivKey crypto.PrivKey,
) (peerConn, error) {

	// TODO: issue PoW challenge

	return newPeerConn(conn, config, false, false, ourNodePrivKey, nil)
}

func newPeerConn(
	rawConn net.Conn,
	cfg *config.P2PConfig,
	outbound, persistent bool,
	ourNodePrivKey crypto.PrivKey,
	originalAddr *NetAddress,
) (pc peerConn, err error) {
	conn := rawConn
	// Set deadline for secret handshake
	dl := time.Now().Add(cfg.HandshakeTimeout)
	if err := conn.SetDeadline(dl); err != nil {
		return pc, cmn.ErrorWrap(
			err,
			"Error setting deadline while encrypting connection",
		)
	}

	// Encrypt connection
	conn, err = tmconn.MakeSecretConnection(conn, ourNodePrivKey)
	if err != nil {
		return pc, cmn.ErrorWrap(err, "Error creating peer")
	}

	// Only the information we already have
	return peerConn{
		config:       cfg,
		outbound:     outbound,
		persistent:   persistent,
		conn:         conn,
		originalAddr: originalAddr,
	}, nil
}

//---------------------------------------------------
// Implements cmn.Service

// SetLogger implements BaseService.
func (p *peer) SetLogger(l log.Logger) {
	p.Logger = l
	p.mconn.SetLogger(l)
}

// OnStart implements BaseService.
func (p *peer) OnStart() error {
	if err := p.BaseService.OnStart(); err != nil {
		return err
	}
	err := p.mconn.Start()
	return err
}

// OnStop implements BaseService.
func (p *peer) OnStop() {
	p.BaseService.OnStop()
	p.mconn.Stop() // stop everything and close the conn
}

//---------------------------------------------------
// Implements Peer

// ID returns the peer's Uniquely identifies
func (p *peer) ID() string {
	return p.nodeInfo.ID()
}

// IsOutbound returns true if the connection is outbound, false otherwise.
func (p *peer) IsOutbound() bool {
	return p.peerConn.outbound
}

// IsPersistent returns true if the peer is persitent, false otherwise.
func (p *peer) IsPersistent() bool {
	return p.peerConn.persistent
}

// NodeInfo returns a copy of the peer's NodeInfo.
func (p *peer) NodeInfo() NodeInfo {
	return p.nodeInfo
}

// OriginalAddr returns the original address, which was used to connect with
// the peer. Returns nil for inbound peers.
func (p *peer) OriginalAddr() *NetAddress {
	if p.peerConn.outbound {
		return p.peerConn.originalAddr
	}
	return nil
}

//RemoteAddr returns the address of remote connetion
func (p *peer) RemoteAddr() net.Addr {
	return p.peerConn.conn.RemoteAddr()
}

// Status returns the peer's ConnectionStatus.
func (p *peer) Status() tmconn.ConnectionStatus {
	return p.mconn.Status()
}

// Send msg bytes to the channel identified by chID byte. Returns false if the
// send queue is full after timeout, specified by MConnection.
func (p *peer) Send(chID byte, msgBytes []byte) bool {
	if !p.IsRunning() {
		// see Switch#Broadcast, where we fetch the list of peers and loop over
		// them - while we're looping, one peer may be removed and stopped.
		return false
	} else if !p.hasChannel(chID) {
		return false
	}
	return p.mconn.Send(chID, msgBytes)
}

// TrySend msg bytes to the channel identified by chID byte. Immediately returns
// false if the send queue is full.
func (p *peer) TrySend(chID byte, msgBytes []byte) bool {
	if !p.IsRunning() {
		return false
	} else if !p.hasChannel(chID) {
		return false
	}
	return p.mconn.TrySend(chID, msgBytes)
}

// Get the data for a given key.
func (p *peer) Get(key string) interface{} {
	return p.Data.Get(key)
}

// Set sets the data for the given key.
func (p *peer) Set(key string, data interface{}) {
	p.Data.Set(key, data)
}

// hasChannel returns true if the peer reported
// knowing about the given chID.
func (p *peer) hasChannel(chID byte) bool {
	for _, ch := range p.channels {
		if ch == chID {
			return true
		}
	}
	// NOTE: probably will want to remove this
	// but could be helpful while the feature is new
	p.Logger.Debug(
		"Unknown channel for peer",
		"channel",
		chID,
		"channels",
		p.channels,
	)
	return false
}

func (p *peer) Close() error {
	err := p.Stop()
	return err
}

//---------------------------------------------------
// methods used by the Switch

// CloseConn should be called by the Switch if the peer was created but never
// started.
func (pc *peerConn) CloseConn() {
	pc.conn.Close() // nolint: errcheck
}

// HandshakeTimeout performs the P2P handshake between a given node
// and the peer by exchanging their NodeInfo. It sets the received nodeInfo on
// the peer.
// NOTE: blocking
func (pc *peerConn) HandshakeTimeout(
	ourNodeInfo NodeInfo,
	timeout time.Duration,
) (peerNodeInfo NodeInfo, err error) {
	// Set deadline for handshake so we don't block forever on conn.ReadFull
	if err := pc.conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		return peerNodeInfo, cmn.ErrorWrap(err, "Error setting deadline")
	}

	var trs, _ = cmn.Parallel(
		func(_ int) (val interface{}, err error, abort bool) {
			_, err = ser.EncodeWriterWithType(pc.conn, ourNodeInfo)
			return
		},
		func(_ int) (val interface{}, err error, abort bool) {
			_, err = ser.DecodeReaderWithType(
				pc.conn,
				&peerNodeInfo,
				int64(MaxNodeInfoSize()),
			)
			return
		},
	)
	if err := trs.FirstError(); err != nil {
		return peerNodeInfo, cmn.ErrorWrap(err, "Error during handshake")
	}

	// Remove deadline
	if err := pc.conn.SetDeadline(time.Time{}); err != nil {
		return peerNodeInfo, cmn.ErrorWrap(err, "Error removing deadline")
	}

	return peerNodeInfo, nil
}

// Addr returns peer's remote network address.
func (p *peer) Addr() net.Addr {
	return p.peerConn.conn.RemoteAddr()
}

// CanSend returns true if the send queue is not full, false otherwise.
func (p *peer) CanSend(chID byte) bool {
	if !p.IsRunning() {
		return false
	}
	return p.mconn.CanSend(chID)
}

// String representation.
func (p *peer) String() string {
	if p.outbound {
		return fmt.Sprintf("Peer{%v %v out}", p.mconn, p.ID())
	}

	return fmt.Sprintf("Peer{%v %v in}", p.mconn, p.ID())
}

// Type returns peer's type
func (p *peer) Type() string {
	return p.nodeInfo.Type.String()
}

//------------------------------------------------------------------
// helper funcs

func dial(addr *NetAddress, cfg *config.P2PConfig) (net.Conn, error) {
	conn, err := addr.DialTimeout(cfg.DialTimeout)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func createMConnection(
	conn net.Conn,
	p *peer,
	reactorsByCh map[byte]Reactor,
	chDescs []*tmconn.ChannelDescriptor,
	onPeerError func(Peer, interface{}),
	config tmconn.MConnConfig,
) *tmconn.MConnection {

	onReceive := func(chID byte, msgBytes []byte) {
		reactor := reactorsByCh[chID]
		if reactor == nil {
			// Note that its ok to panic here as it's caught in the conn._recover,
			// which does onPeerError.
			panic(cmn.Fmt("Unknown channel %X", chID))
		}
		reactor.Receive(chID, p, msgBytes)
	}

	onError := func(r interface{}) {
		onPeerError(p, r)
	}

	return tmconn.NewMConnectionWithConfig(
		conn,
		chDescs,
		onReceive,
		onError,
		config,
	)
}
