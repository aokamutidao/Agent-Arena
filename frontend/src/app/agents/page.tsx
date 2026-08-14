import { Card, CardContent } from "@/components/ui/card";

export default function AgentsPage() {
  const agents = [
    {
      name: "Berserker",
      personality: "狂战士",
      description: "最好的防守就是进攻。",
      wins: 12,
      losses: 5,
      winRate: 70.6,
    },
    {
      name: "Tactician",
      personality: "战术家",
      description: "耐心是胜利的关键。",
      wins: 8,
      losses: 9,
      winRate: 47.1,
    },
    {
      name: "Trickster",
      personality: "诡术师",
      description: "出其不意，攻其不备。",
      wins: 6,
      losses: 3,
      winRate: 66.7,
    },
    {
      name: "Defender",
      personality: "守护者",
      description: "坚如磐石，不动如山。",
      wins: 10,
      losses: 7,
      winRate: 58.8,
    },
  ];

  return (
    <div className="space-y-6">
      <h1 className="text-3xl font-bold">Agent 列表</h1>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {agents.map((agent) => (
          <Card key={agent.name}>
            <CardContent className="pt-6">
              <div className="flex justify-between items-start">
                <div>
                  <h3 className="font-semibold text-lg">{agent.name}</h3>
                  <p className="text-sm text-muted-foreground">
                    {agent.personality}
                  </p>
                  <p className="text-sm mt-2">&quot;{agent.description}&quot;</p>
                </div>
                <div className="text-right">
                  <p className="text-2xl font-bold">{agent.winRate}%</p>
                  <p className="text-xs text-muted-foreground">
                    {agent.wins}W {agent.losses}L
                  </p>
                </div>
              </div>
              {/* Win Rate Bar */}
              <div className="mt-4">
                <div className="h-2 bg-secondary rounded-full overflow-hidden">
                  <div
                    className="h-full bg-primary rounded-full"
                    style={{ width: `${agent.winRate}%` }}
                  />
                </div>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}
