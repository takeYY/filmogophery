export async function GET() {
  const url = `http://127.0.0.1:8000/media`;

  console.log("app api から情報を取得しました。");

  const res = await fetch(url);
  return res;
}
