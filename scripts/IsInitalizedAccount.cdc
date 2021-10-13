// SPDX-License-Identifier: MIT

import NonFungibleToken from "../contracts/NonFungibleToken.cdc"
import TrartContractNFT from "../contracts/TrartTemplateNFT.cdc"

pub fun main(_ address: Address) : UInt64 {

    let nftOwner = getAccount(address)

    let capability = nftOwner.getCapability<&{TrartContractNFT.ICardCollectionPublic}>(TrartContractNFT.CollectionPublicPath)

    let receiver = capability.borrow()??nil
    if receiver == nil {
        return 0
    }

    return 1
}