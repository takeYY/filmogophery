import { app } from "@/core/app";

(await import("./handlers/create-user/create-user.handler")).default(app);
