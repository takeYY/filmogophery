import { defineConfig } from "drizzle-kit";

export default defineConfig({
  dialect: "mysql",
  schema: "./src/core/schemas/schema.ts",
  out: "./src/core/drizzle",
  dbCredentials: {
    host: "localhost",
    port: 3306,
    database: "db4dev",
    user: "user",
    password: "password",
  },
});
