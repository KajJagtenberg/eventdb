export type IEvent = {
  id: string;
  stream: string;
  version: number;
  type: string;
  data: Record<string, any>;
  metadata: Record<string, any>;
  causation_id: string;
  correlation_id: string;
  ts: Date;
};
