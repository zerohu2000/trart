import TrartContractNFT from 0xTRARTNFTADDRESS

pub fun main() : {String:String}? {
    return TrartContractNFT.getMetadataForCardID(cardID: %nftID)?.data
}