package bootcli

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	"crypto/md5"

	"strings"

	"io/ioutil"

	"github.com/lianxiangcloud/linkchain/libs/crypto"
	"github.com/lianxiangcloud/linkchain/libs/hexutil"
	"github.com/lianxiangcloud/linkchain/libs/log"
	"github.com/lianxiangcloud/linkchain/libs/p2p/common"
	"github.com/lianxiangcloud/linkchain/types"
)

var (
	LocalNodeType  types.NodeType
	nodeTypeLocker sync.RWMutex //Just for unit test
)

const (
	Succ = 0
)

const (
	UDP           = "udp"
	TCP           = "tcp"
	RouteGetSeeds = "api/bootnode"
)

type Rnode struct {
	PubKey   string   `json:"pubkey"`
	Endpoint Endpoint `json:"endpoint"`
}

type GeetSeedsReq struct {
	Time   int64  `json:"time"`
	Sign   string `json:"sign"`
	Pubkey string `json:"pubkey"`
}

type GeetSeedsResp struct {
	Code    int     `json:"code"` //0:success，other:failed
	Message string  `json:"message"`
	Type    int     `json:"type"` //The identity type of this node, reference NodeType
	Seeds   []Rnode `json:"nodes"`
}

func buildGetSeedsURL(url string) string {
	return fmt.Sprintf("%s/%s", url, RouteGetSeeds)
}

func GetSeeds(bootSouce string, priv crypto.PrivKey, logger log.Logger) (nodes []*common.Node, localNodeType types.NodeType, err error) {
	if strings.HasPrefix(bootSouce, "http") || strings.HasPrefix(bootSouce, "https") {
		nodes, localNodeType, err = GetSeedsFromBootSvr(bootSouce, priv, logger)
	} else {
		nodes, localNodeType, err = getSeedsFromFile(bootSouce, logger)
	}
	nodeTypeLocker.Lock()
	LocalNodeType = localNodeType
	nodeTypeLocker.Unlock()
	return
}

func getSeedsFromFile(bootSouce string, logger log.Logger) (nodes []*common.Node, localNodeType types.NodeType, err error) {
	if len(bootSouce) == 0 {
		err = fmt.Errorf("getSeedsFromFile len(bootSouce) == 0")
		return
	}

	data, err := ioutil.ReadFile(bootSouce)
	if err != nil {
		return nil, 0, fmt.Errorf("getSeedsFromFile bootSouce:%v err:%v", bootSouce, err)
	}
	if nodes, localNodeType, err = parseAccounts(data, logger); err != nil {
		return
	}
	return
}

func parseAccounts(data []byte, logger log.Logger) (nodes []*common.Node, localNodeType types.NodeType, err error) {
	var resp GeetSeedsResp
	if err = json.Unmarshal(data, &resp); err != nil {
		return
	}
	nodes = RapNodes(resp.Seeds, logger)
	return nodes, types.NodeType(resp.Type), nil
}

func GetSeedsFromBootSvr(bootSvr string, priv crypto.PrivKey, logger log.Logger) (nodes []*common.Node, localNodeType types.NodeType, err error) {
	timeNowSecond := time.Now().Unix()
	timeString := fmt.Sprintf("time=%d", timeNowSecond)
	hash := md5.Sum([]byte(timeString))
	sign, err := priv.Sign(hash[:]) //crypto.Sign(hash[:], priv)
	if err != nil {
		logger.Error("GetSeedsFromBootSvr", "Sign err", err)
		return
	}

	postContent := GeetSeedsReq{
		Time:   timeNowSecond,
		Sign:   hex.EncodeToString(sign.Bytes()),
		Pubkey: hexutil.Encode(priv.PubKey().Bytes()),
	}

	var respBytes []byte
	var retry int
	for {
		respBytes, err = HttpPost(buildGetSeedsURL(bootSvr), postContent)
		if err != nil {
			logger.Error("GetSeedsFromBootSvr", "retry", retry, "HttpPost err", err)

			retry++
			if retry > 3 {
				return
			}

			time.Sleep(time.Second * time.Duration(3*retry))
			continue
		}
		break
	}

	var resp GeetSeedsResp
	err = json.Unmarshal(respBytes, &resp)
	if err != nil {
		logger.Error("GetSeedsFromBootSvr", "Unmarshal err", err)
		return
	}
	if resp.Code != Succ {
		err = fmt.Errorf("GetSeedsFromBootSvr getValidatorsAndSeeds code:%v != success,Retmsg:%v", resp.Code, resp.Message)
		return
	}
	nodes = RapNodes(resp.Seeds, logger)
	localNodeType = types.NodeType(resp.Type)
	return
}

func RapNodes(seeds []Rnode, logger log.Logger) (nodes []*common.Node) {
	logger.Debug("RapNodes", "len(seeds)", len(seeds))
	for i := 0; i < len(seeds); i++ {
		pubkey, err := hexutil.Decode(seeds[i].PubKey)
		if err != nil {
			logger.Error("GetSeedsFromBootSvr", "Decode err", err)
			continue
		}
		nodeID := common.NodeID(crypto.Keccak256Hash(pubkey))

		if len(seeds[i].Endpoint.IP) == 0 {
			tmpNode := &common.Node{ID: nodeID}
			nodes = append(nodes, tmpNode)
		} else {
			for j := 0; j < len(seeds[i].Endpoint.IP); j++ {
				tmpip := net.ParseIP(seeds[i].Endpoint.IP[j])
				tmpNode := &common.Node{IP: tmpip, ID: nodeID}
				for k, v := range seeds[i].Endpoint.Port {
					switch k {
					case UDP:
						tmpNode.UDP_Port = uint16(v)
					case TCP:
						tmpNode.TCP_Port = uint16(v)
					}
				}
				nodes = append(nodes, tmpNode)
			}
		}
	}
	return
}

func GetLocalNodeType() types.NodeType {
	nodeTypeLocker.RLock()
	defer nodeTypeLocker.RUnlock()
	return LocalNodeType
}
