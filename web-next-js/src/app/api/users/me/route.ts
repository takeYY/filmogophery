import { APIBaseURL } from "@/constants/api";

export async function GET(request: Request) {
  /* const token = request.headers.get("Authorization");

  const res = await fetch(`${APIBaseURL}/users/me`, {
    headers: {
      Authorization: token || "",
    },
  });

  return res; */

  // モックを返す
  const mockUser = {
    id: 1,
    username: "test user",
  };
  return Response.json(mockUser);
}
