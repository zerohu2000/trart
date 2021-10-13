// SPDX-License-Identifier: MIT

import NonFungibleToken from "../contracts/NonFungibleToken.cdc"
import TrartContractNFT from "../contracts/TrartTemplateNFT.cdc"

pub struct NFTItem {
  pub let ID: UInt64
  pub let Metadata: {String:String}

  init(id: UInt64, metadata: {String:String}) {
    self.ID = id
    self.Metadata = metadata
  }
}

pub fun main(_ address: Address) : [NFTItem] {

    // Get both public account objects
    let account1 = getAccount(address)

    // Find the public Receiver capability for their Collections
    let receiver1Ref = account1.getCapability(TrartContractNFT.CollectionPublicPath).borrow<&{TrartContractNFT.ICardCollectionPublic}>()
        ?? panic("Could not borrow account receiver reference:　TrartContractNFT.ICardCollectionPublic")

    let ids = receiver1Ref.getIDs()

    let ret : [NFTItem] = []
    
    for id in ids {
        if let nft = receiver1Ref.borrowCard(id: id) {

            if let metadata = nft.getCardMetadata() {
                ret.append(NFTItem(id: nft.id, metadata: metadata.data))
            }
        }
    }

    return ret
}
