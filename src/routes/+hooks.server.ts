import type { Handle } from "@sveltejs/kit";
import { DeleteExpiredFiles } from "../tasks/delete_expired";

export const handle: Handle = async ({ event, resolve }) => {
  await DeleteExpiredFiles();

  let resp = await resolve(event);
  return resp;
};
