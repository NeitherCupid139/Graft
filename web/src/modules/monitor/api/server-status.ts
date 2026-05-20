import { request } from '@/utils/request';

import { MONITOR_API_PATH } from '../contract/paths';
import type { ServerStatusResponse } from '../types/server-status';

export function getServerStatus() {
  return request.get<ServerStatusResponse>({
    url: MONITOR_API_PATH.SERVER_STATUS,
  });
}
