declare module '@babel/standalone' {
  export interface TransformOptions {
    presets?: string[];
    filename?: string;
  }

  export function transform(
    code: string,
    options?: TransformOptions,
  ): { code?: string };
}
