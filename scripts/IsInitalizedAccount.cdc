import NonFungibleToken from 0xNFTADDRESS
import TrartContractNFT from 0xTRARTNFTADDRESS

pub fun main(_ address: Address) : UInt64 {

    let nftOwner = getAccount(address)

    let capability = nftOwner.getCapability<&{TrartContractNFT.ICardCollectionPublic}>(TrartContractNFT.CollectionPublicPath)

    let receiver = capability.borrow()??nil
    if receiver == nil {
        return 0
    }

    return 1
}