# 隐私链钱包客户端
<!-- TOC -->

- [隐私链钱包客户端](#隐私链钱包客户端)
    - [编译](#编译)
    - [依赖](#依赖)
    - [启动钱包客户端](#启动钱包客户端)
    - [停止钱包客户端](#停止钱包客户端)
    - [测试账户](#测试账户)
    - [RPC接口](#rpc接口)
        - [ltk_blockHeight](#ltk_blockheight)
        - [ltk_createSubAccount](#ltk_createsubaccount)
        - [ltk_balance](#ltk_balance)
        - [ltk_getAccountInfo](#ltk_getaccountinfo)
        - [ltk_selectAddress](#ltk_selectaddress)
        - [ltk_setRefreshBlockInterval](#ltk_setrefreshblockinterval)
        - [ltk_autoRefreshBlockchain](#ltk_autorefreshblockchain)
        - [ltk_rescanBlockchain](#ltk_rescanblockchain)
        - [ltk_status](#ltk_status)
        - [ltk_signUTXOTransaction](#ltk_signutxotransaction)
        - [ltk_sendUTXOTransaction](#ltk_sendutxotransaction)
        - [ltk_sendUTXOTransactionSplit](#ltk_sendutxotransactionsplit)
        - [ltk_getTxKey](#ltk_gettxkey)
        - [ltk_getMaxOutput](#ltk_getmaxoutput)
        - [personal_newAccount](#personal_newaccount)
        - [personal_getCue](#personal_getcue)
        - [personal_listAccounts](#personal_listaccounts)
        - [personal_unlockAccount](#personal_unlockaccount)
        - [personal_lockAccount](#personal_lockaccount)
        - [ltk_getProofKey](#ltk_getproofkey)
        - [ltk_checkProofKey](#ltk_checkproofkey)
        - [ltk_getBlockTransactionCountByNumber](#ltk_getblocktransactioncountbynumber)
        - [ltk_getBlockTransactionCountByHash](#ltk_getblocktransactioncountbyhash)
        - [ltk_getTransactionByBlockNumberAndIndex](#ltk_gettransactionbyblocknumberandindex)
        - [ltk_getTransactionByBlockHashAndIndex](#ltk_gettransactionbyblockhashandindex)
        - [ltk_getRawTransactionByBlockNumberAndIndex](#ltk_getrawtransactionbyblocknumberandindex)
        - [ltk_getRawTransactionByBlockHashAndIndex](#ltk_getrawtransactionbyblockhashandindex)
        - [ltk_getTransactionCount](#ltk_gettransactioncount)
        - [ltk_getTransactionByHash](#ltk_gettransactionbyhash)
        - [ltk_getRawTransactionByHash](#ltk_getrawtransactionbyhash)
        - [ltk_getTransactionReceipt](#ltk_gettransactionreceipt)

<!-- /TOC -->

## 编译

```shell
make
```

## 依赖

依赖 libxcrypto.a 库文件

## 启动钱包客户端

```shell
cd sbin && ./wallet.sh start
```

## 停止钱包客户端

```shell
./wallet.sh stop
```

## 测试账户

当钱包连接到测试网络对应的peer节点时（peer test模式运行），可以使用测试账户**0xa73810e519e1075010678d706533486d8ecc8000**,然后执行**personal_unlockAccount**命令进行解锁账户  

```shell
cp ../tests/UTC--2019-07-08T10-03-04.871669363Z--a73810e519e1075010678d706533486d8ecc8000 ./testdata/keystore/
```

## RPC接口

钱包默认RPC端口为 18082，也可自己修改启动参数指定

### ltk_blockHeight

功能：获取钱包当前同步的区块高度  
参数：无  
返回：  
    remote_height 字符串，十六进制，链节点高度  
    local_height 字符串，十六进制，钱包处理区块高度  
示例：  

```shell
curl -s -X POST http://127.0.0.1:18082 -d '{"jsonrpc":"2.0","id":"0","method":"ltk_blockHeight","params":[]}' -H 'Content-Type: application/json'|json_pp
{
   "jsonrpc" : "2.0",
   "id" : "0",
   "result" : {
      "local_height" : "0x1",
      "remote_height" : "0x1"
   }
}
```

### ltk_createSubAccount

功能：按指定最大子地址数，创建子账户  
参数：maxSubIdx 字符串，十六进制 最大子账户数  
返回：  
   bool 执行是否成功
示例：

```shell
curl -s -X POST http://127.0.0.1:18082 -d '{"jsonrpc":"2.0","id":"0","method":"ltk_createSubAccount","params":["0x2"]}' -H 'Content-Type: application/json'|json_pp
{
   "id" : "0",
   "jsonrpc" : "2.0",
   "result" : true
}
```

### ltk_balance

功能：查询余额  
参数：
    index 字符串，十六进制，子账户序号  
    token 字符串，十六进制，资产token标识，默认为"0x0000000000000000000000000000000000000000" 链克  
返回：  
    balance 字符串，十六进制，账户余额  
    address 账户地址，可以作为签交易的to地址  
    token 字符串，十六进制，资产token标识  
示例：

```shell
curl -s -X POST http://127.0.0.1:18082 -d '{"jsonrpc":"2.0","id":"0","method":"ltk_balance","params":[{"index":"0x0","token":"0x0000000000000000000000000000000000000000"}]}' -H 'Content-Type: application/json'|json_pp
{
   "id" : "0",
   "result" : {
      "balance" : "0x3fe0b6f90fc000",
      "token" : "0x0000000000000000000000000000000000000000",
      "address" : "BbzBXAnB8GwTsqVobB1Yg7SkbU1QUGFThPko8hkFP9VnQZ8WaVaLe4siM1r7tdKkrnFWXkHxZKPuj2gVjd6KZeoo1kyXwbu"
   },
   "jsonrpc" : "2.0"
}
```

### ltk_getAccountInfo

功能：查询所有子账户信息  
参数：  
   token 字符串，十六进制，资产token标识，默认为"0x0000000000000000000000000000000000000000" 链克  
返回：  
   total_balance 账户总余额
   token 字符串，十六进制，资产token标识，默认为"0x0000000000000000000000000000000000000000" 链克  
   eth_account 以太坊账户列表  
      address 账户地址  
      balance 账户余额  wei  
      nonce 账户nonce值  
  
   utxo_accounts账户信息列表  
      index 子账户序号  
      balance 账户余额  
      address 账户地址，可以作为签交易的to地址  
示例：

```shell
curl -s -X POST http://127.0.0.1:18082 -d '{"jsonrpc":"2.0","id":"0","method":"ltk_getAccountInfo","params":["0x0000000000000000000000000000000000000000"]}' -H 'Content-Type: application/json'|json_pp
{
   "result" : {
      "token" : "0x0000000000000000000000000000000000000000",
      "total_balance" : "0x1ed09bead87bfe65a1de98e51aa8b",
      "eth_account" : {
         "address" : "0xa73810e519e1075010678d706533486d8ecc8000",
         "balance" : "0x1ed09bead87bfd9c0bffaf7d85243",
         "nonce" : "0x75"
      },
      "utxo_accounts" : [
         {
            "balance" : "0x1f78f21ce1c600",
            "address" : "9x7envctz6N8oPwtstBddpgLoMvT2YmeU79z2A8ZMbf4hxvV2GFUrwPKmT6ko4YgTwMWEmNT1tFDg3DcTSNydftUHLnzj66",
            "index" : "0x0"
         },
         {
            "address" : "BbzBXAnB8GwTsqVobB1Yg7SkbU1QUGFThPko8hkFP9VnQZ8WaVaLe4siM1r7tdKkrnFWXkHxZKPuj2gVjd6KZeoo1kyXwbu",
            "index" : "0x1",
            "balance" : "0x1ff05b7c87e000"
         }
      ]
   }
   "id" : "0",
   "jsonrpc" : "2.0"
}
```

### ltk_selectAddress

功能：设置当前钱包默认账户地址  
参数：  
    address 十六进制字符串，账户地址  
返回：  
    true 设置成功，false 设置失败  
示例：

```shell
curl -s -X POST http://127.0.0.1:18082 -d '{"jsonrpc":"2.0","method":"ltk_selectAddress","params":["0xa73810e519e1075010678d706533486d8ecc8000"],"id":67}' -H 'Content-Type:application/json' |json_pp
{
   "id" : 67,
   "jsonrpc" : "2.0",
   "result" : true
}
```

### ltk_setRefreshBlockInterval

功能：设置钱包获取新区块高度的间隔
参数：  
    interval 正整数，十进制，区块刷新间隔，单位：秒
返回：  
    true 设置成功，false 设置失败  
示例：

```shell
curl -s -X POST http://127.0.0.1:18082 -d '{"jsonrpc":"2.0","method":"ltk_setRefreshBlockInterval","params":[15],"id":0}' -H 'Content-Type:application/json' |json_pp
{
   "id" : 0,
   "result" : true,
   "jsonrpc" : "2.0"
}
```

### ltk_autoRefreshBlockchain

功能：设置钱包自动刷新区块数据属性  
参数：  
    true 允许自动刷新，false 禁止自动刷新  
返回：  
    true 设置成功，false 设置失败  
示例：

```shell
curl -s -X POST http://127.0.0.1:18082 -d '{"jsonrpc":"2.0","method":"ltk_autoRefreshBlockchain","params":[true],"id":67}' -H 'Content-Type:application/json' |json_pp
{
   "result" : true,
   "id" : 67,
   "jsonrpc" : "2.0"
}
```

### ltk_rescanBlockchain

功能：重新扫描区块(对当前选中的账户生效)  
参数：  
    无  
返回：  
    true 设置成功，false 设置失败  
示例：

```shell
curl -s -X POST http://127.0.0.1:18082 -d '{"jsonrpc":"2.0","id":"0","method":"ltk_rescanBlockchain","params":[]}' -H 'Content-Type: application/json'|json_pp
{
   "id" : "0",
   "jsonrpc" : "2.0",
   "result" : true
}
```

### ltk_status

功能：获取钱包当前状态  
参数：  
    无  
返回：  
    local_height 字符串，十六进制，钱包同步的区块高度  
    remote_height 字符串，十六进制，连接的peer节点高度  
    eth_address 当前解锁的钱包地址  
    auto_refresh 是否自动刷新区块数据  
    wallet_open 钱包是否解锁了账户  
    wallet_version 钱包当前版本号  
    chain_version 链节点版本号  
    refresh_block_interval 账户刷新区块间隔秒数  
示例：  

```shell
curl -s -X POST http://127.0.0.1:18082 -d '{"jsonrpc":"2.0","id":"0","method":"ltk_status","params":[]}' -H 'Content-Type: application/json'|json_pp
{
   "result" : {
      "local_height" : "0x2f6",
      "refresh_block_interval" : 5,
      "wallet_version" : "0.1.0",
      "chain_version" : "0.1.0",
      "auto_refresh" : false,
      "eth_address" : "0xa73810e519e1075010678d706533486d8ecc8000",
      "remote_height" : "0x2f8",
      "wallet_open" : true
   },
   "jsonrpc" : "2.0",
   "id" : "0"
}
```

### ltk_signUTXOTransaction

功能：生成交易签名，可以使用 隐私链节点的 send_raw_transaction 接口发送交易  
参数：  
    from 字符串 以太坊from地址(使用以太坊账户转账必填,from字段不填，默认使用utxo账户签交易)  
    nonce 字符串，十六进制 以太坊账户nonce序号 (使用以太坊账户转账必填)  
    token 字符串，十六进制，资产token标识，默认为"0x0000000000000000000000000000000000000000" 链克  
    subaddrs 可以作为input的子账户序号数组(使用utxo账户转账,选填,默认可以使用所有子账户余额)  
    dests 转账目标用户，to地址与金额，可以是数组(如果是调用以太坊合约，可以带data字段)  
      dest结构示例：

```json
//转账目标为utxo
{
	"addr": "9x7envctz6N8oPwtstBddpgLoMvT2YmeU79z2A8ZMbf4hxvV2GFUrwPKmT6ko4YgTwMWEmNT1tFDg3DcTSNydftUHLnzj66",
	"amount": "0x56bc75e2d63100000", //十六进制字符串
	"remark": "0x9900000009" //附加交易留言，可选，最多32字节有效
},
//转账目标为以太坊账户或者合约
{
	"addr": "0xe50ab035b1cc691b84e415ff0931867f6a71b091",
	"amount": "0x56bc75e2d63100000", //十六进制字符串
	"data": "0x123344"//合约参数，可选
}

```

返回：txs 数组  
   raw 交易raw  
   hash 交易hash  
   gas 本次交易花费的gas值
   err_msg 交易错误信息
   err_code 交易错误码，0表示没有错误
示例：

```shell
#account 转 单utxo
curl -s -X POST http://127.0.0.1:18082 -d '{"jsonrpc":"2.0","id":"0","method":"ltk_signUTXOTransaction","params":[{"from":"0xa73810e519e1075010678d706533486d8ecc8000","nonce":"0x0","dests":[{"addr":"9x7envctz6N8oPwtstBddpgLoMvT2YmeU79z2A8ZMbf4hxvV2GFUrwPKmT6ko4YgTwMWEmNT1tFDg3DcTSNydftUHLnzj66","amount":"0x21e19e0c9bab2400000","remark":"0xa73810e519e1075010678d706533486d8ecc8000"}]}]}' -H 'Content-Type: application/json'|json_pp

{
   "jsonrpc" : "2.0",
   "id" : "0",
   "result" : {
      "txs" : [
         {
            "hash" : "0x4a45b985a6948b8ef46b6d0521d8cdd8bf28fbd9511f9204027a62002854bab9",
            "raw" : "0xf904ccf857cf4852c16c32...",
            "gas" : "0x1dcd6500"
         }
      ]
   }
}

#account 转 多utxo
curl -s -X POST http://127.0.0.1:18082 -d '{"jsonrpc":"2.0","id":"0","method":"ltk_signUTXOTransaction","params":[{"from":"0xa73810e519e1075010678d706533486d8ecc8000","nonce":"0x1","dests":[{"addr":"9x7envctz6N8oPwtstBddpgLoMvT2YmeU79z2A8ZMbf4hxvV2GFUrwPKmT6ko4YgTwMWEmNT1tFDg3DcTSNydftUHLnzj66","amount":"0x21e19e0c9bab2400000"},{"addr":"BgURGSDabvhCyrjungujZLVkf74FRxVaDccSXZY47TwQLpg9GPAdfiKZet5uuWyToZW4cJYKHuyXBZCjgB48qAkAKmW4B8g","amount":"0x21e19e0c9bab2400000","remark":"0x9900000009"}]}]}' -H 'Content-Type: application/json'|json_pp

#account 单account
curl -s -X POST http://127.0.0.1:18082 -d '{"jsonrpc":"2.0","id":"0","method":"ltk_signUTXOTransaction","params":[{"from":"0xa73810e519e1075010678d706533486d8ecc8000","nonce":"0x2","dests":[{"addr":"0xe50ab035b1cc691b84e415ff0931867f6a71b091","amount":"0x21e19e0c9bab2400000","data":"0x123344"}]}]}' -H 'Content-Type: application/json'|json_pp

#account 单account+utxo
curl -s -X POST http://127.0.0.1:18082 -d '{"jsonrpc":"2.0","id":"0","method":"ltk_signUTXOTransaction","params":[{"from":"0xa73810e519e1075010678d706533486d8ecc8000","nonce":"0x3","dests":[{"addr":"9x7envctz6N8oPwtstBddpgLoMvT2YmeU79z2A8ZMbf4hxvV2GFUrwPKmT6ko4YgTwMWEmNT1tFDg3DcTSNydftUHLnzj66","amount":"0x21e19e0c9bab2400000"},{"addr":"0xe50ab035b1cc691b84e415ff0931867f6a71b091","amount":"0x21e19e0c9bab2400000","data":"0x123344"}]}]}' -H 'Content-Type: application/json'|json_pp

#utxo账户转单utxo账户
curl -s -X POST http://127.0.0.1:18082 -d '{"jsonrpc":"2.0","id":"0","method":"ltk_signUTXOTransaction","params":[{"subaddrs":[0],"dests":[{"addr":"9x7envctz6N8oPwtstBddpgLoMvT2YmeU79z2A8ZMbf4hxvV2GFUrwPKmT6ko4YgTwMWEmNT1tFDg3DcTSNydftUHLnzj66","amount":"0x4563918244f40000"}]}]}' -H 'Content-Type: application/json'|json_pp
{
   "jsonrpc" : "2.0",
   "id" : "0",
   "result" : {
      "tx" : [
         {
            "hash" : "0x00b42cd0f53c32098e74b0170460130ebe291b4497c8dbf4bbd2f9741039e0a4",
            "raw" : "0xf9063deb10c698b1d373...",
            "gas" : "0x7a120",
            "err_msg" : "",
            "err_code" : 0
         }
      ]
   }
}

#utxo 转 utxo
curl -s -X POST http://127.0.0.1:18082 -d '{"jsonrpc":"2.0","id":"0","method":"ltk_signUTXOTransaction","params":[{"subaddrs":[0],"dests":[{"addr":"9x7envctz6N8oPwtstBddpgLoMvT2YmeU79z2A8ZMbf4hxvV2GFUrwPKmT6ko4YgTwMWEmNT1tFDg3DcTSNydftUHLnzj66","amount":"0x4563918244f40000"}]}]}' -H 'Content-Type: application/json'|json_pp

#utxo 转 单account
curl -s -X POST http://127.0.0.1:18082 -d '{"jsonrpc":"2.0","id":"0","method":"ltk_signUTXOTransaction","params":[{"subaddrs":[0],"dests":[{"addr":"0xe50ab035b1cc691b84e415ff0931867f6a71b091","amount":"0x4563918244f40000","data":"0x123344"}]}]}' -H 'Content-Type: application/json'|json_pp

#utxo 转 单account+utxo
curl -s -X POST http://127.0.0.1:18082 -d '{"jsonrpc":"2.0","id":"0","method":"ltk_signUTXOTransaction","params":[{"subaddrs":[0],"dests":[{"addr":"9x7envctz6N8oPwtstBddpgLoMvT2YmeU79z2A8ZMbf4hxvV2GFUrwPKmT6ko4YgTwMWEmNT1tFDg3DcTSNydftUHLnzj66","amount":"0x4563918244f40000"},{"addr":"0xe50ab035b1cc691b84e415ff0931867f6a71b091","amount":"0x4563918244f40000","data":"0x123344"}]}]}' -H 'Content-Type: application/json'|json_pp

#utxo 转 多utxo
curl -s -X POST http://127.0.0.1:18082 -d '{"jsonrpc":"2.0","id":"0","method":"ltk_signUTXOTransaction","params":[{"subaddrs":[0],"dests":[{"addr":"9x7envctz6N8oPwtstBddpgLoMvT2YmeU79z2A8ZMbf4hxvV2GFUrwPKmT6ko4YgTwMWEmNT1tFDg3DcTSNydftUHLnzj66","amount":"0x4563918244f40000"},{"addr":"BbzBXAnB8GwTsqVobB1Yg7SkbU1QUGFThPko8hkFP9VnQZ8WaVaLe4siM1r7tdKkrnFWXkHxZKPuj2gVjd6KZeoo1kyXwbu","amount":"0x4563918244f40000"}]}]}' -H 'Content-Type: application/json'|json_pp
```

### ltk_sendUTXOTransaction

功能：直接通过钱包发送交易，返还交易hash与交易数据（ltk_sendUTXOTransaction接口，只能发送一笔交易，如果一次转账被拆分成多个子交易，需要使用ltk_sendUTXOTransactionSplit接口）  
参数：  同 **ltk_signUTXOTransaction**  
返回：
   raw 交易raw数组  
   hash 对应的交易hash数组  
   err_code 错误码 0 没有错误  
   err_msg 错误信息  
示例：  

```shell
curl -s -X POST http://127.0.0.1:18082 -d '{"jsonrpc":"2.0","id":"0","method":"ltk_sendUTXOTransaction","params":[{"subaddrs":[0],"dests":[{"addr":"9xdUJ5cXrD3f4V7EzwgUrE5Uiku6VERRcP9xsMYpKjk5REmA8fL8AvpWaznNpL7WuwXGqdoPoDhSyW5oQjyiZNQEQsVYZLx","amount":"0x9900000009"}]}]}' -H 'Content-Type: application/json'|json_pp
{
   "id" : "0",
   "jsonrpc" : "2.0",
   "result" : {
      "tx" : [
         {
            "hash" : "0x821c89178db0dc87d96f6bd1cdd71552660b7b8d1b3d24e19894cf8209a14b44",
            "err_code" : 0,
            "err_msg" : "",
            "raw" : "f90c0c80f..."
         }
      ]
   }
}
```

### ltk_sendUTXOTransactionSplit

功能：直接通过钱包发送交易，返还交易hash与交易数据（如果一次转账被拆分成多个子交易，需要使用ltk_sendUTXOTransactionSplit接口）  
参数：  同 **eth_signUTXOTransaction**  
返回：
   raw 交易raw数组  
   hash 对应的交易hash数组  
   err_code 错误码 0 没有错误  
   err_msg 错误信息  
示例：  

```shell
curl -s -X POST http://127.0.0.1:18082 -d '{"jsonrpc":"2.0","id":"0","method":"ltk_sendUTXOTransactionSplit","params":[{"subaddrs":[0],"dests":[{"addr":"9x7envctz6N8oPwtstBddpgLoMvT2YmeU79z2A8ZMbf4hxvV2GFUrwPKmT6ko4YgTwMWEmNT1tFDg3DcTSNydftUHLnzj66","amount":"0x9900000009"}]}]}' -H 'Content-Type: application/json'|json_pp
{
   "jsonrpc" : "2.0",
   "result" : {
      "tx" : [
         {
            "err_code" : -1,
            "raw" : "f90c0c80f...",
            "err_msg" : "double spend",
            "hash" : "0x0000000000000000000000000000000000000000000000000000000000000000"
         }
      ]
   },
   "id" : "0"
}
```

### ltk_getTxKey

功能：获取交易私钥  
参数：  
    交易hash  
返回：  
    交易私钥 字符串  
示例：

```shell
curl -s -X POST http://127.0.0.1:18082 -d '{"jsonrpc":"2.0","method":"ltk_getTxKey","params":["0xbd7c414769329e6c7511ffcba8567349310da0db8ff4703ed35a8cd471fd2a68"],"id":67}' -H 'Content-Type:application/json' |json_pp
{
   "jsonrpc" : "2.0",
   "id" : 67,
   "result" : "0xe1793e7859af982be02ba0fe288ac2a9a881c0c179833015c68951fe3e00dc0c"
}
```

### ltk_getMaxOutput

功能：获取交易私钥  
参数：  
    tokenID string，0x0000000000000000000000000000000000000000 表示链克  
返回：  
    maxOutput 字符串,十六进制。该token对应的output个数。  
示例：

```shell
curl -s -X POST http://127.0.0.1:18082 -d '{"jsonrpc":"2.0","method":"ltk_getMaxOutput","params":["0x0000000000000000000000000000000000000000"],"id":67}' -H 'Content-Type:application/json' |json_pp
{
   "jsonrpc" : "2.0",
   "result" : "0x35",
   "id" : 67
}
```

### personal_newAccount

功能：创建一个钱包账户  
参数：  password 钱包密码  
返回：
   以太坊格式钱包地址，钱包文件保持在keystore目录下  
示例：  

```shell
curl -s -X POST http://127.0.0.1:18082 -d '{"jsonrpc":"2.0","method":"personal_newAccount","params":["1234","my name is lihua"],"id":67}' -H 'Content-Type:application/json' |json_pp
{
   "id" : 67,
   "result" : "0xe2c3f791a6fbc16252333fa5c89b5bcff9cf27ea",
   "jsonrpc" : "2.0"
}
```

### personal_getCue

功能：获取钱包密钥线索  
参数：  address 钱包账户地址  
返回：
   创建钱包账户时的密码提示信息  
示例：  

```shell
curl -s -X POST http://127.0.0.1:18082 -d '{"jsonrpc":"2.0","method":"personal_getCue","params":["0xe2c3f791a6fbc16252333fa5c89b5bcff9cf27ea"],"id":1}' -H 'Content-Type:application/json' |json_pp
{
   "id" : 1,
   "result" : "my name is lihua",
   "jsonrpc" : "2.0"
}
```

### personal_listAccounts

功能：查看钱包列表  
参数：  
返回：  
  保持在keystore目录下的以太坊格式钱包地址列表  
示例：

```shell
curl -s -X POST http://127.0.0.1:18082 -d '{"jsonrpc":"2.0","method":"personal_listAccounts","params":[],"id":67}' -H 'Content-Type:application/json' |json_pp
{
   "id" : 67,
   "jsonrpc" : "2.0",
   "result" : [
      "0xad2f67755eca0a24e7c587a2f494652bafe82c89",
      "0x67fa55ace3288da59d4e57c70ed97f04500d8000",
      "0x177a2ab1ffb9b6ba2b2de14abcc7e47d67d94c60"
   ]
}
```

### personal_unlockAccount

功能：解锁一个钱包账户  
参数：  
address 以太坊钱包地址  
password 钱包密码  
duration 解锁时间，单位秒  
返回：  
   true 解锁成功  
示例：

```shell
curl -s -X POST http://127.0.0.1:18082 -d '{"jsonrpc":"2.0","method":"personal_unlockAccount","params":["0xa73810e519e1075010678d706533486d8ecc8000","1234",3600],"id":67}' -H 'Content-Type:application/json' |json_pp
{
   "result" : true,
   "id" : 67,
   "jsonrpc" : "2.0"
}
```

### personal_lockAccount

功能：锁定一个钱包账户  
参数：  
address 以太坊钱包地址  
返回：  
   true 锁定成功  
示例：

```shell
curl -s -X POST http://127.0.0.1:18082 -d '{"jsonrpc":"2.0","id":"0","method":"personal_lockAccount","params":["0xa73810e519e1075010678d706533486d8ecc8000"]}' -H 'Content-Type: application/json'|json_pp
{
   "result" : true,
   "jsonrpc" : "2.0",
   "id" : "0"
}
```

### ltk_getProofKey

功能：获取交易验证私钥
参数：  
hash 交易hash  
addr utxo地址
返回：  
     proof_key 字符串
示例：

```shell
curl -s -X POST http://127.0.0.1:18082 -d '{"jsonrpc":"2.0","id":"0","method":"ltk_getProofKey","params":[{"addr":"9tCfNdQ4VKsFN8f2fKioqhHVQCjF2UREfUsxWZd8tmS96MBRt6qho4xRpvS2fGd8yQUZ9CTEQVeQcczyzSu53fHQKZLUARs","hash":"0x5ceeeab09fa5fabd52fc2f4c966b75c7ed2cf41e64f51057d8f22a553363ba12"}]}' -H 'Content-Type: application/json'|json_reformat
{
    "jsonrpc": "2.0",
    "id": "0",
    "result": {
        "proof_key": "086465f04dc90d5aaae3099e1f740431706c04cd6479d61f2d205c76a7177736"
    }
}
```

### ltk_checkProofKey

功能：验证交易是否包含转账
参数：  
hash 交易hash  
addr utxo地址
key  验证私钥
返回：  
     hash 交易hash
     addr 收款utxo地址
     amount 金额
示例：

```shell
curl -s -X POST http://127.0.0.1:18082 -d '{"jsonrpc":"2.0","id":"0","method":"ltk_checkProofKey","params":[{"addr":"9tCfNdQ4VKsFN8f2fKioqhHVQCjF2UREfUsxWZd8tmS96MBRt6qho4xRpvS2fGd8yQUZ9CTEQVeQcczyzSu53fHQKZLUARs","hash":"0x5ceeeab09fa5fabd52fc2f4c966b75c7ed2cf41e64f51057d8f22a553363ba12","key":"086465f04dc90d5aaae3099e1f740431706c04cd6479d61f2d205c76a7177736"}]}' -H 'Content-Type: application/json'|json_reformat
{
    "jsonrpc": "2.0",
    "id": "0",
    "result": {
        "records": [
            {
                "hash": "0x5ceeeab09fa5fabd52fc2f4c966b75c7ed2cf41e64f51057d8f22a553363ba12",
                "addr": "9tCfNdQ4VKsFN8f2fKioqhHVQCjF2UREfUsxWZd8tmS96MBRt6qho4xRpvS2fGd8yQUZ9CTEQVeQcczyzSu53fHQKZLUARs",
                "amount": "0x635c9adc5dea00000"
            }
        ]
    }
}
```

### ltk_getBlockTransactionCountByNumber

功能：查询区块中交易数
参数：  
blockNumber 字符串，十六进制 区块高度
返回：  
     txCount 字符串，十六进制 区块中包含的交易数
示例：

```shell
curl -s -X POST http://127.0.0.1:18082 -d '{"jsonrpc":"2.0","id":"0","method":"ltk_getBlockTransactionCountByNumber","params":["0x2"]}' -H 'Content-Type: application/json'|json_pp  
{
   "id" : "0",
   "jsonrpc" : "2.0",
   "result" : "0x1"
}
```

### ltk_getBlockTransactionCountByHash

功能：查询区块中交易数
参数：  
blockHash 字符串，十六进制 区块哈希
返回：  
     txCount 字符串，十六进制 区块中包含的交易数
示例：

```shell
curl -s -X POST http://127.0.0.1:18082 -d '{"jsonrpc":"2.0","id":"0","method":"ltk_getBlockTransactionCountByHash","params":["0x93de14f0a3ecc931f508bd155d67d79c59146c18caf48573cfa45dbc8b7556d8"]}' -H 'Content-Type: application/json'|json_pp  
{
   "id" : "0",
   "jsonrpc" : "2.0",
   "result" : "0x1"
}
```

### ltk_getTransactionByBlockNumberAndIndex

功能：查询区块中交易数
参数：  
blockHash 字符串，十六进制 区块哈希
返回：  
     txCount 字符串，十六进制 区块中包含的交易数
示例：

```shell
curl -s -X POST http://127.0.0.1:18082 -d '{"jsonrpc":"2.0","id":"0","method":"ltk_getTransactionByBlockNumberAndIndex","params":["0x1","0x0"]}' -H 'Content-Type: application/json'|json_pp  
{
   "jsonrpc" : "2.0",
   "id" : "0",
   "result" : {
      "txEntry" : {
         "txIndex" : "0",
         "blockHeight" : "1",
         "blockHash" : "0x93de14f0a3ecc931f508bd155d67d79c59146c18caf48573cfa45dbc8b7556d8"
      },
      "txType" : "utx",
      "txHash" : "0x2f9c05caefe4372118d951a58e5f6992c597d0bb063b63d3aa0399c0f1e520d6",
      "tx" : {
         "type" : "utx",
         "value" : {
            "extra" : null,
            "outputs" : [
               {
                  "value" : {
                     "otaddr" : "hRR7vvCs/ZZ2abg10OHQzFPMpKVFm44OFF9J9jugYXo=",
                     "amount" : 0,
                     "remark" : "0Eo0CGfbfd2z8CHVbsFIJqACxHA+mClvWVV2E/VCHxY="
                  },
                  "type" : "UTXOOutput"
               }
            ],
            "r_key" : "YX5tj5ADji5N+pN1GHEQi6wKMg+hJFJU84TLmSIXyxI=",
            "rct_signature" : {
               "P" : {
                  "RangeSigs" : null,
                  "MGs" : null,
                  "PseudoOuts" : null,
                  "Ss" : null,
                  "Bulletproofs" : [
                     {
                        "Aa" : "z0GgaCbTDtZLVIRxBVsKwov0bYai8NBMv7Iv52XjCwU=",
                        "T1" : "92cTIcIZfw/tBb/rX3t806gXV1t/2Ntfg3lK2P9M5ck=",
                        "A" : "NGwZhk7+TT2ie/ZiLNgV3OJ4sP9I1oEv5a2dNF3Mpzo=",
                        "Taux" : "v54uVTgmrWWSupA2YY8fNc/8U370OaGemUm7vLrVxgQ=",
                        "T2" : "nsHok3PJG2kLhdhjoaQFNelsqMQJuojQxL5VyoN9gts=",
                        "S" : "G1di0pUV40UuRFbeVW6MnEJjFT882XaQrJAaUIzKmtk=",
                        "B" : "IjDfHztddRZBRqkAFMbvEeVgWZFHJmk7LCrFBusqiQg=",
                        "Mu" : "qcA0IgMS1bm8Weto3bZ+DIDqQS/qLB12uSSVhVxlfgI=",
                        "R" : [
                           "udH5Sy8pL795/N+Pe+yqDNarxGqxx/Tb8ooxf04O+yo=",
                           "RlzAq9Nzr3UkC1PPMxq8ZKl85abmfv3Cj+4D8Nbpx34=",
                           "ql34CNfS4CjJ6sHJJtxLg261I9L+jIsc9qdbt8LqDQU=",
                           "PAAIQO8jnfmpube7YeFJ9oXL4VdCy+UxNAhSPpPiFvk=",
                           "geI7ysA0KFimi4lh/8Ug99DfK3je/COPTsylIUBlTp8=",
                           "dd3XLnsiMgbGd51VG5gSL8JHwhdSqkYLbbx81FKiZ2M=",
                           "ddqkCygAoq/Ph0VOpwTeGmeE74K+PwfGXNYXVjprhH4="
                        ],
                        "L" : [
                           "gWV5p3WvZFInut+BN0/fvEPg6S02vPGkTBG6Ok2RW3M=",
                           "urK2YqXys4QsQY3S0BdAY3/OmJecVqGBULD6DIphkNU=",
                           "4eMC1HJIaXAXfxLLtPRx9SZcJsuJVJ1O9BB8BZixPdQ=",
                           "o/efGtzKxm6HgbpfFCDnOd/elcsDSN2F9CJy37yBfns=",
                           "T+6+xybn02B+kTfma8JZfp52qylD6Sjsm57mOlbG3MM=",
                           "agFNcH5VJ7n2BhVuQBvZyFrE+dujLJAZXer5ZtBDZYI=",
                           "JUjP7a/Vp+nOrXNHEZR6D2a1BVeGOyhoUx+xeAckpFU="
                        ],
                        "T" : "JjPniUObCVkEKnUp5mW29YusoFx21nZLTh2nsNET8w4="
                     }
                  ]
               },
               "RctSigBase" : {
                  "OutPk" : [
                     {
                        "Mask" : "ypDcUwW0mzj7sDEQ0m9CkkhUlPT3HHQOdUbYrY7zeMQ=",
                        "Dest" : "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="
                     }
                  ],
                  "EcdhInfo" : [
                     {
                        "Mask" : "SoAUZyYa87/Y/4WB/H6HhE+G+FK6cyBTwiV637gaJQw=",
                        "SenderPK" : "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
                        "Amount" : "NAnJdh7dy4B8lpFO9GUV5UeFltknTyVXMwm9xTFt8Aw="
                     }
                  ],
                  "Type" : 0,
                  "TxnFee" : "0",
                  "PseudoOuts" : null
               }
            },
            "signature" : {
               "s" : "21060230996165888414594080395778246988172145535280638705203221867395606692834",
               "v" : 58342,
               "r" : "41175346618973896416810484304580252287215810708491173410791088609882574681352"
            },
            "inputs" : [
               {
                  "value" : {
                     "nonce" : "0",
                     "commit" : "XPR8dM5fkgpAZJ8JDIOgRgXaKRbz+OCV/tCkHC5wyjU=",
                     "amount" : 59000000000000000,
                     "cf" : "ZwnWu9mhhzr7q1tPbLcdc6+DNOJ72/bjaNADzMPXBQY="
                  },
                  "type" : "AccountInput"
               }
            ],
            "fee" : 50000000000000000,
            "add_keys" : [
               "Ml7VIYxDJC+DYvkx1SpCz8ustXzLmtgyX/3pakJc0Fg="
            ],
            "token_id" : "0x0000000000000000000000000000000000000000"
         }
      },
      "from" : "0xa73810e519e1075010678d706533486d8ecc8000"
   }
}
```

### ltk_getTransactionByBlockHashAndIndex

功能：用区块哈希与交易index查询交易
参数：  
blockHash 字符串，十六进制 区块哈希
index  字符串，十六进制 交易在区块中的序号
返回：  
     参考**ltk_getTransactionByBlockNumberAndIndex**
示例：

```shell
curl -s -X POST http://127.0.0.1:18082 -d '{"jsonrpc":"2.0","id":"0","method":"ltk_getTransactionByBlockHashAndIndex","params":["0x93de14f0a3ecc931f508bd155d67d79c59146c18caf48573cfa45dbc8b7556d8","0x0"]}' -H 'Content-Type: application/json'|json_pp
```

### ltk_getRawTransactionByBlockNumberAndIndex

功能：用区块哈希与交易index查询交易
参数：  
blockNumber 字符串，十六进制 区块编号
index  字符串，十六进制 交易在区块中的序号
返回：  
     raw 字符串，十六进制 交易raw
示例：

```shell
curl -s -X POST http://127.0.0.1:18082 -d '{"jsonrpc":"2.0","id":"0","method":"ltk_getRawTransactionByBlockNumberAndIndex","params":["0x1","0x0"]}' -H 'Content-Type: application/json'|json_pp
{
   "jsonrpc" : "2.0",
   "result" : "0xf90509f85...",
   "id" : "0"
}
```

### ltk_getRawTransactionByBlockHashAndIndex

功能：用区块哈希与交易index查询交易raw
参数：  
blockHash 字符串，十六进制 区块哈希
index  字符串，十六进制 交易在区块中的序号
返回：  
     raw 字符串，十六进制 交易raw
示例：

```shell
curl -s -X POST http://127.0.0.1:18082 -d '{"jsonrpc":"2.0","id":"0","method":"ltk_getRawTransactionByBlockHashAndIndex","params":["0x93de14f0a3ecc931f508bd155d67d79c59146c18caf48573cfa45dbc8b7556d8","0x0"]}' -H 'Content-Type: application/json'|json_pp
{
   "jsonrpc" : "2.0",
   "result" : "0xf90509f85...",
   "id" : "0"
}
```

### ltk_getTransactionCount

功能：查询账户交易nonce值
参数：  
address 字符串，十六进制 eth格式账户地址
blockNumber  字符串， 区块中的序号，目前只支持 "latest"  
返回：  
     nonce 字符串，十六进制 账户的交易nonce值
示例：

```shell
curl -s -X POST http://127.0.0.1:18082 -d '{"jsonrpc":"2.0","id":"0","method":"ltk_getTransactionCount","params":["0xa73810e519e1075010678d706533486d8ecc8000","latest"]}' -H 'Content-Type: application/json'|json_pp
{
   "jsonrpc" : "2.0",
   "id" : "0",
   "result" : "0x41c"
}
```

### ltk_getTransactionByHash

功能：根据交易哈希查询交易结构
参数：  
hash 字符串，十六进制 交易哈希
返回：  
     参考**ltk_getTransactionByBlockNumberAndIndex**
示例：

```shell
curl -s -X POST http://127.0.0.1:18082 -d '{"jsonrpc":"2.0","id":"0","method":"ltk_getTransactionByHash","params":["0x2f9c05caefe4372118d951a58e5f6992c597d0bb063b63d3aa0399c0f1e520d6"]}' -H 'Content-Type: application/json'|json_pp
```

### ltk_getRawTransactionByHash

功能：根据交易哈希查询交易raw
参数：  
hash 字符串，十六进制 交易哈希
返回：  
     raw 字符串，十六进制 交易raw
示例：

```shell
curl -s -X POST http://127.0.0.1:18082 -d '{"jsonrpc":"2.0","id":"0","method":"ltk_getRawTransactionByHash","params":["0x2f9c05caefe4372118d951a58e5f6992c597d0bb063b63d3aa0399c0f1e520d6"]}' -H 'Content-Type: application/json'|json_pp
{
   "jsonrpc" : "2.0",
   "id" : "0",
   "result" : "0xf90509f854cf48..."
}
```

### ltk_getTransactionReceipt

功能：根据交易哈希查询交易凭证
参数：  
hash 字符串，十六进制 交易哈希
返回：  
示例：

```shell
curl -s -X POST http://127.0.0.1:18082 -d '{"jsonrpc":"2.0","id":"0","method":"ltk_getTransactionReceipt","params":["0x2f9c05caefe4372118d951a58e5f6992c597d0bb063b63d3aa0399c0f1e520d6"]}' -H 'Content-Type: application/json'|json_pp
{
   "id" : "0",
   "result" : {
      "logsBloom" : "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
      "blockHash" : "0x93de14f0a3ecc931f508bd155d67d79c59146c18caf48573cfa45dbc8b7556d8",
      "gasUsed" : "0x7a120",
      "cumulativeGasUsed" : "0x7a120",
      "transactionHash" : "0x2f9c05caefe4372118d951a58e5f6992c597d0bb063b63d3aa0399c0f1e520d6",
      "logs" : [],
      "blockNumber" : "0x1",
      "contractAddress" : null,
      "transactionIndex" : "0x0",
      "from" : "0xa73810e519e1075010678d706533486d8ecc8000",
      "tokenAddress" : "0x0000000000000000000000000000000000000000",
      "status" : "0x1",
      "to" : []
   },
   "jsonrpc" : "2.0"
}
```
