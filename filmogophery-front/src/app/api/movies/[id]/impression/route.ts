import { APIBaseURL } from "@/constants/api";
import { NextResponse } from "next/server";

// Deprecated
export async function PUT(
  req: Request,
  { params }: { params: { id: string } }
) {
  const url = `${APIBaseURL}/movies/${params.id}/impression`;

  try {
    const requestData = await req.json();
    const processedData = {
      ...requestData,
      rating: parseFloat(requestData.rating), // float にする
    };
    console.log("リクエストデータ:", processedData);

    const res = await fetch(url, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        "Access-Control-Allow-Origin": "*",
      },
      body: JSON.stringify(processedData),
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
