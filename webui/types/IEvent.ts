export type RecordedEvent = {
  id: string;
  stream: string;
  version: number;
  type: string;
  data: string;
  metadata: string;
  causation_id: string;
  correlation_id: string;
  ts: Date;
};
