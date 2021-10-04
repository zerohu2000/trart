import NonFungibleToken from 0xNFTADDRESS
import TrartContractNFT from 0xTRARTNFTADDRESS

transaction {
    prepare(signer: AuthAccount) {
        if signer.borrow<&TrartContractNFT.Collection>(from: TrartContractNFT.CollectionStoragePath) != nil {
            return
        }

        let collection <- TrartContractNFT.createEmptyCollection()
        signer.save(<-collection, to: TrartContractNFT.CollectionStoragePath)

        signer.link<&TrartContractNFT.Collection{NonFungibleToken.CollectionPublic, TrartContractNFT.ICardCollectionPublic}>(TrartContractNFT.CollectionPublicPath, target: TrartContractNFT.CollectionStoragePath)
    }
}
