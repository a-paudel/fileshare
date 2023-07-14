import { error } from "@sveltejs/kit";
import { db } from "../../../database";
import * as path from "node:path";
import fs from "node:fs";
import { Readable } from "node:stream";
import type { RequestHandler } from "./$types";

export const GET: RequestHandler = async ({ params }) => {
  // get code

  let file = await db.file.findFirst({
    where: {
      code: params.code,
    },
  });
  if (!file) throw error(404, "File not found");

  let filePath = path.resolve("./data/uploads", params.code, file.name);
  //   check if file exists
  if (!fs.existsSync(filePath)) throw error(404, "File not found");

  //   serve file
  let stream = fs.createReadStream(filePath);
  let rStream = Readable.toWeb(stream);
  // @ts-ignore
  return new Response(rStream, {
    headers: { "Content-Disposition": "attachment" },
  });
};
