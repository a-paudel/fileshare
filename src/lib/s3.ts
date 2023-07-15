import * as minio from "minio";

let endpoint = process.env.S3_ENDPOINT ?? "";
let accessKey = process.env.S3_ACCESS_KEY ?? "";
let secretKey = process.env.S3_SECRET_KEY ?? "";
let bucket = process.env.S3_BUCKET ?? "";
let useSSL = endpoint?.includes("https://") ? true : false;
let parsedEndpoint = endpoint.replace("https://", "").replace("http://", "");

let client = new minio.Client({
  endPoint: parsedEndpoint,
  accessKey,
  secretKey,
  useSSL,
});

export const s3 = {
  client,
  endpoint,
  accessKey,
  secretKey,
  bucket,
};
