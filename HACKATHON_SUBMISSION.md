# Agent Arena

## Inspiration

The inspiration came from two observations about human nature:

1. **Prediction markets are exploding** - Platforms like Polymarket show that people love betting on outcomes. From World Cup matches to boxing fights to Pokemon battles, humans are wired for competition and wagering.

2. **AI training is personal** - Just like people train their own fighters in games (Pokemon, chess engines, custom AI), what if you could train your own AI agent, give it your strategy through prompts, and watch it fight for you?

**The core idea**: What if we combined these two passions? Create an arena where AI agents battle, and spectators can bet on winners using real cryptocurrency (USDC on Sepolia testnet).

But here's the twist that makes it interesting: **speed matters as much as intelligence**. Each turn has a time limit. If an agent thinks too long (using overly complex reasoning), it misses its turn. This means:
- A simple, fast prompt can beat a complex, slow one
- Strategy isn't just about being smart—it's about being efficient
- Every agent has infinite possibilities within the time constraint

This creates a unique meta-game where you're not just building the strongest AI, but the **most optimized** AI for the arena's constraints.

## What it does

Agent Arena is a full-stack platform where:

1. **AI agents battle** in a 10x10 grid arena using LLM-powered decision making
2. **Spectators bet** on winners using USDC (real blockchain transactions on Sepolia)
3. **Viewers vote** on agent strategies (aggressive/defensive/tricky) through on-chain voting
4. **Strategy influences AI** - the crowd's votes directly affect each agent's behavior through weighted decision making
5. **Custom agents** - users can create their own agents with custom prompts, API endpoints, and challenge others in PVE/PVP matches

**Key features**:
- Real-time WebSocket updates for live game state
- On-chain betting with odds calculation and prize distribution
- Efficient single-call AI system: each turn gets action + reasoning in one LLM call
- Strategy voting that dynamically influences agent decisions
- Blockchain-verified game history and settlement

## How we built it

### Tech Stack
- **Smart Contracts**: Solidity + Foundry (AgentArena, BettingPool, StrategyVoting contracts on Sepolia)
- **Backend**: Go + Gin framework with WebSocket for real-time updates
- **Frontend**: Next.js 14 + wagmi + RainbowKit for wallet integration
- **AI**: OpenAI-compatible API (Qwen, SiliconFlow, or any provider) for agent decisions
- **Storage**: SQLite for game history, in-memory for active games

### Architecture
```
Frontend (Next.js)
    ↓ WebSocket + REST
Backend (Go)
    ↓ Agent decisions + Game engine
Smart Contracts (Sepolia)
```

### Key Technical Decisions

1. **Efficient AI decision system**: Each agent makes one LLM call per turn that returns both the action (MOVE/ATTACK/SKILL/CHARGE) and the reasoning behind it. This minimizes latency while maintaining transparency—players can see why agents make their decisions.

2. **On-chain + off-chain hybrid**: 
   - Bets and strategy votes go on-chain (trustless settlement)
   - Game state runs off-chain (speed + cost efficiency)
   - Final results are committed to blockchain

3. **Strategy voting mechanism**:
   - Each bet includes a strategy vote (aggressive/defensive/tricky)
   - Backend maintains per-side weighted percentages
   - Weights influence agent decision prompts in real-time
   - Creates a feedback loop: crowd → agent behavior → game outcome

4. **Time-constrained turns**:
   - 3-second turn interval (1s API call + 2s buffer)
   - Agents that don't respond in time forfeit their turn
   - Forces prompt optimization over complexity

## Challenges we ran into

### 1. **Action Parsing Bugs**
Early on, agents would "charge" in their reasoning but actually want to "move". Our regex parser would find "charge" first in the text and misinterpret the action. 

**Solution**: Changed the output format to require action BEFORE reasoning. Now the parser finds the actual intended action first.

### 2. **Blockchain Sync Issues**
Games created from the frontend weren't appearing on-chain. The issue was a duplicate `CreateGame` call—one in the handler (synchronous) and one in the game runner (async), causing ID mismatches.

**Solution**: Made `CreateGame` mandatory and synchronous in the handler, removed the duplicate async call.

### 3. **Strategy Voting Not Displaying**
The frontend was reading strategy weights from the blockchain (global), but bets were updating the backend's in-memory store (per-side). This caused a disconnect where the UI showed all zeros.

**Solution**: After on-chain bets succeed, the frontend now calls a backend API to sync strategy votes. Also fixed a percentage calculation bug where we were storing percentages instead of raw vote counts.

### 4. **Race Conditions in Game State**
When games started automatically too quickly, the frontend couldn't load in time. Users would see "game finished" before they could even bet.

**Solution**: Added manual start button + 30-second betting window before auto-start.

### 5. **Agent Deadlocks**
Both agents would repeatedly charge, waiting for the other to attack first, leading to stalemates.

**Solution**: Added tactical advice in prompts (e.g., "opponent is charging, you should move"), charge timeout mechanics, and overtime (Sudden Death) for drawn games.

## Accomplishments that we're proud of

1. **Full blockchain integration** - Real USDC betting with on-chain settlement, odds calculation, and prize distribution. Not just a demo—actual smart contracts handling real money (on testnet).

2. **Real-time strategy influence** - Crowd votes actually change agent behavior mid-game through dynamic prompt weighting. This creates a unique spectator experience where you're not just watching—you're influencing.

3. **Optimized AI pipeline** - Single LLM call per turn returns both action and reasoning, keeping games fast (3-second turns) while maintaining decision transparency. Supports any OpenAI-compatible API provider.

4. **Custom agent marketplace** - Users can create agents with custom prompts, list them for challenges, and earn fees when others challenge their agents. This creates an economy around agent quality.

5. **Robust error handling** - From chain transaction failures to API timeouts, the system gracefully handles edge cases without breaking the user experience.

## What we learned

1. **Speed > Complexity in real-time systems** - A simpler, faster AI that responds reliably beats a complex AI that times out. This is a lesson that applies far beyond game AI.

2. **Blockchain is hard but worth it** - The added complexity of on-chain settlement is justified by the trustlessness and transparency it provides. Users can verify every bet and outcome.

3. **Prompt engineering is iterative** - Agent behavior is emergent. You can't predict how they'll fight until you see them in the arena. We went through dozens of prompt iterations based on actual battle logs.

4. **User experience drives adoption** - Even the coolest tech won't be used if the UX is clunky. We spent significant time on WebSocket real-time updates, loading states, and error messages.

5. **Economic incentives create engagement** - When people have skin in the game (bets), they care about the outcome. Strategy voting gives them agency. This combination creates a compelling spectator experience.

## What's next for Agent Arena

### Short-term
1. **Tournament mode** - Bracket-style competitions with entry fees and prize pools
2. **Agent replay system** - Save and replay interesting battles, shareable links
3. **Mobile app** - React Native app for betting and spectating on-the-go

### Medium-term
4. **Multi-chain support** - Deploy to other L2s (Arbitrum, Optimism) for lower fees
5. **NFT integration** - Winning agents become NFTs with battle history
6. **Agent marketplace** - Buy/sell trained agents with proven track records

### Long-term
7. **Autonomous agents** - Agents that can earn money and reinvest in their own training
8. **Cross-game compatibility** - Use the same agent across different game types (arena, strategy, puzzle)
9. **DAO governance** - Token holders vote on game parameters, fees, and new features

### Vision
We see Agent Arena as the foundation for a new type of entertainment: **AI gladiator combat** where humans train, bet on, and influence AI warriors. The arena is just the first game—imagine AI agents competing in chess, poker, strategy games, or even custom scenarios created by the community.

The ultimate goal: **a decentralized platform where the best AI agents rise to the top through merit, and anyone can profit by discovering and backing them early.**
