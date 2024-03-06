import { INSTALLER_URL } from "$lib/utils/constants";
import { api_fetch } from "$lib/utils/fetch";

export function api_root_query() {
  return api_fetch(`${INSTALLER_URL}/`);
}
