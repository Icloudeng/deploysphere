export function api_fetch(
  input: RequestInfo | URL,
  init?: RequestInit | undefined
) {
  return fetch(input, {
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
    },
    ...(init || {}),
  });
}
