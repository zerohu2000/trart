package trartnft

import (
	"log"
	"regexp"
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

var (
	TEMPLATE_NFTADDRESS        = regexp.MustCompile(`"[^"\s].*/NonFungibleToken.cdc"`)
	TEMPLATE_TRARTNFTADDRESS   = regexp.MustCompile(`"[^"\s].*/TrartTemplateNFT.cdc"`)
	TEMPLATE_TRARTCONTRACTNAME = regexp.MustCompile(`TrartContractNFT`)
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

func replaceAddressPlaceholders(code, nftAddress, trartNftAddress, trartContractName string) []byte {
	return []byte(flowutils.ReplaceImports(
		code,
		map[string]*regexp.Regexp{
			nftAddress:        TEMPLATE_NFTADDRESS,
			trartNftAddress:   TEMPLATE_TRARTNFTADDRESS,
			trartContractName: TEMPLATE_TRARTCONTRACTNAME,
		},
	))
}

func placeholderTrartNFTCode(nftAddress, contractName string) []byte {
	code := []byte(flowutils.ReplaceImports(
		string(flowutils.ReadFile(trartContractPath)),
		map[string]*regexp.Regexp{
			"0x" + nftAddress: TEMPLATE_NFTADDRESS,
			contractName:      TEMPLATE_TRARTCONTRACTNAME,
		},
	))

	//fmt.Println(string(code))

	return code
}

func SetupAccountScript(nftAddress, trartNftAddress, contractName string) []byte {
	code := replaceAddressPlaceholders(
		string(flowutils.ReadFile(trartSetupAccountPath)),
		"0x"+nftAddress,
		"0x"+trartNftAddress,
		contractName,
	)

	//fmt.Println(string(code))

	return []byte(code)
}

func BatchMintTrartScript(nftAddress, trartNftAddress, contractName string) []byte {
	code := replaceAddressPlaceholders(
		string(flowutils.ReadFile(trartBatchMintNFTPath)),
		"0x"+nftAddress,
		"0x"+trartNftAddress,
		contractName,
	)

	//fmt.Println(string(code))

	return []byte(code)
}

func TransferNFTScript(nftAddress, trartNftAddress, contractName string) []byte {
	code := replaceAddressPlaceholders(
		string(flowutils.ReadFile(trartTransferNFTPath)),
		"0x"+nftAddress,
		"0x"+trartNftAddress,
		contractName,
	)

	//fmt.Println(string(code))

	return []byte(code)
}

func IsInitalizedAccountScript(nftAddress, trartNftAddress, contractName string) []byte {
	code := replaceAddressPlaceholders(
		string(flowutils.ReadFile(trartIsInitializedAccountPath)),
		"0x"+nftAddress,
		"0x"+trartNftAddress,
		contractName,
	)

	//fmt.Println(string(code))

	return []byte(code)
}

func QueryAccountNFTScript(nftAddress, trartNftAddress, contractName string) []byte {
	code := replaceAddressPlaceholders(
		string(flowutils.ReadFile(trartQueryAccountNFTPath)),
		"0x"+nftAddress,
		"0x"+trartNftAddress,
		contractName,
	)

	//fmt.Println(string(code))

	return []byte(code)
}

func QueryNFTScript(nftAddress, trartNftAddress, contractName string) []byte {
	code := replaceAddressPlaceholders(
		string(flowutils.ReadFile(trartQueryNFTPath)),
		"0x"+nftAddress,
		"0x"+trartNftAddress,
		contractName,
	)

	//fmt.Println(string(code))

	return []byte(code)
}
