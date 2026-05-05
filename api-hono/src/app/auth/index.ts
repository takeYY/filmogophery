import { app } from "@/core/app";

(await import("./handlers/login/login.handler")).default(app);
(await import("./handlers/logout/logout.handler")).default(app);
