import { error } from "@sveltejs/kit";
import { db } from "../../database";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ params }) => {
  let file = await db.file.findFirst({
    where: {
      code: params.code,
    },
  });
  if (!file) throw error(404, "File not found");

  return {
    file,
  };
};
