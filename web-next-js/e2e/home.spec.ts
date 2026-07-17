import { expect, test } from "@playwright/test";

test.describe("Home page", () => {
  test("visual snapshot", async ({ page }) => {
    await page.goto("/");
    await expect(page).toHaveScreenshot("home.png");
  });
});
