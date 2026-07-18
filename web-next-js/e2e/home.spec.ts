import { expect, test } from "@playwright/test";

/**
 * 未認証状態でのホームページアクセステスト
 * Middleware が access_token Cookie の有無を確認し /login にリダイレクトする
 */
test.describe("Home page (unauthenticated)", () => {
  test("未認証でアクセスするとログイン画面にリダイレクトされる", async ({
    page,
  }) => {
    // Cookie なしでアクセス → Middleware がリダイレクト
    await page.goto("/");
    await expect(page).toHaveURL(/\/login/);
  });

  test("ログイン画面に正しい要素が表示される", async ({ page }) => {
    await page.goto("/");
    await expect(page).toHaveURL(/\/login/);
    await expect(page.getByRole("heading", { name: "ログイン" })).toBeVisible();
    await expect(page.getByLabel("メールアドレス")).toBeVisible();
    await expect(page.locator("#password")).toBeVisible();
    await expect(page.getByRole("button", { name: "ログイン" })).toBeVisible();
  });

  test("visual snapshot", async ({ page }) => {
    await page.goto("/");
    await expect(page).toHaveURL(/\/login/);
    await expect(page).toHaveScreenshot("home-unauthenticated.png");
  });
});
