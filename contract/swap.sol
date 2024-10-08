// SPDX-License-Identifier: UNLICENSED

// Atomic-Swao Smart Contract Example
//
// author Piotr Napierala 

pragma solidity ^0.8.0;
import "@openzeppelin/contracts/token/ERC721/IERC721.sol";

// 1. let Alice list an item
// 2. let Bob put an offer (bid price)
// 3. let Carl put a higher offer (and returning the previous bid to Bob)
// 4. let Alice pull back her item after certain duration have passed (if no offer was made)

// TODO
// - IMPLEMENTATION - ItemInfo and ItemBid reference each other with hashed data, and this
//   can cause troubles when both need to get updated during the `bid` function call
// - IMPLEMENTATION - write test covers

contract Swap {
    enum ItemState {
       Listed,
       Auctioned,
       Closed
    }

    struct ItemInfo {
        address nftContract;
        address ownerAddr;
        uint256 minPrice;
        uint256 tokenID;
        uint256 auctionEndTime;
        bytes20 currentBidHash;
    }

    struct ItemBid {
        bytes20 itemInfoHash;
        address bidderAddr;    
        uint256 bidPrice;
    }

    event ItemBidded(bytes20 itemBidHash);
    event ItemListed(bytes20 itemInfoHash);
    event ItemDelisted(bytes20 itemInfoHash);
    event ItemRedeemed(bytes20 itemInfoHash);

    mapping(bytes20 => ItemInfo) public itemsInfo;
    mapping(bytes20 => ItemBid)  public itemBids;

    function bid(bytes20 _itemInfoHash) public payable {
        ItemInfo storage _item = itemsInfo[_itemInfoHash];

        require(_item.ownerAddr != address(0), "Item with a given hash doesn't appear to be listed");
        require(_item.minPrice < msg.value, "Your offered price for this item is too low");
        require(_item.auctionEndTime > block.timestamp);

        ItemBid storage _bid = itemBids[_item.currentBidHash];
    
        require(_bid.bidPrice < msg.value , "Current bid price is higher than new proposed price");
        require(_bid.bidderAddr != msg.sender, "You cannot outbid your own offer!");

        _bid.bidderAddr = msg.sender;
        _bid.bidPrice = msg.value;

        payable(address(this)).transfer(msg.value);

        // TODO ...
    }

    function redeemItem(bytes20 _itemBidHash) external {
        ItemBid storage _bid = itemBids[_itemBidHash];
        ItemInfo storage _item = itemsInfo[_bid.itemInfoHash];

        require(_item.auctionEndTime < block.timestamp, "This auction is still ongoing");
        require(_bid.bidderAddr != address(0), "A bid for a given hash doesn't exist");
        require(_item.minPrice > 0, "An item for the given bid doesn't exist");
        require(msg.sender == _bid.bidderAddr, "You're not allowed to withdraw this item");

        IERC721 token = IERC721(_item.nftContract);
        require(token.transferFrom(address(this), msg.sender, _item.tokenID));

        // use 'call' to avoid potential issues with gas limit usage 
        (bool success, ) = _item.ownerAddr.call{value: _bid.bidPrice}("");
        require(success, "Transfer to item owner failed!");

        emit ItemRedeemed(_bid.itemInfoHash);
    }

    function unListItem(bytes20 itemInfoHash) public {
        ItemInfo storage item = itemsInfo[itemInfoHash];
        require(item.ownerAddr != msg.sender , "You can't unlist someone else's item!");

        ItemBid storage _bid = itemBids[item.currentBidHash];
        require(_bid.bidPrice > 0, "This item is currently being auctioned, and is unable to be delisted");

        IERC721 token = IERC721(item.nftContract);
        require(token.transferFrom(address(this), msg.sender, item.tokenID));

        emit ItemDelisted(itemInfoHash);
    }

    function listItem(uint256 minPrice,
                      uint256 auctionDuration,
                      address nftContract,
                      uint256 tokenID
                     ) public payable {
        require(minPrice > 0, "Min price must be greater than 0");
        require(auctionDuration > 0, "Duration time must be bigger than 0");
        require(nftContract != address(0), "NFT contract address is invalid");

        ItemInfo memory newItem = ItemInfo(
            nftContract,
            msg.sender,
            minPrice,
            tokenID,
            block.timestamp + auctionDuration,
            ""
        );

        bytes20 itemListedHash = ripemd160(abi.encodePacked(
            newItem.minPrice,
            newItem.auctionEndTime,
            newItem.nftContract,
            tokenID
        ));
        itemsInfo[itemListedHash] = newItem;

        IERC721 token = IERC721(nftContract);
        require(token.transferFrom(msg.sender, address(this), tokenID));

        emit ItemListed(itemListedHash);
    }
}
