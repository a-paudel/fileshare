import { nanoid } from "nanoid";
import { db } from "../database";
// import * as fs from "fs-extra";
import type { RequestHandler } from "./$types";
import z from "zod";
import { s3 } from "../lib/s3";

async function getCode() {
  let code = nanoid(3);
  // check if already in database
  // if not, return code
  let file = await db.file.findFirst({
    where: {
      code: code,
    },
  });
  if (!file) return code;
  return getCode();
}

let requestSchema = z.object({
  filename: z.string(),
  filesize: z.number(),
});

let responseSchema = z.object({
  code: z.string().max(3),
  uploadurl: z.string(),
});
type responseSchema = z.infer<typeof responseSchema>;

export const POST: RequestHandler = async ({ request }) => {
  let data = requestSchema.safeParse(await request.json());
  if (!data.success) {
    return new Response(JSON.stringify(data.error.formErrors), { status: 400 });
  }
  let { filename, filesize } = data.data;
  if (filesize > 1024 * 1024 * 100) {
    // 100mb
    return new Response(JSON.stringify("file is too large"), { status: 403 });
  }

  let code = await getCode();
  // create entry in db
  let file = await db.file.create({
    data: {
      code,
      name: filename,
      url: `${s3.endpoint}/${s3.bucket}/${code}/${filename}`,
    },
  });

  // get upload url
  let uploadurl = await s3.client.presignedPutObject(
    s3.bucket,
    `${code}/${filename}`,
    60 * 60 * 1
  );

  let response: responseSchema = {
    code,
    uploadurl,
  };

  return new Response(JSON.stringify(response), { status: 200 });
};
