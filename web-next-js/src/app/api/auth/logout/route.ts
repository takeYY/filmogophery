import { APIBaseURL } from "@/constants/api";

export async function POST(request: Request) {
  const token = request.headers.get("Authorization");

  const res = await fetch(`${APIBaseURL}/auth/logout`, {
    method: "POST",
    headers: {
      Authorization: token || "",
    },
  });

  return res;
}
