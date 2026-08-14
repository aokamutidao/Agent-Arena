package engine

import (
	"crypto/sha256"
	"encoding/json"
	"math/rand"
)

// GameRunner 游戏运行器
type GameRunner struct {
	engine *Engine
}

// NewGameRunner 创建游戏运行器
func NewGameRunner() *GameRunner {
	return &GameRunner{engine: NewEngine()}
}

// RunFullGame 跑一场完整对局（使用随机动作模拟）
func (r *GameRunner) RunFullGame(gameID uint64, redName, blueName string) *GameState {
	game := r.engine.NewGame(gameID, redName, blueName, "berserker", "tactician")
	game.Status = StatusBetting
	game.Status = StatusPlaying

	for game.Status == StatusPlaying {
		redAction := r.randomAction(game.AgentRed, game.AgentBlue, game.Obstacles)
		blueAction := r.randomAction(game.AgentBlue, game.AgentRed, game.Obstacles)
		r.engine.ExecuteTurn(game, redAction, blueAction)
	}

	return game
}

// randomAction 生成一个合法的随机动作
func (r *GameRunner) randomAction(self, enemy AgentState, obstacles []Position) Action {
	dist := ManhattanDistance(self.Position, enemy.Position)

	// 收集合法动作
	var actions []Action

	// MOVE: 尝试几个随机目标
	for i := 0; i < 5; i++ {
		dx := rand.Intn(MOV*2+1) - MOV
		dy := rand.Intn(MOV*2+1) - MOV
		if abs(dx)+abs(dy) == 0 || abs(dx)+abs(dy) > MOV {
			continue
		}
		target := Position{
			X: clamp(int(self.Position.X)+dx, 0, MapWidth-1),
			Y: clamp(int(self.Position.Y)+dy, 0, MapHeight-1),
		}
		action := Action{Type: ActionMove, Target: target}
		if r.engine.ValidateAction(action, self, enemy, obstacles) == nil {
			actions = append(actions, action)
			break
		}
	}

	// ATTACK
	if dist <= AtkRange {
		actions = append(actions, Action{Type: ActionAttack})
	}

	// SKILL
	if dist <= SkillRange && self.SkillCooldown == 0 {
		actions = append(actions, Action{Type: ActionSkill})
	}

	// CHARGE (if not already charging and HP > 30)
	if !self.IsCharging && self.HP > 30 {
		actions = append(actions, Action{Type: ActionCharge})
	}

	// WAIT
	actions = append(actions, Action{Type: ActionWait})

	// 随机选一个
	return actions[rand.Intn(len(actions))]
}

// ComputeActionsHash 计算所有回合动作的 hash（用于链上验证）
func ComputeActionsHash(history []TurnRecord) [32]byte {
	data, _ := json.Marshal(history)
	return sha256.Sum256(data)
}

func clamp(val, min, max int) uint8 {
	if val < min {
		return uint8(min)
	}
	if val > max {
		return uint8(max)
	}
	return uint8(val)
}
