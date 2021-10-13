// SPDX-License-Identifier: MIT

import NonFungibleToken from "../contracts/NonFungibleToken.cdc"
import TrartContractNFT from "../contracts/TrartTemplateNFT.cdc"

transaction(recipient: Address, cardIDs: [UInt64], metadatas: [{String:String}]) {
    
    let minter: &TrartContractNFT.NFTMinter

    prepare(signer: AuthAccount) {

        self.minter = signer.borrow<&TrartContractNFT.NFTMinter>(from: TrartContractNFT.MinterStoragePath)
            ?? panic("Could not borrow a reference to the NFT minter")
    }

    execute {
 
        let recipient = getAccount(recipient)

        let receiver = recipient
            .getCapability(TrartContractNFT.CollectionPublicPath)!
            .borrow<&{TrartContractNFT.ICardCollectionPublic}>()
            ?? panic("Could not get receiver reference to the NFT Collection")

        //receiver.deposit(token: <- self.minter.newNFT(cardID: 1, data: {"Name": "Card1", "uri": "ipfs://1-123456789"}))

        let collection <- TrartContractNFT.createEmptyCollection()

        var i = 0
        for id in cardIDs {
            collection.deposit(token: <- self.minter.newNFT(cardID: id, data: metadatas[i]))
            i = i + 1
        }

        receiver.batchDeposit(tokens: <- collection)
    }
}
