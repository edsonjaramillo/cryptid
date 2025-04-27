import { env } from 'node:process';
import * as v from 'valibot';

const envSchema = v.object({
  WEB_PORT: v.optional(
    v.pipe(
      v.string(),
      v.regex(/^\d+$/),
      v.transform(value => Number(value)),
    ),
  ),
});

export const ENV = v.parse(envSchema, {
  WEB_PORT: env.WEB_PORT,
});
