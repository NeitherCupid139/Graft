import { request } from '@/utils/request';

export const REALTIME_API_PATH = {
  SUBSCRIPTIONS: '/api/realtime/subscriptions',
} as const;

export type RealtimeSubscriptionRequest = {
  topic: string;
};

export type RealtimeSubscriptionResponse = {
  topic: string;
  ticket: string;
  websocket_url: string;
  expires_at: string;
};

export function postRealtimeSubscription(body: RealtimeSubscriptionRequest) {
  return request.post<RealtimeSubscriptionResponse>({
    url: REALTIME_API_PATH.SUBSCRIPTIONS,
    data: body,
  }) as Promise<RealtimeSubscriptionResponse>;
}
