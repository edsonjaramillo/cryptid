import * as v from 'valibot';

export const universalSchema = v.object({
  file: v.pipe(v.instance(FileList), v.minLength(1, 'File is required')),
  password: v.pipe(v.string(), v.minLength(1, 'Password is required')),
});

export type UniversalFormData = v.InferOutput<typeof universalSchema>;
