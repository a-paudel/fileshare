import { db } from "../database";
import { s3 } from "../lib/s3";

export async function DeleteExpiredFiles() {
  let files = await db.file.findMany({
    where: {
      createdAt: {
        lt: new Date(Date.now() - 60 * 60 * 1000),
      },
    },
  });
  let ids = files.map((f) => f.id);

  for (const file of files) {
    // delete from s3
    await s3.client.removeObject(s3.bucket, `${file.code}/${file.name}`);
  }
  // delete from the database
  await db.file.deleteMany({
    where: {
      id: {
        in: ids,
      },
    },
  });
}
