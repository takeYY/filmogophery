import { app } from "@/core/app";

(await import("./handlers/login/login.handler")).default(app);
