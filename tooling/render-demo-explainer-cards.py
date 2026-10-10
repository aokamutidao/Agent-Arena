#!/usr/bin/env python3
"""Render the reproducible English explainer cards used in the course demo.

The 1920 x 900 artwork sits above a separate 180 px bilingual caption area.
Requires Pillow. No source footage, application data, or network is modified.
"""
from __future__ import annotations

import argparse
import math
from pathlib import Path

from PIL import Image, ImageDraw, ImageFont

WIDTH, HEIGHT, SCALE = 1920, 900, 2
BG = "#F7FAF8"
WHITE = "#FFFFFF"
INK = "#162D35"
MUTED = "#526B70"
GREEN = "#167D63"
PALE = "#E8F4EE"
LINE = "#D9E5DF"
BLUE = "#2465A0"
BLUE_PALE = "#EAF2F9"
FONT_DIR = Path("/System/Library/Fonts/Supplemental")


def font(size: int, bold: bool = False):
    name = "Arial Bold.ttf" if bold else "Arial.ttf"
    return ImageFont.truetype(str(FONT_DIR / name), size * SCALE)


class Card:
    def __init__(self):
        self.image = Image.new("RGB", (WIDTH * SCALE, HEIGHT * SCALE), BG)
        self.draw = ImageDraw.Draw(self.image)

    @staticmethod
    def rect(box):
        return tuple(round(v * SCALE) for v in box)

    def box(self, x, y, w, h, fill=WHITE, outline=LINE, radius=22, width=2):
        self.draw.rounded_rectangle(self.rect((x, y, x + w, y + h)),
                                    radius=radius * SCALE, fill=fill,
                                    outline=outline, width=width * SCALE)

    def line(self, points, color=LINE, width=3):
        self.draw.line([(int(x*SCALE), int(y*SCALE)) for x, y in points],
                       fill=color, width=width*SCALE, joint="curve")

    def text(self, text, x, y, size=30, bold=False, color=INK, anchor="lt"):
        self.draw.text((x*SCALE, y*SCALE), text, font=font(size, bold),
                       fill=color, anchor=anchor)

    def wrap(self, text, x, y, width, size=28, bold=False, color=MUTED, leading=1.35):
        current, lines = "", []
        for word in text.split():
            trial = f"{current} {word}".strip()
            if self.draw.textlength(trial, font=font(size, bold)) > width*SCALE and current:
                lines.append(current)
                current = word
            else:
                current = trial
        if current:
            lines.append(current)
        for line in lines:
            self.text(line, x, y, size, bold, color)
            y += size*leading
        return y

    def pill(self, label, x, y, color=GREEN, fill=PALE, size=22):
        w = self.draw.textlength(label, font=font(size, True))/SCALE + 36
        self.box(x, y, w, 44, fill, fill, 22, 1)
        self.text(label, x+18, y+10, size, True, color)
        return w

    def arrow(self, points, color=GREEN, width=3):
        self.line(points, color, width)
        x0, y0 = points[-2]
        x1, y1 = points[-1]
        a = math.atan2(y1-y0, x1-x0)
        size = 12
        left = (x1-size*math.cos(a-.52), y1-size*math.sin(a-.52))
        right = (x1-size*math.cos(a+.52), y1-size*math.sin(a+.52))
        self.draw.polygon([(x1*SCALE,y1*SCALE),
                           (left[0]*SCALE,left[1]*SCALE),
                           (right[0]*SCALE,right[1]*SCALE)], fill=color)

    def header(self, section, title, subtitle):
        self.text(f"AGENT ARENA  /  {section.upper()}", 96, 104, 22, True, GREEN)
        self.text(title, 96, 151, 52, True)
        self.text(subtitle, 96, 223, 28, color=MUTED)

    def footer(self, text="LOCAL DEMONSTRATION  •  WALLET SIGNATURES  •  ON-CHAIN SETTLEMENT"):
        self.line([(96,807),(1824,807)], LINE, 2)
        self.text(text, 96, 771, 19, True, MUTED)

    def save(self, path):
        self.image.resize((WIDTH,HEIGHT), Image.Resampling.LANCZOS).save(path)


def title_card():
    c = Card()
    c.text("COURSE PROJECT DEMONSTRATION", 96, 112, 23, True, GREEN)
    c.text("Agent Arena", 96, 225, 94, True)
    c.text("AI tournaments.", 100, 362, 54, True)
    c.text("Transparent prize settlement.", 100, 430, 54, True)
    c.wrap("From a funded prize pool to a reviewed result and a winner’s claim.",
           100, 531, 1050, 32)
    c.pill("Local environment · Fast, focused demonstration", 100, 667, size=23)
    # A compact geometric arena motif, entirely vector drawn.
    c.box(1310,235,490,390,WHITE,LINE,36)
    c.box(1350,275,410,310,PALE,PALE,24)
    for i in range(1,5):
        c.line([(1350+i*82,275),(1350+i*82,585)], "#CFE3D8", 2)
    for i in range(1,4):
        c.line([(1350,275+i*77.5),(1760,275+i*77.5)], "#CFE3D8", 2)
    c.box(1391,315,64,64,"#C56059","#C56059",16)
    c.box(1658,477,64,64,BLUE,BLUE,16)
    c.arrow([(1469,347),(1545,347),(1545,507),(1641,507)],GREEN,5)
    c.text("STRATEGY → EVIDENCE → REWARD",1555,661,21,True,GREEN,"mt")
    c.footer()
    return c


def role_flow():
    c=Card()
    c.header("Business flow", "Five steps. One traceable tournament.",
             "Every participant sees the next step, the current state, and the financial record.")
    items=[("01","Organizer","Sets the schedule and assigns an independent auditor."),
           ("02","Two Agents","Participants register their saved Agent definitions."),
           ("03","Sponsor","Approves DUSD, then deposits it into the prize escrow."),
           ("04","Auditor","Reviews the match evidence and confirms the result."),
           ("05","Winner","Claims the approved prize through their own wallet.")]
    for i,(n,title,body) in enumerate(items):
        x=96+i*349
        c.box(x,329,330,281)
        c.pill(n,x+25,352,size=21)
        c.text(title,x+25,418,35,True)
        c.wrap(body,x+25,479,279,27)
        if i<4:c.arrow([(x+331,470),(x+346,470)],GREEN,3)
    c.box(96,655,1726,76,PALE,PALE,18)
    c.text("FUNDING",122,680,23,True,GREEN)
    c.text("Approve exact amount",310,680,26,True)
    c.arrow([(627,694),(677,694)])
    c.text("Deposit into escrow",702,680,26,True)
    c.arrow([(1002,694),(1052,694)])
    c.text("Independent review",1077,680,26,True)
    c.arrow([(1382,694),(1432,694)])
    c.text("Winner claims prize",1457,680,26,True)
    c.footer()
    return c


def architecture():
    c=Card()
    c.header("System architecture", "A full stack with clear responsibilities.",
             "The application connects the user experience, AI execution, persistent evidence, and contracts.")
    # Top row: UI, API, and AI execution.
    c.box(96,324,396,195)
    c.text("Next.js + Wallet",122,351,35,True)
    c.wrap("Dashboard, MetaMask, signed login, transaction previews",122,409,342,28)
    c.box(706,324,504,195)
    c.text("Go / Gin REST API",735,351,35,True)
    c.wrap("Authentication, tournament queries, match jobs and event indexing",735,409,446,28)
    c.box(1424,324,400,195)
    c.text("LLM + Game Engine",1448,351,32,True)
    c.wrap("Strategy actions, rule execution, round results",1450,409,346,28)
    c.arrow([(494,409),(704,409)])
    c.text("REST API",599,370,22,True,GREEN,"mt")
    c.arrow([(1212,409),(1422,409)])
    c.text("Run match",1317,370,22,True,GREEN,"mt")
    # Bottom row: blockchain and durable store.
    c.box(96,637,396,118,PALE,PALE)
    c.text("Solidity Contracts",122,660,33,True)
    c.text("Escrow + tournament state",122,708,25,color=MUTED)
    c.box(706,637,504,118,PALE,PALE)
    c.text("SQLite",735,660,33,True)
    c.text("Profiles, rounds, evidence, events",735,708,25,color=MUTED)
    c.arrow([(200,521),(200,635)],GREEN,4)
    c.text("Wallet-signed",226,550,23,True,GREEN)
    c.text("transactions",226,581,23,True,GREEN)
    c.arrow([(958,521),(958,635)],GREEN,3)
    c.text("Persist",985,565,23,True,GREEN)
    c.arrow([(706,485),(580,485),(580,697),(494,697)],MUTED,3)
    c.text("Read",605,561,23,True,MUTED)
    c.text("chain",605,592,23,True,MUTED)
    c.wrap("Wallets authorize fund movements. The API reads chain state and confirmed events.",
           1424,644,399,28,color=INK)
    c.footer()
    return c


def contracts():
    c=Card()
    c.header("Smart contracts", "Separate modules. Coordinated state changes.",
             "TournamentManager is the application entry point for tournament operations.")
    c.box(615,314,690,128,PALE,PALE)
    c.text("TournamentManager",960,340,39,True,GREEN,"mt")
    c.text("Roles · Lifecycle · Events",960,392,27,color=MUTED,anchor="mt")
    cards=[(96,"TournamentRegistry","Schedule, participants, Agent hashes and tournament state."),
           (694,"PrizeEscrow","Sponsor contributions, prize claims and individual refunds."),
           (1292,"ResultReview","Result evidence, review window and dispute decisions.")]
    c.line([(362,487),(1558,487)],GREEN,3)
    c.line([(960,444),(960,487)],GREEN,3)
    for x,title,body in cards:
        c.arrow([(x+266,487),(x+266,520)])
        c.box(x,522,532,176)
        c.text(title,x+27,550,33,True)
        c.wrap(body,x+27,603,478,27)
    c.pill("DUSD · ERC-20 demo token",695,710,size=21)
    c.text("Token transfer and accounting in PrizeEscrow",1292,719,22,color=MUTED)
    c.footer()
    return c


def agent_evidence():
    c=Card()
    c.header("Agent execution", "Every round becomes reviewable evidence.",
             "Saved strategy definitions and structured model actions feed a deterministic game engine.")
    items=[("01","Game state","Position, health and legal moves"),
           ("02","LLM actions","Each Agent chooses an action"),
           ("03","Game engine","Rules resolve movement and damage"),
           ("04","Saved rounds","Actions and state changes persist"),
           ("05","Replay + hash","Inspect rounds and commit evidence"),
           ("06","Result review","Auditor checks the proposed winner")]
    xs=[96,704,1312]
    for i,(n,title,body) in enumerate(items):
        # Snake layout expresses sequence without crossing paths.
        col=i%3 if i<3 else 2-i%3
        x=xs[col];y=318 if i<3 else 562
        c.box(x,y,510,172)
        c.pill(n,x+22,y+22,size=19)
        c.text(title,x+94,y+26,33,True)
        c.wrap(body,x+25,y+91,461,27)
    c.arrow([(608,404),(701,404)])
    c.arrow([(1216,404),(1309,404)])
    c.arrow([(1567,492),(1567,559)])
    c.arrow([(1309,648),(1216,648)])
    c.arrow([(701,648),(608,648)])
    c.footer("AGENT DEFINITIONS  •  STRUCTURED ACTIONS  •  PERSISTENT ROUNDS  •  CANONICAL EVIDENCE HASH")
    return c


def security():
    c=Card()
    c.header("Security design", "Controls follow identity, state, and money.",
             "Authorization is enforced by the API and contracts at the operation boundary.")
    controls=[("01","Wallet authentication","One-time login nonce and a wallet signature bind the session to the account."),
              ("02","Role and state checks","Organizer, participant and designated auditor permissions are checked on-chain."),
              ("03","Escrow accounting","Integer token units, SafeERC20 and reentrancy protection guard fund movements."),
              ("04","Claim protection","A finalized winner claims once. Sponsors claim their own contributions after cancellation.")]
    for i,(n,title,body) in enumerate(controls):
        x=96+(i%2)*885;y=319+(i//2)*217
        c.box(x,y,842,194)
        c.pill(n,x+25,y+27,size=20)
        c.text(title,x+107,y+34,35,True)
        c.wrap(body,x+28,y+100,783,28)
    c.footer()
    return c


def end_card():
    c=Card()
    c.header("Demonstration complete", "A funded tournament, carried through to payout.",
             "Tournament #13 · agent fight test 4 · Real model execution and signed local transactions")
    metrics=[("100","DUSD prize"),("15","completed rounds"),("30","successful model calls")]
    for i,(value,label) in enumerate(metrics):
        x=96+i*587
        c.box(x,334,553,221)
        c.text(value,x+35,374,82,True,GREEN)
        c.text(label,x+37,488,30,True)
    c.box(96,599,1727,129,BLUE_PALE,BLUE_PALE)
    c.text("BLUE WINS",132,626,39,True,BLUE)
    c.text("Result reviewed. Prize claimed. Complete financial history retained.",
           438,642,31,True)
    c.footer("AGENT ARENA  •  AI TOURNAMENTS  •  TRANSPARENT PRIZE SETTLEMENT")
    return c


def main():
    parser=argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--output",type=Path,
                        default=Path(__file__).resolve().parents[1]/".local/user-demo-edit/cards")
    args=parser.parse_args()
    args.output.mkdir(parents=True,exist_ok=True)
    makers=[("01-title",title_card),("02-role-flow",role_flow),
            ("03-architecture",architecture),("04-contracts",contracts),
            ("05-agent-evidence",agent_evidence),("06-security",security),
            ("07-end",end_card)]
    for name,maker in makers:
        path=args.output/f"{name}.png"
        maker().save(path)
        print(path)


if __name__=="__main__":
    main()
