export async function api(
    input: RequestInfo | URL,
    init?: RequestInit
): Promise<Response> {
    const response = await fetch(input, {
        credentials: "include",
        ...init,
    });

    if (response.status === 401) {
        const redirect = encodeURIComponent(window.location.pathname);
        window.location.replace(`/auth?redirect=${redirect}`);
        throw new Error("unauthorized");
    }

    return response;
}

export async function apiJson<T>(
    input: RequestInfo | URL,
    init?: RequestInit
): Promise<T> {
    const response = await api(input, init);
    return response.json();
}
