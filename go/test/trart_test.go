package test

import (
	"testing"

	jsoncdc "github.com/onflow/cadence/encoding/json"

	"github.com/onflow/cadence"
	"github.com/stretchr/testify/assert"
	"github.com/trart/test/flowutils"
	"github.com/trart/test/trartnft"
)

//test contract name
const trartContractName string = "Trart100LimitedEdittion"

//test nft data
var mintNFTs []trartnft.NFTData = []trartnft.NFTData{
	{ID: 1, Metadata: map[string]string{"CARD ID": "NFT1", "CARD NAME": "TRART NFT1"}},
	{ID: 2, Metadata: map[string]string{"CARD ID": "NFT2", "CARD NAME": "TRART NFT2"}},
	{ID: 3, Metadata: map[string]string{"CARD ID": "NFT3", "CARD NAME": "TRART NFT3"}},
}

//deploy our contract, and mint nfts later
func TestTrartNFTDeployContracts(t *testing.T) {
	b := flowutils.NewBlockchain()

	//deploy our contract
	nftAddress, trartNFTAddr, trartNFTSigner := trartnft.DeployContracts(t, b, trartContractName)

	t.Run("Should be able to mint NFTs", func(t *testing.T) {

		//mint NFTs
		trartnft.BatchMintNFT(t, b, nftAddress, trartNFTAddr, trartNFTSigner, trartContractName, mintNFTs)

		//query NFTs
		nfts := flowutils.ExecuteScriptAndCheck(
			t, b,
			trartnft.QueryAccountNFTScript(nftAddress.String(), trartNFTAddr.String(), trartNFTAddr.String(), trartContractName),
			nil,
		)
		userNfts, _ := nfts.ToGoValue().([]interface{})
		assert.EqualValues(t, len(mintNFTs), len(userNfts))

		//query a NFT's metadata
		metadata := flowutils.ExecuteScriptAndCheck(
			t, b,
			trartnft.QueryNFTScript(nftAddress.String(), trartNFTAddr.String(), trartContractName, int64(1)),
			nil,
		)
		nftMetadata, _ := metadata.ToGoValue().(map[interface{}]interface{})
		assert.EqualValues(t, "NFT1", nftMetadata["CARD ID"].(string))
	})
}

//transfer a NFT
func TestTransferNFT(t *testing.T) {
	b := flowutils.NewBlockchain()

	//deploy our contract
	nftAddress, trartNFTAddr, trartNFTSigner := trartnft.DeployContracts(t, b, trartContractName)

	//create receiver
	userAddress, userSigner, _ := flowutils.CreateAccount(t, b)

	// create a new Collection for receiver
	t.Run("Should be able to create a new empty NFT Collection", func(t *testing.T) {

		//setup receiver account
		trartnft.SetupAccount(t, b, userAddress, userSigner, nftAddress, trartNFTAddr, trartContractName)

		//test receiver's Collection
		result := flowutils.ExecuteScriptAndCheck(
			t, b,
			trartnft.IsInitalizedAccountScript(nftAddress.String(), trartNFTAddr.String(), trartContractName),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(userAddress))},
		)
		assert.EqualValues(t, cadence.NewUInt64(1), result)
	})

	// transfer an non-existing NFT
	t.Run("Should not be able to withdraw an NFT that does not exist in a collection", func(t *testing.T) {

		//non-existing NFT
		nonExistentID := uint64(3333333)

		//transder NFT
		trartnft.TransferNFT(
			t, b,
			nftAddress, trartNFTAddr, trartNFTSigner,
			nonExistentID, userAddress, true, trartContractName,
		)
	})

	// transfer an NFT
	t.Run("Should be able to withdraw an NFT and deposit to another accounts collection", func(t *testing.T) {

		//mint NFTs
		trartnft.BatchMintNFT(t, b, nftAddress, trartNFTAddr, trartNFTSigner, trartContractName, mintNFTs)

		//NFT id
		nftID := uint64(1)

		//transfer NFT
		trartnft.TransferNFT(
			t, b,
			nftAddress, trartNFTAddr, trartNFTSigner,
			nftID, userAddress, false, trartContractName,
		)
	})
}
