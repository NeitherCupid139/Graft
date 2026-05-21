import { request } from '@/utils/request';

import { MONITOR_API_PATH } from '../contract/paths';
import type { MonitorTrendRange } from '../contract/trend';
import type { ServerStatusResponse } from '../types/server-status';

export function getServerStatus(trendRange: MonitorTrendRange) {
  return request.get<ServerStatusResponse>({
    url: MONITOR_API_PATH.SERVER_STATUS,
    params: {
      trend_range: trendRange,
    },
  });
}
