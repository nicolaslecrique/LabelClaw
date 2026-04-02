export type JsonValue =
  | null
  | boolean
  | number
  | string
  | JsonValue[]
  | { [key: string]: JsonValue };

export interface GenerateConfigurationRequest {
  sampleSchema: JsonValue;
  labelSchema: JsonValue;
  uiPrompt: string;
}

export interface GenerateConfigurationResponse {
  componentSource: string;
  sampleData: JsonValue;
}

export interface SavedConfiguration extends GenerateConfigurationRequest {
  sampleData: JsonValue;
  componentSource: string;
  updatedAt: string;
}

