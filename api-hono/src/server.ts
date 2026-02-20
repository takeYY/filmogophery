import { app } from "./core/app";

import "@/app/health";

export default {
  port: process.env.SERVER_PORT,
  fetch: app.fetch,
};
