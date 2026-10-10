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
  lang: "zh" | "en" = "zh",
): TournamentPresentation {
  const l = (zh: string, en: string) => (lang === "zh" ? zh : en);
  if (t.status === "finalized") {
    return t.prize_claimed
      ? {
          label: l("已领奖", "Prize claimed"),
          tone: "success",
          hint: l(
            "冠军已领取奖金，可查看比赛证据与资金流水。",
            "The winner has claimed the prize. Match evidence and fund records are available.",
          ),
        }
      : {
          label: l("已结算 · 待领奖", "Settled \u00b7 Ready to claim"),
          tone: "success",
          hint: l(
            "审核已完成，冠军可用自己的钱包签名领取奖金。",
            "Review is complete. The winner can sign with their own wallet to claim the prize.",
          ),
        };
  }
  if (t.status === "cancelled") {
    return {
      label: l("已取消", "Cancelled"),
      tone: "neutral",
      hint: l(
        "赞助者可各自领取剩余可退的出资；赛事取消不代表所有退款已经领取。",
        "Each sponsor can claim their remaining refundable contribution. Cancellation does not mean all refunds have been claimed.",
      ),
    };
  }

  if (!Number.isFinite(now) || now <= 0) {
    const labels = {
      open: l("报名时间待确认", "Checking registration time"),
      running: l("对战中", "Match in progress"),
      result_submitted: l("结果审核", "Result review"),
      disputed: l("存在异议", "Disputed"),
    };
    return {
      label: labels[t.status],
      tone: "neutral",
      hint: l(
        "正在等待网络时间，暂不能判断报名、审核窗口或最终期限是否已到。",
        "Waiting for network time to check registration, review and settlement deadlines.",
      ),
    };
  }

  if (now >= t.settlement_deadline) {
    return {
      label: l("已到期 · 待退款", "Expired \u00b7 Refund pending"),
      tone: "danger",
      hint: l(
        "最终结算期限已到；任意钱包先触发到期取消，确认取消后赞助者才能各自领取剩余退款。",
        "The settlement deadline has passed. Any wallet must first trigger cancellation; after confirmation, each sponsor can claim their remaining refund.",
      ),
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
            ? l(
                "参赛名额已满；截止前仍可赞助，已报名者可退出报名。",
                "Both places are filled. Sponsorship and withdrawals remain available before registration closes.",
              )
            : l(
                "参赛者可选择自己的 Agent 报名，赞助者可向奖池出资。",
                "Participants can register their own Agent, and sponsors can fund the prize pool.",
              );
        return {
          label: l("报名中", "Registration open"),
          tone: "info",
          hint:
            registration +
            (funded
              ? ""
              : l(
                  "当前奖池为零，开赛前需要非零赞助奖金。",
                  " The prize pool is empty. Funding is required before the match can start.",
                )),
        };
      }
      if (t.participants.length < 2) {
        return {
          label: l("报名不足 · 可取消", "Underfilled \u00b7 Can cancel"),
          tone: "warning",
          hint: l(
            "报名已截止且不足两位参赛者；任意钱包可触发到期取消，随后赞助者各自领取退款。",
            "Registration has closed with fewer than two participants. Any wallet can trigger cancellation, then sponsors can claim their refunds.",
          ),
        };
      }
      if (t.participants.length !== 2) {
        return {
          label: l("报名信息待核对", "Check registrations"),
          tone: "danger",
          hint: l(
            "参赛人数不符合两人比赛规则，请先核对链上报名记录。",
            "The participant count does not meet the two-player rule. Check on-chain registrations first.",
          ),
        };
      }
      if (!funded) {
        return {
          label: l("奖池不足 · 无法开赛", "No funding \u00b7 Cannot start"),
          tone: "warning",
          hint: l(
            "报名已截止，但奖池没有有效的正数奖金，无法开赛；主办方可取消，或最终期限后由任意钱包取消。",
            "Registration has closed without a valid positive prize pool. The organizer can cancel, or any wallet can cancel after the settlement deadline.",
          ),
        };
      }
      if (now + t.review_period >= t.settlement_deadline) {
        return {
          label: l(
            "审核时间不足 · 无法开赛",
            "Review time too short \u00b7 Cannot start",
          ),
          tone: "warning",
          hint: l(
            "剩余时间已不足一个完整审核窗口，无法开赛；主办方可取消赛事，或最终期限到达后由任意钱包取消。",
            "There is not enough time for a full review window. The organizer can cancel, or any wallet can cancel after the settlement deadline.",
          ),
        };
      }
      return {
        label: l("等待开赛", "Ready to start"),
        tone: "info",
        hint: l(
          "报名已截止，两位参赛者和非零奖池已就绪，等待主办方签名开赛。",
          "Registration is closed, with two participants and a funded prize pool. The organizer can sign to start the match.",
        ),
      };
    }
    case "running":
      return {
        label: l("对战中", "Match in progress"),
        tone: "info",
        hint: l(
          "主办方运行 Agent 对战并核对证据，再提交冠军与证据哈希；审核完成前不能领奖。",
          "The organizer runs the Agent match, checks its evidence, then submits the winner and evidence hash. Prizes cannot be claimed before approval.",
        ),
      };
    case "result_submitted":
      return now < t.review_deadline
        ? {
            label: l("审核窗口中", "Review window open"),
            tone: "info",
            hint: l(
              "参赛者或实际赞助者可在窗口内提交异议；窗口结束后由指定审核员确认结果。",
              "Participants and contributing sponsors can dispute the result during this window. The assigned auditor can approve it after the window closes.",
            ),
          }
        : {
            label: l("等待审核", "Awaiting approval"),
            tone: "info",
            hint: l(
              "异议窗口已结束，等待指定审核员核对证据并确认结果；确认后冠军才能领奖。",
              "The dispute window is closed. The assigned auditor must check the evidence and approve the result before the winner can claim.",
            ),
          };
    case "disputed":
      return now < t.review_deadline
        ? {
            label: l("存在异议", "Disputed"),
            tone: "warning",
            hint: l(
              "异议已记录，等待审核窗口结束后由指定审核员根据证据裁决。",
              "A dispute has been recorded. The assigned auditor can resolve it after the review window closes.",
            ),
          }
        : {
            label: l("等待裁决", "Awaiting resolution"),
            tone: "warning",
            hint: l(
              "审核窗口已结束，指定审核员需根据证据确认冠军或取消赛事；取消后赞助者可领取退款。",
              "The review window is closed. The assigned auditor must confirm a winner or cancel based on the evidence. Sponsors can claim refunds after cancellation.",
            ),
          };
  }
}

/** Global contract roles do not imply registration in any particular tournament. */
export function globalRoleLabel(
  roles: string[],
  lang: "zh" | "en" = "zh",
): string {
  const l = (zh: string, en: string) => (lang === "zh" ? zh : en);
  const labels = [
    ["admin", l("管理员", "Admin")],
    ["organizer", l("主办方", "Organizer")],
    ["auditor", l("审核员", "Auditor")],
  ]
    .filter(([role]) => roles.includes(role))
    .map(([, label]) => label);
  return labels.length ? labels.join(" · ") : l("普通用户", "Standard user");
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
