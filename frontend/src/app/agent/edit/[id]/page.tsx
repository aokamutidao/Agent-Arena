"use client";

import { useState, useEffect } from "react";
import { useRouter, useParams } from "next/navigation";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useAuth } from "@/lib/auth";

export default function EditAgentPage() {
  const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
  const router = useRouter();
  const params = useParams();
  const agentId = params.id as string;
  const { token } = useAuth();
  const [loading, setLoading] = useState(false);
  const [fetching, setFetching] = useState(true);
  const [error, setError] = useState("");

  const [formData, setFormData] = useState({
    name: "",
    personality: "",
    api_endpoint: "",
    api_key: "", // 服务端会返回（owner 自己查看）
    model: "",
    challenge_fee: "10",
  });

  const [testing, setTesting] = useState(false);
  const [testResult, setTestResult] = useState<{success: boolean, message: string, response?: string} | null>(null);

  useEffect(() => {
    if (!token || !agentId) return;

    const fetchAgent = async () => {
      try {
        const res = await fetch(`${API_URL}/api/marketplace/agents/${agentId}`, {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (!res.ok) {
          throw new Error("Agent not found");
        }

        const agent = await res.json();
        setFormData({
          name: agent.name,
          personality: agent.personality,
          api_endpoint: agent.api_endpoint || "",
          api_key: agent.api_key || "",
          model: agent.model || "",
          challenge_fee: agent.challenge_fee?.toString() || "10",
        });
      } catch (err: any) {
        setError(err.message);
      } finally {
        setFetching(false);
      }
    };

    fetchAgent();
  }, [token, agentId]);

  const handleTestAPI = async () => {
    if (!formData.api_endpoint || !formData.api_key) {
      setTestResult({
        success: false,
        message: "请先填写 API Endpoint 和 API Key",
      });
      return;
    }

    setTesting(true);
    setTestResult(null);

    try {
      const res = await fetch(`${API_URL}/api/auth/agents/test-api`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          api_endpoint: formData.api_endpoint,
          api_key: formData.api_key,
          model: formData.model || "gpt-3.5-turbo",
        }),
      });

      const data = await res.json();
      setTestResult({
        success: data.success,
        message: data.success ? data.message : data.error,
        response: data.response,
      });
    } catch (err: any) {
      setTestResult({
        success: false,
        message: "测试请求失败: " + err.message,
      });
    } finally {
      setTesting(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!token) {
      setError("请先登录");
      return;
    }

    setLoading(true);
    setError("");

    try {
      const res = await fetch(`${API_URL}/api/auth/agents/${agentId}`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          name: formData.name,
          personality: formData.personality,
          api_endpoint: formData.api_endpoint,
          api_key: formData.api_key,
          model: formData.model,
          challenge_fee: parseInt(formData.challenge_fee),
        }),
      });

      if (!res.ok) {
        const err = await res.json();
        throw new Error(err.error || "更新失败");
      }

      router.push(`/profile`);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  if (fetching) {
    return <div className="max-w-2xl mx-auto">加载中...</div>;
  }

  return (
    <div className="max-w-2xl mx-auto space-y-6">
      <h1 className="text-3xl font-bold">编辑 Agent</h1>

      <Card>
        <CardHeader>
          <CardTitle>Agent 信息</CardTitle>
          <CardDescription>
            修改你的 AI Agent 的性格和配置
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-6">
            {/* 名称 */}
            <div className="space-y-2">
              <label className="text-sm font-medium">名称 *</label>
              <Input
                placeholder="例如：Shadow Assassin"
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                maxLength={50}
                required
              />
              <p className="text-xs text-muted-foreground">
                {formData.name.length}/50 字符
              </p>
            </div>

            {/* 性格描述 */}
            <div className="space-y-2">
              <label className="text-sm font-medium">性格描述 *</label>
              <textarea
                className="w-full min-h-[150px] rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                placeholder="描述你的 Agent 的战斗风格、策略偏好、性格特点..."
                value={formData.personality}
                onChange={(e) => setFormData({ ...formData, personality: e.target.value })}
                maxLength={500}
                required
              />
              <p className="text-xs text-muted-foreground">
                {formData.personality.length}/500 字符 — 这会影响 Agent 的 AI 决策
              </p>
            </div>

            {/* API Endpoint */}
            <div className="space-y-2">
              <label className="text-sm font-medium">
                AI API Endpoint *
              </label>
              <Input
                placeholder="https://api.openai.com/v1/chat/completions"
                value={formData.api_endpoint}
                onChange={(e) => setFormData({ ...formData, api_endpoint: e.target.value })}
                required
              />
              <p className="text-xs text-muted-foreground">
                你的 AI 服务 endpoint（兼容 OpenAI 或 Anthropic 格式）。
                例如：OpenAI、Azure OpenAI、Qwen DashScope 等。
              </p>
            </div>

            {/* API Key */}
            <div className="space-y-2">
              <label className="text-sm font-medium">
                API Key *
              </label>
              <Input
                type="password"
                placeholder="sk-..."
                value={formData.api_key}
                onChange={(e) => setFormData({ ...formData, api_key: e.target.value })}
                required
              />
              <p className="text-xs text-muted-foreground">
                你的 AI 服务 API key。服务端使用，不会泄露给其他用户。
              </p>
            </div>

            {/* Model */}
            <div className="space-y-2">
              <label className="text-sm font-medium">
                模型名称 *
              </label>
              <Input
                placeholder="例如：gpt-3.5-turbo, Qwen/Qwen2.5-7B-Instruct"
                value={formData.model}
                onChange={(e) => setFormData({ ...formData, model: e.target.value })}
                required
              />
              <p className="text-xs text-muted-foreground">
                你的 AI 服务支持的模型名称。例如：OpenAI 的 <code className="bg-muted px-1 rounded">gpt-3.5-turbo</code>，SiliconFlow 的 <code className="bg-muted px-1 rounded">Qwen/Qwen2.5-7B-Instruct</code>。
              </p>
            </div>

            {/* Test API Button */}
            <div className="space-y-2">
              <Button
                type="button"
                variant="outline"
                onClick={handleTestAPI}
                disabled={testing}
              >
                {testing ? "测试中..." : "🔍 测试 API 连通性"}
              </Button>
              {testResult && (
                <div className={`p-3 rounded-lg text-sm ${
                  testResult.success
                    ? "bg-green-50 border border-green-200 text-green-800"
                    : "bg-red-50 border border-red-200 text-red-800"
                }`}>
                  <div className="font-medium">
                    {testResult.success ? "✅ 成功" : "❌ 失败"}
                  </div>
                  <div className="mt-1">{testResult.message}</div>
                  {testResult.response && (
                    <div className="mt-2 p-2 bg-white rounded border text-xs font-mono">
                      <div className="text-muted-foreground mb-1">AI 回复:</div>
                      {testResult.response}
                    </div>
                  )}
                </div>
              )}
            </div>

            {/* 挑战费用 */}
            <div className="space-y-2">
              <label className="text-sm font-medium">挑战费用</label>
              <Input
                type="number"
                min="0"
                value={formData.challenge_fee}
                onChange={(e) => setFormData({ ...formData, challenge_fee: e.target.value })}
              />
            </div>

            {error && (
              <p className="text-sm text-red-500">{error}</p>
            )}

            <div className="flex gap-3">
              <Button type="submit" disabled={loading}>
                {loading ? "保存中..." : "保存修改"}
              </Button>
              <Button
                type="button"
                variant="outline"
                onClick={() => router.back()}
              >
                取消
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
