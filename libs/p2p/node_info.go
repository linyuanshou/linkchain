package p2p

import (
	"encoding/hex"
	"fmt"
	"strings"

	cmn "github.com/lianxiangcloud/linkchain/libs/common"
	"github.com/lianxiangcloud/linkchain/libs/crypto"
	"github.com/lianxiangcloud/linkchain/types"
)

const (
	maxNodeInfoSize = 10240 // 10Kb
	maxNumChannels  = 16    // plenty of room for upgrades, for now
)

// MaxNodeInfoSize return Max size of the NodeInfo struct
func MaxNodeInfoSize() int {
	return maxNodeInfoSize
}

// NodeInfo is the basic node information exchanged
// between two peers during the P2P handshake.
type NodeInfo struct {
	PubKey     crypto.PubKeyEd25519 `json:"pub_key"`
	ListenAddr string               `json:"listen_addr"` // listen tcp addr,accepting tcp protocol incoming

	// Check compatibility.
	// Channels are HexBytes so easier to read as JSON
	Network  string       `json:"network"`  // network/chain ID
	Version  string       `json:"version"`  // major.minor.revision
	Channels cmn.HexBytes `json:"channels"` // channels this node knows about

	// ASCIIText fields
	Moniker string   `json:"moniker"` // arbitrary moniker
	Other   []string `json:"other"`   // other application specific data

	Type       types.NodeType `json:"type"`  // node type
	LocalAddrs []string       `json:"addrs"` // all address of node
}

// Validate checks the self-reported NodeInfo is safe.
// It returns an error if there
// are too many Channels, if there are any duplicate Channels,
// if the ListenAddr is malformed, or if the ListenAddr is a host name
// that can not be resolved to some IP.
// TODO: constraints for Moniker/Other? Or is that for the UI ?
// JAE: It needs to be done on the client, but to prevent ambiguous
// unicode characters, maybe it's worth sanitizing it here.
// In the future we might want to validate these, once we have a
// name-resolution system up.
// International clients could then use punycode (or we could use
// url-encoding), and we just need to be careful with how we handle that in our
// clients. (e.g. off by default).
func (info NodeInfo) Validate() error {
	if len(info.Channels) > maxNumChannels {
		return fmt.Errorf("info.Channels is too long (%v). Max is %v", len(info.Channels), maxNumChannels)
	}

	// Sanitize ASCII text fields.
	if !cmn.IsASCIIText(info.Moniker) || cmn.ASCIITrim(info.Moniker) == "" {
		return fmt.Errorf("info.Moniker must be valid non-empty ASCII text without tabs, but got %v", info.Moniker)
	}
	for i, s := range info.Other {
		if !cmn.IsASCIIText(s) || cmn.ASCIITrim(s) == "" {
			return fmt.Errorf("info.Other[%v] must be valid non-empty ASCII text without tabs, but got %v", i, s)
		}
	}

	channels := make(map[byte]struct{})
	for _, ch := range info.Channels {
		_, ok := channels[ch]
		if ok {
			return fmt.Errorf("info.Channels contains duplicate channel id %v", ch)
		}
		channels[ch] = struct{}{}
	}

	// ensure ListenAddr is good
	if info.ListenAddr == "" {
		return nil
	}

	_, err := NewNetAddressString(info.ListenAddr)
	return err
}

// CompatibleWith checks if two NodeInfo are compatible with eachother.
// CONTRACT: two nodes are compatible if the major version matches and network match
// and they have at least one channel in common.
func (info NodeInfo) CompatibleWith(other NodeInfo) error {
	iMajor, _, _, iErr := splitVersion(info.Version)
	oMajor, _, _, oErr := splitVersion(other.Version)

	// if our own version number is not formatted right, we messed up
	if iErr != nil {
		return iErr
	}

	// version number must be formatted correctly ("x.x.x")
	if oErr != nil {
		return oErr
	}

	// major version must match
	if iMajor != oMajor {
		return fmt.Errorf("Peer is on a different major version. Got %v, expected %v", oMajor, iMajor)
	}

	// nodes must be on the same network
	if info.Network != other.Network {
		return fmt.Errorf("Peer is on a different network. Got %v, expected %v", other.Network, info.Network)
	}

	// if we have no channels, we're just testing
	if len(info.Channels) == 0 {
		return nil
	}

	// for each of our channels, check if they have it
	found := false
OUTER_LOOP:
	for _, ch1 := range info.Channels {
		for _, ch2 := range other.Channels {
			if ch1 == ch2 {
				found = true
				break OUTER_LOOP // only need one
			}
		}
	}
	if !found {
		return fmt.Errorf("Peer has no common channels. Our channels: %v ; Peer channels: %v", info.Channels, other.Channels)
	}
	return nil
}

// NetAddress returns a NetAddress derived from the NodeInfo - ListenAddr
// Note that the ListenAddr is not authenticated and may not match that
// address actually dialed if its an outbound peer.
func (info NodeInfo) NetAddress() *NetAddress {
	netAddr, err := NewNetAddressString(info.ListenAddr)
	if err != nil {
		switch err.(type) {
		case ErrNetAddressLookup:
			// XXX If the peer provided a host name  and the lookup fails here
			// we're out of luck.
			// TODO: use a NetAddress in NodeInfo
		default:
			//panic(err) // everything should be well formed by now
		}
	}
	return netAddr
}

func (info NodeInfo) String() string {
	return fmt.Sprintf("NodeInfo{pk: %v, id: %v, moniker: %v, type: %v, network: %v [listen %v], addrs: %v, version: %v (%v)}",
		info.PubKey.String(), info.ID(), info.Moniker, info.Type, info.Network, info.ListenAddr, info.LocalAddrs, info.Version, info.Other)
}

// ID returns the peer's Uniquely identifies
func (info NodeInfo) ID() string {
	//return hex.EncodeToString(info.PubKey.Address())
	return hex.EncodeToString(crypto.Keccak256Hash(info.PubKey.Bytes()).Bytes())
}

func splitVersion(version string) (string, string, string, error) {
	spl := strings.Split(version, ".")
	if len(spl) != 3 {
		return "", "", "", fmt.Errorf("Invalid version format %v", version)
	}
	return spl[0], spl[1], spl[2], nil
}
