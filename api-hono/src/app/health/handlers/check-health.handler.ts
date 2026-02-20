import { AppType } from "@/core/app";
import { StatusCodes } from "http-status-codes";

export default function (app: AppType) {
  app.get(
    `/health`,

    async (c) => {
      const result = {
        message: "system all green",
      };
      return c.json(result, StatusCodes.OK);
    },
  );
}
