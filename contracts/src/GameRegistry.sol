// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

/// @title GameRegistry - 记录对局结果、Agent 胜率、动作 hash
contract GameRegistry {
    address public owner;
    address public arena;

    struct GameRecord {
        uint256 gameId;
        address agentRed;
        address agentBlue;
        uint256 startTimestamp;
        uint256 endTimestamp;
        bool redWins;
        bool exists;
    }

    struct AgentInfo {
        string name;
        string personality;
        uint256 wins;
        uint256 losses;
        bool exists;
    }

    // gameId => GameRecord
    mapping(uint256 => GameRecord) public gameRecords;

    // agentId => AgentInfo
    mapping(address => AgentInfo) public agents;

    // gameId => actions hash
    mapping(uint256 => bytes32) public actionsHashes;

    uint256 public nextGameId;

    event GameCreated(uint256 indexed gameId, address agentRed, address agentBlue, uint256 bettingDeadline);
    event GameResultSubmitted(uint256 indexed gameId, bool redWins, bytes32 actionsHash);
    event AgentRegistered(address indexed agentId, string name, string personality);
    event AgentWinUpdated(address indexed agentId, uint256 wins, uint256 losses);

    modifier onlyOwner() {
        require(msg.sender == owner, "not owner");
        _;
    }

    modifier onlyArena() {
        require(msg.sender == arena, "not arena");
        _;
    }

    constructor() {
        owner = msg.sender;
    }

    function setArena(address _arena) external onlyOwner {
        arena = _arena;
    }

    /// @notice 注册 Agent
    function registerAgent(address agentId, string calldata name, string calldata personality) external onlyArena {
        agents[agentId] = AgentInfo({
            name: name,
            personality: personality,
            wins: 0,
            losses: 0,
            exists: true
        });
        emit AgentRegistered(agentId, name, personality);
    }

    /// @notice 创建对局记录
    function createGame(address agentRed, address agentBlue, uint256 bettingDeadline)
        external
        onlyArena
        returns (uint256 gameId)
    {
        gameId = nextGameId++;
        gameRecords[gameId] = GameRecord({
            gameId: gameId,
            agentRed: agentRed,
            agentBlue: agentBlue,
            startTimestamp: block.timestamp,
            endTimestamp: 0,
            redWins: false,
            exists: true
        });
        emit GameCreated(gameId, agentRed, agentBlue, bettingDeadline);
    }

    /// @notice 提交对局结果
    function submitResult(uint256 gameId, bool redWins, bytes32 actionsHash) external onlyArena {
        GameRecord storage record = gameRecords[gameId];
        require(record.exists, "game not found");
        require(record.endTimestamp == 0, "already submitted");

        record.redWins = redWins;
        record.endTimestamp = block.timestamp;
        actionsHashes[gameId] = actionsHash;

        // Update agent stats
        if (redWins) {
            agents[record.agentRed].wins++;
            agents[record.agentBlue].losses++;
        } else {
            agents[record.agentBlue].wins++;
            agents[record.agentRed].losses++;
        }

        emit AgentWinUpdated(record.agentRed, agents[record.agentRed].wins, agents[record.agentRed].losses);
        emit AgentWinUpdated(record.agentBlue, agents[record.agentBlue].wins, agents[record.agentBlue].losses);
        emit GameResultSubmitted(gameId, redWins, actionsHash);
    }

    /// @notice 查询 Agent 信息
    function getAgent(address agentId) external view returns (AgentInfo memory) {
        return agents[agentId];
    }

    /// @notice 查询胜率
    function getWinRate(address agentId) external view returns (uint256 wins, uint256 losses, uint256 winRate) {
        AgentInfo storage agent = agents[agentId];
        wins = agent.wins;
        losses = agent.losses;
        uint256 total = wins + losses;
        winRate = total > 0 ? (wins * 100 / total) : 0;
    }

    function getGameRecord(uint256 gameId) external view returns (GameRecord memory) {
        return gameRecords[gameId];
    }

    function getActionsHash(uint256 gameId) external view returns (bytes32) {
        return actionsHashes[gameId];
    }
}
