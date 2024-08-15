export async function GET() {
  const url = `http://127.0.0.1:8000/movies`;
  console.log("app apiから情報を取得中...");

  const res = await fetch(url);
  return res;
}
