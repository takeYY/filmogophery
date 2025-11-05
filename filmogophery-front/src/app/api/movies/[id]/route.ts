import { NextResponse } from "next/server";

export async function GET(
  req: Request,
  { params }: { params: { id: string } }
) {
  // 映画詳細取得のロジック
  return NextResponse.json({ message: "Movie details" });
}
