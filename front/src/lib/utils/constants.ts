export const AUTH_COOKIE_NAME = "auth-basic-encoded";

// Job status
export const JOB_STATUS_IDLE = "idle",
  JOB_STATUS_COMPLETED = "completed",
  JOB_STATUS_FAILED = "failed",
  JOB_STATUS_RUNNING = "running";

export const INSTALLER_URL = import.meta.env.VITE_INSTALLER_URL as string;

if (!INSTALLER_URL) {
  throw new Error("VITE_INSTALLER_URL Variable is required");
}
