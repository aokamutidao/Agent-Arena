"use client";
import { useQuery } from "@tanstack/react-query";
import { request, type CourseConfig } from "@/lib/course";
export function useCourseResource<T>(path: string, enabled = true) {
    return useQuery({ queryKey: ["course", path], queryFn: () => request<T>(path), enabled, refetchInterval: 10000, retry: 1 });
}
export function useCourseConfig() { return useCourseResource<CourseConfig>("/api/config"); }
export function useAuthenticatedResource<T>(path: string, token: string | null) {
    return useQuery({ queryKey: ["course", path, token], queryFn: () => request<T>(path, { headers: { Authorization: `Bearer ${token}` } }), enabled: !!token, refetchInterval: 10000, retry: 1 });
}
