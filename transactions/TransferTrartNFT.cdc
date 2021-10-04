import NonFungibleToken from 0xNFTADDRESS
import TrartContractNFT from 0xTRARTNFTADDRESS

transaction(recipient: Address, withdrawID: UInt64) {
    prepare(signer: AuthAccount) {
        
        let recipient = getAccount(recipient)

        let collectionRef = signer.borrow<&TrartContractNFT.Collection>(from: TrartContractNFT.CollectionStoragePath)
            ?? panic("Could not borrow a reference to the owner's collection")

        let depositRef = recipient.getCapability(TrartContractNFT.CollectionPublicPath)!.borrow<&{NonFungibleToken.CollectionPublic}>()!

        let nft <- collectionRef.withdraw(withdrawID: withdrawID)

        depositRef.deposit(token: <-nft)
    }
}

