import NonFungibleToken from 0xNFTADDRESS
import TrartContractNFT from 0xTRARTNFTADDRESS

pub struct NFTItem {
  pub let ID: UInt64
  pub let Metadata: {String:String}

  init(id: UInt64, metadata: {String:String}) {
    self.ID = id
    self.Metadata = metadata
  }
}

pub fun main() : [NFTItem] {

    let account1 = getAccount(0xUSERADDRESS)

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
