type FileRenameMode = 'decrypt' | 'encrypt';

export function downloadFile(mode: FileRenameMode, blob: Blob, filename: string) {
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = getFilename(mode, filename);
  a.click();
  URL.revokeObjectURL(url);
  a.remove();
}

function getFilename(mode: FileRenameMode, filename: string) {
  if (mode === 'decrypt') {
    return filename.replace('.enc', '');
  }

  return `${filename}.enc`;
}
