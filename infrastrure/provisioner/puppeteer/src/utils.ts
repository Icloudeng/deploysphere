import * as https from "https";
import * as fs from "fs/promises";
import * as path from "path";
import { createWriteStream, unlink } from "fs";

export const wait = (sec = 1) =>
  new Promise((resolve) => setTimeout(resolve, sec * 1000));

export const log = (text: string) => console.log(`%%%${text}%%%`);

export const downloadFile = async (
  remoteFileUrl: string,
  localFilePath: string
) => {
  try {
    await fs.access(localFilePath, fs.constants.F_OK);
    return path.join(process.cwd(), localFilePath);
  } catch (_) {}

  const file = createWriteStream(localFilePath);

  return new Promise<string>((resolve, reject) => {
    https
      .get(remoteFileUrl, (response) => {
        if (response.statusCode !== 200) {
          console.error("Failed to download the file:", response.statusMessage);
          return;
        }

        response.pipe(file);

        file.on("finish", () => {
          file.close();
          resolve(path.join(process.cwd(), localFilePath));
        });
      })
      .on("error", (err) => {
        reject("Error while downloading: " + err.message);
        console.error("Error while downloading:", err.message);

        unlink(localFilePath, (unlinkErr) => {
          if (unlinkErr)
            console.error("Failed to delete the incomplete file:", unlinkErr);
        });
      });
  });
};
