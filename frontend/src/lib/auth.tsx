"use client";
import { createContext, useContext, useEffect, useRef, useState, type ReactNode } from "react";
import { useAccount } from "wagmi";
import { ApiError, request, type Role } from "@/lib/course";
interface User {
    address: string;
    roles: Role[];
}
interface Session {
    token: string;
    user: User;
    chain_id: number;
}
interface AuthContextType {
    user: User | null;
    token: string | null;
    login: (address: string, nonce: string, signature: string) => Promise<void>;
    logout: () => void;
    isAuthenticated: boolean;
}
const AuthContext = createContext<AuthContextType | undefined>(undefined);
export function AuthProvider({ children }: {
    children: ReactNode;
}) {
    const { address, chainId, status } = useAccount();
    const [session, setSession] = useState<Session | null>(null);
    const identity = useRef<string | null>(null);
    const restored = useRef(false);
    const currentIdentity = address && chainId ? `${address.toLowerCase()}:${chainId}` : null;
    const liveIdentity = useRef(currentIdentity);
    liveIdentity.current = currentIdentity;
    useEffect(() => {
        if (status === "connecting" || status === "reconnecting")
            return;
        if (!restored.current) {
            restored.current = true;
            try {
                const raw = localStorage.getItem("course_session");
                const saved = raw ? JSON.parse(raw) as Session : null;
                if (saved && currentIdentity === `${saved.user.address.toLowerCase()}:${saved.chain_id}` && Array.isArray(saved.user.roles))
                    setSession(saved);
                else
                    localStorage.removeItem("course_session");
            }
            catch {
                localStorage.removeItem("course_session");
            }
        }
        else if (identity.current !== currentIdentity) {
            setSession(null);
            localStorage.removeItem("course_session");
        }
        identity.current = currentIdentity;
    }, [currentIdentity, status]);
    const savedToken = session?.token;
    useEffect(() => {
        if (!savedToken || !currentIdentity)
            return;
        let active = true;
        const refresh = async () => {
            try {
                const profile = await request<User>("/api/auth/profile", { headers: { Authorization: `Bearer ${savedToken}` } });
                if (active)
                    setSession(previous => previous?.token === savedToken ? { ...previous, user: profile } : previous);
            }
            catch (error) {
                if (active && error instanceof ApiError && error.status === 401) {
                    setSession(null);
                    localStorage.removeItem("course_session");
                }
            }
        };
        void refresh();
        const timer = setInterval(refresh, 30000);
        return () => { active = false; clearInterval(timer); };
    }, [savedToken, currentIdentity]);
    const logout = () => { setSession(null); localStorage.removeItem("course_session"); };
    const login = async (walletAddress: string, nonce: string, signature: string) => {
        if (!chainId)
            throw new Error("请先连接钱包。");
        const response = await request<{
            token: string;
            user: User;
        }>("/api/auth/login", { method: "POST", body: JSON.stringify({ address: walletAddress.toLowerCase(), nonce, signature }) });
        if (liveIdentity.current !== `${walletAddress.toLowerCase()}:${chainId}`)
            throw new Error("签名期间钱包或网络已变化，请重新登录。");
        const next = { ...response, chain_id: chainId };
        setSession(next);
        localStorage.setItem("course_session", JSON.stringify(next));
    };
    return <AuthContext.Provider value={{ user: session?.user || null, token: session?.token || null, login, logout, isAuthenticated: !!session }}>{children}</AuthContext.Provider>;
}
export function useAuth() { const context = useContext(AuthContext); if (!context)
    throw new Error("AuthProvider is missing"); return context; }
