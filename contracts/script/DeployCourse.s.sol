// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {Script, console} from "forge-std/Script.sol";
import {TournamentManager} from "../src/TournamentManager.sol";
import {CourseToken} from "../src/CourseToken.sol";

/// @notice Explicit wallet deployment; this script does not give a backend any signing key.
/// @dev No broadcast happens unless a human invokes Forge with --broadcast.
contract DeployCourse is Script {
    function run() external {
        uint256 deployerKey = vm.envUint("COURSE_DEPLOYER_KEY");
        address admin = vm.envAddress("COURSE_ADMIN");
        address organizer = vm.envAddress("COURSE_ORGANIZER");
        address auditor = vm.envAddress("COURSE_AUDITOR");
        address token = vm.envOr("COURSE_TOKEN", address(0));
        require(vm.addr(deployerKey) == admin, "deployer must be explicit admin");
        require(organizer != address(0) && auditor != address(0) && organizer != auditor, "invalid roles");

        vm.startBroadcast(deployerKey);
        if (token == address(0)) token = address(new CourseToken(admin));
        TournamentManager manager = new TournamentManager(token, admin);
        manager.grantRole(manager.ORGANIZER_ROLE(), organizer);
        manager.grantRole(manager.AUDITOR_ROLE(), auditor);
        vm.stopBroadcast();

        console.log("Chain ID:", block.chainid);
        console.log("Deployment block:", block.number);
        console.log("CourseToken:", token);
        console.log("TournamentManager:", address(manager));
        console.log("TournamentRegistry:", address(manager.tournamentRegistry()));
        console.log("PrizeEscrow:", address(manager.prizeEscrow()));
        console.log("ResultReview:", address(manager.resultReview()));
        console.log("Admin:", admin);
        console.log("Organizer:", organizer);
        console.log("Auditor:", auditor);
    }
}
