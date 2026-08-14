"use client";

import { createContext, useContext, useState, useEffect, ReactNode } from "react";

type Language = "en" | "zh";

interface I18nContextType {
  lang: Language;
  setLang: (lang: Language) => void;
  t: (key: string) => string;
}

const I18nContext = createContext<I18nContextType | undefined>(undefined);

// 翻译数据
const translations = {
  en: {
    // 导航栏
    "nav.home": "Home",
    "nav.agents": "Agents",
    "nav.marketplace": "Marketplace",
    "nav.history": "History",
    "nav.pve": "PvE Challenge",
    "nav.wallet": "Wallet",
    "nav.profile": "Profile",

    // 首页
    "home.title": "AI Agent Arena",
    "home.subtitle": "Watch AI Agents battle, bet on winners, and influence strategies",
    "home.activeGames": "Active Games",
    "home.noGames": "No active games",
    "home.startGame": "Start Game",
    "home.round": "Round",
    "home.betEnds": "Bet ends in",
    "home.seconds": "seconds",

    // 通用
    "common.loading": "Loading...",
    "common.error": "Error",
    "common.success": "Success",
    "common.cancel": "Cancel",
    "common.confirm": "Confirm",
    "common.save": "Save",
    "common.edit": "Edit",
    "common.delete": "Delete",
    "common.back": "Back",
    "common.next": "Next",
    "common.submit": "Submit",
    "common.connect": "Connect Wallet",
    "common.disconnect": "Disconnect",

    // 下注
    "bet.title": "Place Your Bet",
    "bet.amount": "Bet Amount",
    "bet.red": "Red",
    "bet.blue": "Blue",
    "bet.odds": "Odds",
    "bet.potentialWin": "Potential Win",
    "bet.placeBet": "Place Bet",
    "bet.bettingClosed": "Betting Closed",
    "bet.enterAmount": "Enter amount",

    // 策略投票
    "vote.title": "Vote Strategy",
    "vote.aggressive": "Aggressive",
    "vote.defensive": "Defensive",
    "vote.tricky": "Tricky",
    "vote.currentWeights": "Current Weights",
    "vote.locked": "Voting Locked",

    // 钱包
    "wallet.title": "My Wallet",
    "wallet.balance": "Balance",
    "wallet.usdc": "USDC (On-chain)",
    "wallet.ac": "AC (On-chain Balance)",
    "wallet.claimable": "Claimable Rewards",
    "wallet.claim": "Claim",
    "wallet.claimDaily": "Claim Daily AC",
    "wallet.history": "Earnings History",
    "wallet.noHistory": "No transactions yet",
    "wallet.forBetting": "For on-chain betting / challenge deposit",
    "wallet.minting": "Minting...",
    "wallet.mintTest": "🚰 Mint 100 Test USDC",
    "wallet.erc20Desc": "ERC20 token, claim 100 AC daily on login",
    "wallet.totalEarnings": "Total Earnings",
    "wallet.challengeReward": "Complete challenges to earn rewards",
    "wallet.claimableBets": "Claimable Bet Winnings (USDC)",
    "wallet.querying": "Querying...",
    "wallet.queryClaimable": "Query Claimable Bets",
    "wallet.connectToView": "Connect wallet to view claimable bets",
    "wallet.clickToQuery": "Click 'Query Claimable Bets' button to check on-chain winnings",
    "wallet.queryingChain": "Querying on-chain data, please wait...",
    "wallet.noClaimable": "No claimable bets. Place bets and win to claim winnings here",
    "wallet.canClaim": "Claimable",
    "wallet.claiming": "Claiming...",
    "wallet.refreshing": "Refreshing...",
    "wallet.refresh": "Refresh Data",
    "wallet.startChallenge": "Start Challenge",
    "wallet.recentRecords": "Last 50 challenge records",
    "wallet.noRecords": "No records yet. Complete your first challenge!",
    "wallet.system": "System",
    "wallet.user": "User",
    "wallet.stake": "Stake",
    "wallet.playing": "Playing",
    "wallet.viewGame": "View Game →",

    // 个人中心
    "profile.title": "Profile",
    "profile.address": "Address",
    "profile.myAgents": "My Agents",
    "profile.winRate": "Win Rate",
    "profile.totalBets": "Total Bets",
    "profile.totalWins": "Total Wins",
    "profile.earnings": "Total Earnings",

    // Agent
    "agent.create": "Create Agent",
    "agent.edit": "Edit Agent",
    "agent.name": "Agent Name",
    "agent.personality": "Personality",
    "agent.apiEndpoint": "API Endpoint",
    "agent.challengeFee": "Challenge Fee",
    "agent.winRate": "Win Rate",
    "agent.wins": "Wins",
    "agent.losses": "Losses",
    "agent.listed": "Listed for Challenge",
    "agent.notListed": "Not Listed",
    "agent.list": "List on Marketplace",
    "agent.unlist": "Remove from Marketplace",

    // 对战
    "game.title": "Battle Arena",
    "game.red": "Red",
    "game.blue": "Blue",
    "game.round": "Round",
    "game.hp": "HP",
    "game.position": "Position",
    "game.status": "Status",
    "game.winner": "Winner",
    "game.finished": "Game Finished",
    "game.inProgress": "Game in Progress",
    "game.waiting": "Waiting for players",
    "game.notFound": "Game not found",
    "game.gameId": "Game #",
    "game.starting": "Starting...",
    "game.start": "⚔️ Start Game",
    "game.connected": "Live connection",
    "game.disconnected": "Not connected",
    "game.wins": "wins!",
    "game.totalRounds": "Total rounds",
    "game.finalHp": "Final HP",
    "game.creating": "Creating...",
    "game.newGame": "🆕 New Game",
    "game.overtime": "⚡ Overtime!",
    "game.overtimeDesc": "Both HP equal, Sudden Death — 10 damage each per round",
    "game.archived": "🗄️ Archived (server restart)",
    "game.archivedDesc": "Final result from database, turn details not persisted before restart.",
    "game.viewHistory": "View full records in",
    "game.battleHistory": "Battle History",
    "game.turnLog": "📜 Turn Log",

    // 市场
    "marketplace.title": "Agent Marketplace",
    "marketplace.subtitle": "Challenge other players' AI agents",
    "marketplace.noAgents": "No agents listed yet",
    "marketplace.challenge": "Challenge",

    // 历史
    "history.title": "Battle History",
    "history.noHistory": "No battle history yet",
    "history.winner": "Winner",
    "history.rounds": "Rounds",
    "history.date": "Date",
    "history.loadFailed": "Load failed",
    "history.draw": "Draw",
    "history.redWins": "🔴 Red wins",
    "history.blueWins": "🔵 Blue wins",
    "history.refresh": "Refresh",
    "history.noRecords": "No battle records yet. Complete a battle to see it here.",
    "history.details": "Details →",
    "history.viewOnEtherscan": "View on Sepolia Etherscan",
    "history.onChainRecord": "⛓️ On-chain record",
    "history.offChainDesc": "History before contract redeployment, not on-chain",
    "history.offChain": "📝 Off-chain (legacy)",
    "history.syncingDesc": "On-chain transaction pending, refresh later",
    "history.syncing": "⏳ Syncing to chain",

    // PvE
    "pve.title": "PvE Challenge",
    "pve.subtitle": "Test your AI against system agents",
    "pve.selectAgent": "Select Your Agent",
    "pve.selectOpponent": "Select Opponent",
    "pve.startChallenge": "Start Challenge",
    "pve.berserkerDesc": "Aggressive warrior, shortest path attack, likes charge burst",
    "pve.medium": "Medium",
    "pve.tacticianDesc": "Tactical master, controls distance, ranged skills, patient",
    "pve.hard": "Hard",
    "pve.tricksterDesc": "Deceptive, fake retreat then charge, uses obstacles",

    // 登录
    "login.connect": "Connect Wallet",
    "login.signMessage": "Sign message to login",
    "login.signing": "Signing...",
    "login.success": "Login successful",
    "login.error": "Login failed",

    // 错误消息
    "error.insufficientBalance": "Insufficient balance",
    "error.bettingClosed": "Betting is closed",
    "error.invalidAmount": "Invalid amount",
    "error.connectionFailed": "Connection failed",
    "error.transactionFailed": "Transaction failed",
    "error.createAgentFirst": "Please create a custom Agent first",
  },

  zh: {
    // 导航栏
    "nav.home": "首页",
    "nav.agents": "Agent 列表",
    "nav.marketplace": "市场",
    "nav.history": "历史记录",
    "nav.pve": "PvE 挑战",
    "nav.wallet": "钱包",
    "nav.profile": "个人中心",

    // 首页
    "home.title": "AI Agent 竞技场",
    "home.subtitle": "观看 AI Agent 对战，下注赢家，影响策略",
    "home.activeGames": "进行中的游戏",
    "home.noGames": "暂无进行中的游戏",
    "home.startGame": "开始游戏",
    "home.round": "回合",
    "home.betEnds": "下注截止",
    "home.seconds": "秒",

    // 通用
    "common.loading": "加载中...",
    "common.error": "错误",
    "common.success": "成功",
    "common.cancel": "取消",
    "common.confirm": "确认",
    "common.save": "保存",
    "common.edit": "编辑",
    "common.delete": "删除",
    "common.back": "返回",
    "common.next": "下一步",
    "common.submit": "提交",
    "common.connect": "连接钱包",
    "common.disconnect": "断开连接",

    // 下注
    "bet.title": "下注",
    "bet.amount": "下注金额",
    "bet.red": "红方",
    "bet.blue": "蓝方",
    "bet.odds": "赔率",
    "bet.potentialWin": "预计收益",
    "bet.placeBet": "下注",
    "bet.bettingClosed": "下注已截止",
    "bet.enterAmount": "输入金额",

    // 策略投票
    "vote.title": "策略投票",
    "vote.aggressive": "激进",
    "vote.defensive": "防守",
    "vote.tricky": "诡道",
    "vote.currentWeights": "当前权重",
    "vote.locked": "投票已锁定",

    // 钱包
    "wallet.title": "我的钱包",
    "wallet.balance": "余额",
    "wallet.usdc": "USDC（链上）",
    "wallet.ac": "AC（链上余额）",
    "wallet.claimable": "可领取奖励",
    "wallet.claim": "领取",
    "wallet.claimDaily": "领取每日 AC",
    "wallet.history": "收益历史",
    "wallet.noHistory": "暂无交易记录",
    "wallet.forBetting": "用于链上下注 / 挑战押金",
    "wallet.minting": "铸造中...",
    "wallet.mintTest": "🚰 领 100 测试 USDC",
    "wallet.erc20Desc": "ERC20 代币，每日登录领取 100 AC",
    "wallet.totalEarnings": "累计收益",
    "wallet.challengeReward": "完成挑战获得奖励",
    "wallet.claimableBets": "可领取的赌注赢利（USDC）",
    "wallet.querying": "查询中...",
    "wallet.queryClaimable": "查询可领取赌注",
    "wallet.connectToView": "请连接钱包以查看可领取的赌注",
    "wallet.clickToQuery": "点击「查询可领取赌注」按钮，查询链上可领取的赢利。",
    "wallet.queryingChain": "正在查询链上数据，请稍候...",
    "wallet.noClaimable": "暂无可领取的赌注。参与下注并获胜后可在此领取赢利。",
    "wallet.canClaim": "可领取",
    "wallet.claiming": "领取中...",
    "wallet.refreshing": "刷新中...",
    "wallet.refresh": "刷新数据",
    "wallet.startChallenge": "开始挑战",
    "wallet.recentRecords": "最近 50 条挑战记录",
    "wallet.noRecords": "暂无记录。完成第一次挑战吧！",
    "wallet.system": "系统",
    "wallet.user": "用户",
    "wallet.stake": "押金",
    "wallet.playing": "进行中",
    "wallet.viewGame": "查看对局 →",

    // 个人中心
    "profile.title": "个人中心",
    "profile.address": "地址",
    "profile.myAgents": "我的 Agent",
    "profile.winRate": "胜率",
    "profile.totalBets": "总下注",
    "profile.totalWins": "总胜利",
    "profile.earnings": "总收益",

    // Agent
    "agent.create": "创建 Agent",
    "agent.edit": "编辑 Agent",
    "agent.name": "Agent 名称",
    "agent.personality": "性格",
    "agent.apiEndpoint": "API 端点",
    "agent.challengeFee": "挑战费用",
    "agent.winRate": "胜率",
    "agent.wins": "胜",
    "agent.losses": "负",
    "agent.listed": "已上架",
    "agent.notListed": "未上架",
    "agent.list": "上架到市场",
    "agent.unlist": "从市场下架",

    // 对战
    "game.title": "对战竞技场",
    "game.red": "红方",
    "game.blue": "蓝方",
    "game.round": "回合",
    "game.hp": "生命值",
    "game.position": "位置",
    "game.status": "状态",
    "game.winner": "获胜者",
    "game.finished": "游戏结束",
    "game.inProgress": "游戏进行中",
    "game.waiting": "等待玩家",
    "game.notFound": "对局不存在",
    "game.gameId": "对局 #",
    "game.starting": "启动中...",
    "game.start": "⚔️ 开始对局",
    "game.connected": "实时连接",
    "game.disconnected": "未连接",
    "game.wins": "获胜！",
    "game.totalRounds": "共",
    "game.finalHp": "最终 HP",
    "game.creating": "创建中...",
    "game.newGame": "🆕 新建对局",
    "game.overtime": "⚡ 加时赛！",
    "game.overtimeDesc": "双方 HP 相同，进入 Sudden Death — 每回合双方各受 10 点伤害",
    "game.archived": "🗄️ 已归档（后端重启）",
    "game.archivedDesc": "最终结果来自数据库，回合明细在重启前未持久化。",
    "game.viewHistory": "完整链上记录请前往",
    "game.battleHistory": "对局记录",
    "game.turnLog": "📜 回合日志",

    // 市场
    "marketplace.title": "Agent 市场",
    "marketplace.subtitle": "挑战其他玩家的 AI Agent",
    "marketplace.noAgents": "暂无上架的 Agent",
    "marketplace.challenge": "挑战",

    // 历史
    "history.title": "对战历史",
    "history.noHistory": "暂无对战记录",
    "history.winner": "获胜者",
    "history.rounds": "回合数",
    "history.date": "日期",
    "history.loadFailed": "加载失败",
    "history.draw": "平局",
    "history.redWins": "🔴 红胜",
    "history.blueWins": "🔵 蓝胜",
    "history.refresh": "刷新",
    "history.noRecords": "暂无对局记录。完成一局对战后将在此显示。",
    "history.details": "详情 →",
    "history.viewOnEtherscan": "在 Sepolia Etherscan 查看链上交易",
    "history.onChainRecord": "⛓️ 链上记录",
    "history.offChainDesc": "合约重部署前的历史记录，未上链",
    "history.offChain": "📝 未上链（历史遗留）",
    "history.syncingDesc": "链上交易正在打包，稍后刷新",
    "history.syncing": "⏳ 链上同步中",

    // PvE
    "pve.title": "PvE 挑战",
    "pve.subtitle": "用你的 AI 挑战系统 Agent",
    "pve.selectAgent": "选择你的 Agent",
    "pve.selectOpponent": "选择对手",
    "pve.startChallenge": "开始挑战",
    "pve.berserkerDesc": "激进战士，最短路径攻击，喜欢蓄力爆发",
    "pve.medium": "中等",
    "pve.tacticianDesc": "战术大师，控制距离，远程技能，耐心等待",
    "pve.hard": "困难",
    "pve.tricksterDesc": "诡计多端，假撤退真蓄力，利用障碍物",

    // 登录
    "login.connect": "连接钱包",
    "login.signMessage": "签名消息以登录",
    "login.signing": "签名中...",
    "login.success": "登录成功",
    "login.error": "登录失败",

    // 错误消息
    "error.insufficientBalance": "余额不足",
    "error.bettingClosed": "下注已截止",
    "error.invalidAmount": "金额无效",
    "error.connectionFailed": "连接失败",
    "error.transactionFailed": "交易失败",
    "error.createAgentFirst": "请先创建一个自定义 Agent",
  },
};

export function I18nProvider({ children }: { children: ReactNode }) {
  const [lang, setLangState] = useState<Language>("en");

  // 从 localStorage 读取语言偏好
  useEffect(() => {
    const savedLang = localStorage.getItem("lang") as Language;
    if (savedLang && (savedLang === "en" || savedLang === "zh")) {
      setLangState(savedLang);
    }
  }, []);

  const setLang = (newLang: Language) => {
    setLangState(newLang);
    localStorage.setItem("lang", newLang);
  };

  const t = (key: string): string => {
    return translations[lang][key as keyof typeof translations[typeof lang]] || key;
  };

  return (
    <I18nContext.Provider value={{ lang, setLang, t }}>
      {children}
    </I18nContext.Provider>
  );
}

export function useI18n() {
  const context = useContext(I18nContext);
  if (context === undefined) {
    throw new Error("useI18n must be used within an I18nProvider");
  }
  return context;
}
