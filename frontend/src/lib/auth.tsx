"use client";

import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useRef,
  useState,
  type ReactNode,
} from "react";
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

function sessionIdentity(session: Session) {
  return `${session.user.address.toLowerCase()}:${session.chain_id}`;
}

function isSession(value: unknown): value is Session {
  if (!value || typeof value !== "object") return false;
  const saved = value as Partial<Session>;
  return (
    typeof saved.token === "string" &&
    saved.token.length > 0 &&
    typeof saved.chain_id === "number" &&
    Number.isSafeInteger(saved.chain_id) &&
    saved.chain_id > 0 &&
    !!saved.user &&
    typeof saved.user.address === "string" &&
    /^0x[0-9a-f]{40}$/i.test(saved.user.address) &&
    Array.isArray(saved.user.roles) &&
    saved.user.roles.every((role) => typeof role === "string")
  );
}

function forgetPersistedSession() {
  try {
    localStorage.removeItem("course_session");
  } catch {
    // Browsers may block storage; in-memory authentication still works.
  }
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const { address, chainId, status } = useAccount();
  const [session, setSession] = useState<Session | null>(null);
  const identity = useRef<string | null>(null);
  const restored = useRef(false);
  const currentIdentity =
    address && chainId ? `${address.toLowerCase()}:${chainId}` : null;
  const liveIdentity = useRef(currentIdentity);
  const liveSession = useRef<Session | null>(null);
  const generation = useRef(0);

  if (liveIdentity.current !== currentIdentity) {
    liveIdentity.current = currentIdentity;
    generation.current += 1;
  }

  const updateSession = useCallback((next: Session | null) => {
    liveSession.current = next;
    setSession(next);
  }, []);

  const clearSession = useCallback(() => {
    generation.current += 1;
    updateSession(null);
    forgetPersistedSession();
  }, [updateSession]);

  useEffect(() => {
    // WagmiProvider starts in reconnecting, so initial hydration cannot be mistaken for disconnect.
    if (status === "connecting" || status === "reconnecting") return;
    if (!restored.current) {
      restored.current = true;
      try {
        const raw = localStorage.getItem("course_session");
        const saved: unknown = raw ? JSON.parse(raw) : null;
        if (
          status === "connected" &&
          isSession(saved) &&
          sessionIdentity(saved) === currentIdentity
        ) {
          updateSession(saved);
        } else {
          clearSession();
        }
      } catch {
        clearSession();
      }
    } else if (identity.current !== currentIdentity || status === "disconnected") {
      clearSession();
    }
    identity.current = currentIdentity;
  }, [currentIdentity, status, clearSession, updateSession]);

  // Hide a previous identity immediately, before the cleanup effect runs.
  const visibleSession =
    status === "connected" &&
    session &&
    sessionIdentity(session) === currentIdentity
      ? session
      : null;
  const savedToken = visibleSession?.token;

  useEffect(() => {
    if (!savedToken || !currentIdentity) return;
    let active = true;
    const refresh = async () => {
      const requestGeneration = generation.current;
      const stillCurrent = () =>
        active &&
        requestGeneration === generation.current &&
        liveIdentity.current === currentIdentity &&
        liveSession.current?.token === savedToken;
      try {
        const profile = await request<User>("/api/auth/profile", {
          headers: { Authorization: `Bearer ${savedToken}` },
        });
        if (!stillCurrent()) return;
        const current = liveSession.current;
        if (
          !current ||
          profile.address.toLowerCase() !== current.user.address.toLowerCase()
        ) {
          clearSession();
          return;
        }
        updateSession({ ...current, user: profile });
      } catch (error) {
        if (stillCurrent() && error instanceof ApiError && error.status === 401) {
          clearSession();
        }
      }
    };
    void refresh();
    const timer = setInterval(refresh, 30000);
    return () => {
      active = false;
      clearInterval(timer);
    };
  }, [savedToken, currentIdentity, clearSession, updateSession]);

  const logout = () => clearSession();

  const login = async (
    walletAddress: string,
    nonce: string,
    signature: string,
  ) => {
    const requestedIdentity = `${walletAddress.toLowerCase()}:${chainId}`;
    if (
      status !== "connected" ||
      !chainId ||
      requestedIdentity !== liveIdentity.current
    ) {
      throw new Error("请先连接钱包。");
    }
    const requestGeneration = ++generation.current;
    const response = await request<{ token: string; user: User }>(
      "/api/auth/login",
      {
        method: "POST",
        body: JSON.stringify({
          address: walletAddress.toLowerCase(),
          nonce,
          signature,
        }),
      },
    );
    if (
      liveIdentity.current !== requestedIdentity ||
      requestGeneration !== generation.current
    ) {
      throw new Error("签名期间钱包或网络已变化，请重新登录。");
    }
    const next = { ...response, chain_id: chainId };
    if (!isSession(next) || sessionIdentity(next) !== requestedIdentity) {
      throw new Error("签名期间钱包或网络已变化，请重新登录。");
    }
    updateSession(next);
    try {
      localStorage.setItem("course_session", JSON.stringify(next));
    } catch {
      // The session remains usable until reload when persistent storage is unavailable.
    }
  };

  return (
    <AuthContext.Provider
      value={{
        user: visibleSession?.user || null,
        token: visibleSession?.token || null,
        login,
        logout,
        isAuthenticated: !!visibleSession,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) throw new Error("AuthProvider is missing");
  return context;
}
