import { APIBaseURL } from "@/constants/api";

export async function POST(request: Request) {
  const body = await request.json();

  const res = await fetch(`${APIBaseURL}/users`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });

  return res;
}
