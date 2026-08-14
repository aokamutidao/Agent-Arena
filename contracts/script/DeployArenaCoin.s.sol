// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "forge-std/Script.sol";
import "../src/ArenaCoin.sol";

contract DeployArenaCoin is Script {
    function run() external {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");

        vm.startBroadcast(deployerPrivateKey);

        ArenaCoin ac = new ArenaCoin();

        console.log("ArenaCoin deployed at:", address(ac));

        vm.stopBroadcast();
    }
}
