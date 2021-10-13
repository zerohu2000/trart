# TRART Template NFT

## What is the TRART

TRART is an application that provides users with the opportunity to purchase, collect, and showcase digital blockchain collectibles containing exclusive content from the international artists (collectively, the “Artists”).

## What are the contracts

- The template contracts are built with [Cadence](https://docs.onflow.org/cadence), Flow's resource-oriented smart contract programming language. 

- The contracts are template for NonFungibleToken of TRART. For each artist's artworks, TRART will use the template contracts to remake, and deploy a new smart contract on FLOW blockchain.

## What is the program

- The program source in the folder `go/test` is golang testing code. Run testing to check the contracts in conceived cases.

# ✨ Getting Started

## 1. Install Golang and the Flow CLI

- Before you start, install Golang and the [Flow command-line interface (CLI)](https://docs.onflow.org/flow-cli).

- See `/go/test/go.mod` to check Golang version in our environment.

- The testing program needs to run in linux os or mac os, please check for FLOW documents.

## 2. Clone the project

- Download the project files in your workspace folder.

## 3. Install dependencies

- Run `cd go/test` in your workspace folder.
- Run `go mod tidy` to install golang package dependencies.

## 4. Modify setting to yours.

- Modify the value fo variant `trartContractName` to your contract name in the file `go/test/trart_test.go`.

## 5. Testing the project

- Run `cd go/test` in your workspace folder.
- Run `go test`

# Others

## Market contract?

- There is no market contract in TRART. 

## flow.json

- The file `flow.json` is emulator-mode setting for CLI, and using original template contracts to deploy(contract name is `TrartContractNFT`). In release, TRART will remake the contracts to new contracts(other contract name). 
