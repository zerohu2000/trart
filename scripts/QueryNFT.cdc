// SPDX-License-Identifier: MIT

import TrartContractNFT from "../contracts/TrartTemplateNFT.cdc"

pub fun main(_ nftID: UInt64) : {String:String}? {
    return TrartContractNFT.getMetadataForCardID(cardID: nftID)?.data
}