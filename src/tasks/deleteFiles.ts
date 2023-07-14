// delete files older than 1 hour

import { db } from "../database";
import path from "node:path";
import fs from "node:fs";

async function deleteExpiredFiles() {
  console.log("deleting expired files");
  let files = await db.file.findMany({
    where: {
      createdAt: {
        lt: new Date(Date.now() - 1000 * 60 * 60),
      },
    },
  });

  //   for each file delete the file and the database entry
  for (const file of files) {
    // delete the file
    let filePath = path.resolve("./data/uploads", file.code, file.name);
    // delete directory
    let dirPath = path.dirname(filePath);
    if (fs.existsSync(dirPath)) {
      fs.rmSync(dirPath, { recursive: true });
    }
    // delete database entry
    await db.file.delete({
      where: {
        id: file.id,
      },
    });
  }
}

export function TaskDeleteExpiredFiles() {
  // repeat every minute
  deleteExpiredFiles();
  setInterval(deleteExpiredFiles, 1000 * 60);
}
