import type { AuditEvent, Tournament } from "./course";

export type PresentationTone =
  | "neutral"
  | "success"
  | "warning"
  | "danger"
  | "info";

export interface TournamentPresentation {
  label: string;
  tone: PresentationTone;
  hint: string;
}

/** Describe the next business step using contract state and a trusted chain timestamp. */
export function getTournamentPresentation(
  t: Tournament,
  now: number,
): TournamentPresentation {
  if (t.status === "finalized") {
    return t.prize_claimed
      ? {
          label: "已领奖",
          tone: "success",
          hint: "冠军已领取奖金，可查看比赛证据与资金流水。",
        }
      : {
          label: "已结算 · 待领奖",
          tone: "success",
          hint: "审核已完成，冠军可用自己的钱包签名领取奖金。",
        };
  }
  if (t.status === "cancelled") {
    return {
      label: "已取消",
      tone: "neutral",
      hint: "赞助者可各自领取剩余可退的出资；赛事取消不代表所有退款已经领取。",
    };
  }

  if (!Number.isFinite(now) || now <= 0) {
    const labels = {
      open: "报名时间待确认",
      running: "对战中",
      result_submitted: "结果审核",
      disputed: "存在异议",
    };
    return {
      label: labels[t.status],
      tone: "neutral",
      hint: "正在等待网络时间，暂不能判断报名、审核窗口或最终期限是否已到。",
    };
  }

  if (now >= t.settlement_deadline) {
    return {
      label: "已到期 · 待退款",
      tone: "danger",
      hint: "最终结算期限已到；任意钱包先触发到期取消，确认取消后赞助者才能各自领取剩余退款。",
    };
  }

  switch (t.status) {
    case "open": {
      // Raw amounts must never be rounded through Number, including large prize pools.
      const funded =
        /^\d+$/.test(t.prize_pool) && BigInt(t.prize_pool) > BigInt(0);
      if (now < t.registration_deadline) {
        const registration =
          t.participants.length >= 2
            ? "参赛名额已满；截止前仍可赞助，已报名者可退出报名。"
            : "参赛者可选择自己的 Agent 报名，赞助者可向奖池出资。";
        return {
          label: "报名中",
          tone: "info",
          hint:
            registration +
            (funded ? "" : "当前奖池为零，开赛前需要非零赞助奖金。"),
        };
      }
      if (t.participants.length < 2) {
        return {
          label: "报名不足 · 可取消",
          tone: "warning",
          hint: "报名已截止且不足两位参赛者；任意钱包可触发到期取消，随后赞助者各自领取退款。",
        };
      }
      if (t.participants.length !== 2) {
        return {
          label: "报名信息待核对",
          tone: "danger",
          hint: "参赛人数不符合两人比赛规则，请先核对链上报名记录。",
        };
      }
      if (!funded) {
        return {
          label: "奖池不足 · 无法开赛",
          tone: "warning",
          hint: "报名已截止，但奖池没有有效的正数奖金，无法开赛；主办方可取消，或最终期限后由任意钱包取消。",
        };
      }
      if (now + t.review_period >= t.settlement_deadline) {
        return {
          label: "审核时间不足 · 无法开赛",
          tone: "warning",
          hint: "剩余时间已不足一个完整审核窗口，无法开赛；主办方可取消赛事，或最终期限到达后由任意钱包取消。",
        };
      }
      return {
        label: "等待开赛",
        tone: "info",
        hint: "报名已截止，两位参赛者和非零奖池已就绪，等待主办方签名开赛。",
      };
    }
    case "running":
      return {
        label: "对战中",
        tone: "info",
        hint: "主办方运行 Agent 对战并核对证据，再提交冠军与证据哈希；审核完成前不能领奖。",
      };
    case "result_submitted":
      return now < t.review_deadline
        ? {
            label: "审核窗口中",
            tone: "info",
            hint: "参赛者或实际赞助者可在窗口内提交异议；窗口结束后由指定审核员确认结果。",
          }
        : {
            label: "等待审核",
            tone: "info",
            hint: "异议窗口已结束，等待指定审核员核对证据并确认结果；确认后冠军才能领奖。",
          };
    case "disputed":
      return now < t.review_deadline
        ? {
            label: "存在异议",
            tone: "warning",
            hint: "异议已记录，等待审核窗口结束后由指定审核员根据证据裁决。",
          }
        : {
            label: "等待裁决",
            tone: "warning",
            hint: "审核窗口已结束，指定审核员需根据证据确认冠军或取消赛事；取消后赞助者可领取退款。",
          };
  }
}

/** Global contract roles do not imply registration in any particular tournament. */
export function globalRoleLabel(roles: string[]): string {
  const labels = [
    ["admin", "管理员"],
    ["organizer", "主办方"],
    ["auditor", "审核员"],
  ]
    .filter(([role]) => roles.includes(role))
    .map(([, label]) => label);
  return labels.length ? labels.join(" · ") : "普通用户";
}

/** Event participants and beneficiaries are not always the wallet that performed the action. */
export function getEventActor(event: AuditEvent): string {
  const fields: Record<string, string> = {
    TournamentCreated: "organizer",
    ParticipantRegistered: "participant",
    ParticipantWithdrawn: "participant",
    TournamentFunded: "sponsor",
    ResultSubmitted: "submitter",
    ResultDisputed: "challenger",
    ResultApproved: "auditor",
    DisputeResolved: "auditor",
    TournamentCancelled: "caller",
    PrizeClaimed: "winner",
    SponsorRefunded: "sponsor",
    RoleGranted: "sender",
    RoleRevoked: "sender",
  };
  const field = fields[event.event_name];
  const candidate = field ? event.data?.[field] : undefined;
  // Started and other events without an operator field cannot identify the caller.
  // The backend's generic actor may be a beneficiary, so it is not a safe fallback.
  return typeof candidate === "string" &&
    /^0x[0-9a-fA-F]{40}$/.test(candidate) &&
    !/^0x0{40}$/.test(candidate)
    ? candidate
    : "";
}
