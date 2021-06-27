import { Application, Context, log } from "@utils/deps.ts";
import * as configs from "@utils/const.ts";

const app = new Application();

app.use((ctx: Context) => {
  ctx.response.body = "Hello world!";
});

await app.listen({ port: configs.PORT });
log.info(`${configs.APP_NAME} running on localhost:${configs.PORT}`);
