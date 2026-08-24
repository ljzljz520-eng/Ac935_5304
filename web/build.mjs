import { mkdir, writeFile } from 'node:fs/promises';

const dashboard = `<!doctype html><html lang="zh"><head><meta charset="utf-8"><title>ClinicDesk</title></head><body><main><h1>科室制度盘</h1><p>制度、培训记录和排班文件统一入口。</p></main></body></html>`;
await mkdir('dist', { recursive: true });
await writeFile('dist/index.html', dashboard);
console.log('built clinicdesk console');
