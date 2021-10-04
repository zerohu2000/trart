package trartnft

import (
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/onflow/cadence"
	emulator "github.com/onflow/flow-emulator"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	sdktemplates "github.com/onflow/flow-go-sdk/templates"
	sdktest "github.com/onflow/flow-go-sdk/test"

	nftcontracts "github.com/onflow/flow-nft/lib/go/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	flowutils "github.com/trart/test/flowutils"
)

const (
	trartTransactionsRootPath = "../../transactions"
	trartScriptsRootPath      = "../../scripts"

	trartContractPath             = "../../contracts/TrartTemplateNFT.cdc"
	trartSetupAccountPath         = trartTransactionsRootPath + "/SetupTrartNFT.cdc"
	trartBatchMintNFTPath         = trartTransactionsRootPath + "/BatchMintNFT.cdc"
	trartTransferNFTPath          = trartTransactionsRootPath + "/TransferTrartNFT.cdc"
	trartIsInitializedAccountPath = trartScriptsRootPath + "/IsInitalizedAccount.cdc"
	trartQueryAccountNFTPath      = trartScriptsRootPath + "/QueryAccountNFT.cdc"
	trartQueryNFTPath             = trartScriptsRootPath + "/QueryNFT.cdc"
)

func DeployContracts(
	t *testing.T,
	b *emulator.Blockchain,
	contractName string,
) (flow.Address, flow.Address, crypto.Signer) {
	accountKeys := sdktest.AccountKeyGenerator()

	//
	nftCode := nftcontracts.NonFungibleToken()
	nftAddress, err := b.CreateAccount(
		nil,
		[]sdktemplates.Contract{
			{
				Name:   "NonFungibleToken",
				Source: string(nftCode),
			},
		},
	)
	require.NoError(t, err)

	_, err = b.CommitBlock()
	assert.NoError(t, err)

	//
	trartAccountKey, trartSigner := accountKeys.NewWithSigner()
	trartCode := placeholderTrartNFTCode(nftAddress.String(), contractName)
	trartAddr, err := b.CreateAccount(
		[]*flow.AccountKey{trartAccountKey},
		[]sdktemplates.Contract{
			{
				Name:   contractName,
				Source: string(trartCode),
			},
		},
	)
	assert.NoError(t, err)

	_, err = b.CommitBlock()
	assert.NoError(t, err)

	return nftAddress, trartAddr, trartSigner
}

func SetupAccount(
	t *testing.T,
	b *emulator.Blockchain,
	userAddress flow.Address,
	userSigner crypto.Signer,
	nftAddress flow.Address,
	trartAddress flow.Address,
	contractName string,
) {
	tx := flow.NewTransaction().
		SetScript(SetupAccountScript(nftAddress.String(), trartAddress.String(), contractName)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(userAddress)

	flowutils.SignAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, userAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), userSigner},
		false,
	)
}

type NFTData struct {
	ID       uint
	Metadata map[string]string
}

func BatchMintNFT(
	t *testing.T, b *emulator.Blockchain,
	nftAddress, trartAddr flow.Address,
	trartSigner crypto.Signer, contractName string, mintNFTs []NFTData,
) {
	tx := flow.NewTransaction().
		SetScript(BatchMintTrartScript(nftAddress.String(), trartAddr.String(), contractName)).
		SetGasLimit(9999).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(trartAddr)

	//1
	_ = tx.AddArgument(cadence.NewAddress(trartAddr))

	//2
	cardIDs := make([]cadence.Value, len(mintNFTs))
	for i, mintNFT := range mintNFTs {
		cardIDs[i] = cadence.NewUInt64(uint64(mintNFT.ID))
	}
	_ = tx.AddArgument(cadence.NewArray(cardIDs))

	//3
	metadatas := make([]cadence.Value, len(mintNFTs))
	for i, mintNFT := range mintNFTs {
		//
		datas := make([]cadence.KeyValuePair, 0)
		for k, v := range mintNFT.Metadata {
			key, err := cadence.NewValue(k)
			if err != nil {
				log.Println(err)
				continue
			}

			value, err := cadence.NewValue(v)
			if err != nil {
				log.Println(err)
				continue
			}

			datas = append(datas, cadence.KeyValuePair{
				Key:   key,
				Value: value,
			})
		}
		metadatas[i] = cadence.NewDictionary(datas)
	}
	_ = tx.AddArgument(cadence.NewArray(metadatas))

	flowutils.SignAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, trartAddr},
		[]crypto.Signer{b.ServiceKey().Signer(), trartSigner},
		false,
	)
}

func TransferNFT(
	t *testing.T, b *emulator.Blockchain,
	nftAddress, trartAddr flow.Address, trartSigner crypto.Signer,
	typeID uint64, recipientAddr flow.Address, shouldFail bool, contractName string,
) {

	tx := flow.NewTransaction().
		SetScript(TransferNFTScript(nftAddress.String(), trartAddr.String(), contractName)).
		SetGasLimit(1000).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(trartAddr)

	_ = tx.AddArgument(cadence.NewAddress(recipientAddr))
	_ = tx.AddArgument(cadence.NewUInt64(typeID))

	flowutils.SignAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, trartAddr},
		[]crypto.Signer{b.ServiceKey().Signer(), trartSigner},
		shouldFail,
	)
}

func placeholderTrartNFTCode(nftAddress, contractName string) []byte {
	code := string(flowutils.ReadFile(trartContractPath))
	code = strings.ReplaceAll(code, "0xNFTADDRESS", "0x"+nftAddress)
	code = strings.ReplaceAll(code, "TrartContractNFT", contractName)

	//fmt.Println(code)

	return []byte(code)
}

func SetupAccountScript(nftAddress, trartNftAddress, contractName string) []byte {
	code := string(flowutils.ReadFile(trartSetupAccountPath))
	code = strings.ReplaceAll(code, "0xNFTADDRESS", "0x"+nftAddress)
	code = strings.ReplaceAll(code, "0xTRARTNFTADDRESS", "0x"+trartNftAddress)
	code = strings.ReplaceAll(code, "TrartContractNFT", contractName)

	//fmt.Println(code)

	return []byte(code)
}

func BatchMintTrartScript(nftAddress, trartNftAddress, contractName string) []byte {

	code := string(flowutils.ReadFile(trartBatchMintNFTPath))
	code = strings.ReplaceAll(code, "0xNFTADDRESS", "0x"+nftAddress)
	code = strings.ReplaceAll(code, "0xTRARTNFTADDRESS", "0x"+trartNftAddress)
	code = strings.ReplaceAll(code, "TrartContractNFT", contractName)

	//fmt.Println(code)

	return []byte(code)
}

func TransferNFTScript(nftAddress, trartNftAddress, contractName string) []byte {
	code := string(flowutils.ReadFile(trartTransferNFTPath))
	code = strings.ReplaceAll(code, "0xNFTADDRESS", "0x"+nftAddress)
	code = strings.ReplaceAll(code, "0xTRARTNFTADDRESS", "0x"+trartNftAddress)
	code = strings.ReplaceAll(code, "TrartContractNFT", contractName)

	//fmt.Println(code)

	return []byte(code)
}

func IsInitalizedAccountScript(nftAddress, trartNftAddress, contractName string) []byte {
	code := string(flowutils.ReadFile(trartIsInitializedAccountPath))
	code = strings.ReplaceAll(code, "0xNFTADDRESS", "0x"+nftAddress)
	code = strings.ReplaceAll(code, "0xTRARTNFTADDRESS", "0x"+trartNftAddress)
	code = strings.ReplaceAll(code, "TrartContractNFT", contractName)

	//fmt.Println(code)

	return []byte(code)
}

func QueryAccountNFTScript(nftAddress, trartNftAddress, userAddress, contractName string) []byte {
	code := string(flowutils.ReadFile(trartQueryAccountNFTPath))
	code = strings.ReplaceAll(code, "0xNFTADDRESS", "0x"+nftAddress)
	code = strings.ReplaceAll(code, "0xTRARTNFTADDRESS", "0x"+trartNftAddress)
	code = strings.ReplaceAll(code, "TrartContractNFT", contractName)
	code = strings.ReplaceAll(code, "0xUSERADDRESS", "0x"+userAddress)

	//fmt.Println(code)

	return []byte(code)
}

func QueryNFTScript(nftAddress, trartNftAddress, contractName string, id int64) []byte {
	code := string(flowutils.ReadFile(trartQueryNFTPath))
	code = strings.ReplaceAll(code, "0xNFTADDRESS", "0x"+nftAddress)
	code = strings.ReplaceAll(code, "0xTRARTNFTADDRESS", "0x"+trartNftAddress)
	code = strings.ReplaceAll(code, "TrartContractNFT", contractName)
	code = strings.ReplaceAll(code, "%nftID", strconv.FormatInt(id, 10))

	//fmt.Println(code)

	return []byte(code)
}