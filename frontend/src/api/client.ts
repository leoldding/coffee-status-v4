export async function api(
    input: RequestInfo | URL,
    init?: RequestInit
): Promise<Response> {
    let response = await fetch(input, {
        credentials: "include",
        ...init,
    });

    if (response.status === 401) {
        // Try refreshing.
        const refresh = await fetch("/api/v1/auth/refresh", {
            method: "POST",
            credentials: "include",
        });

        if (!refresh.ok) {
            const redirect = encodeURIComponent(window.location.pathname);
            window.location.replace(`/auth?redirect=${redirect}`);
            throw new Error("unauthorized");
        }

        // Retry original request.
        response = await fetch(input, {
            credentials: "include",
            ...init,
        });

        if (response.status === 401) {
            const redirect = encodeURIComponent(window.location.pathname);
            window.location.replace(`/auth?redirect=${redirect}`);
            throw new Error("unauthorized");
        }
    }

    if (!response.ok) {
        throw new Error(`request failed with status ${response.status}`);
    }

    return response;
}

export async function apiJson<T>(
    input: RequestInfo | URL,
    init?: RequestInit
): Promise<T> {
    const response = await api(input, init);
    return response.json() as Promise<T>;
}
