// SPDX-License-Identifier: MIT

import NonFungibleToken from "../contracts/NonFungibleToken.cdc"
import TrartContractNFT from "../contracts/TrartTemplateNFT.cdc"

transaction {
    prepare(signer: AuthAccount) {
        // Return early if the account already has a collection
        if signer.borrow<&TrartContractNFT.Collection>(from: TrartContractNFT.CollectionStoragePath) != nil {
            return
        }

        // create a new empty collection
        let collection <- TrartContractNFT.createEmptyCollection()

        // save it to the account
        signer.save(<-collection, to: TrartContractNFT.CollectionStoragePath)

        // create a public capability for the collection
        signer.link<&TrartContractNFT.Collection{NonFungibleToken.CollectionPublic, TrartContractNFT.ICardCollectionPublic}>(TrartContractNFT.CollectionPublicPath, target: TrartContractNFT.CollectionStoragePath)
    }
}
