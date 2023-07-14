import { redirect, type Actions } from "@sveltejs/kit";
import { nanoid } from "nanoid";
import { db } from "../database";
// import * as fs from "fs-extra";
import fs from "node:fs/promises";
import * as path from "node:path";

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

export const actions: Actions = {
  default: async ({ request }) => {
    let data = await request.formData();
    let uploadFile = data.get("file") as File | null;
    if (!uploadFile) return { error: "No file uploaded" };
    // create code
    let code = await getCode();
    // create file in database
    let file = await db.file.create({
      data: {
        code: code,
        name: uploadFile.name,
      },
    });

    // save file to ./data/uploads/code/filename
    let filePath = path.resolve("./data/uploads", code, uploadFile.name);
    // make parent directory
    // check if dir exists
    // if (!fs.existsSync(path.dirname(filePath))) {
    await fs.mkdir(path.dirname(filePath), { recursive: true });
    // }
    // write file
    // to arraybufferview
    let fileData = await uploadFile.arrayBuffer();
    await fs.writeFile(filePath, Buffer.from(fileData));

    // redirect to code page
    throw redirect(302, `/${code}`);
  },
};
