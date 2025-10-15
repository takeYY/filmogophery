import { APIBaseURL } from "@/constants/api";
import { NextResponse } from "next/server";

export async function POST(
  req: Request,
  { params }: { params: { id: string } }
) {
  console.log("POST route.ts called!");

  const url = `${APIBaseURL}/movies/${params.id}/records`;
  console.log(`リクエスト先URL: ${url}`);

  try {
    const requestData = await req.json();
    console.log("リクエストデータ:", requestData);

    const res = await fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Access-Control-Allow-Origin": "*",
      },
      body: JSON.stringify(requestData),
    });
    console.log(`APIレスポンスステータス: ${res.status}`);

    return res;
  } catch (error) {
    console.error("Error in API route:", error);
    return NextResponse.json(
      { error: "Internal Server Error" },
      { status: 500 }
    );
  }
}
