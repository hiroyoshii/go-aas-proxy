import http from 'k6/http'
import { check, sleep } from 'k6'
import { thresholds, outputReport } from './common.js';

export const options = {
  duration: '10s',
  vus: 50,
  thresholds: thresholds,
};

export default function () {
  let res = http.get('http://localhost:8081/shells')

  check(res, { 'success': (r) => r.status === 200 })

  sleep(0.3)
}

// レポート出力設定
export function handleSummary(data) {
  console.log('Preparing the end-of-test summary...');

  return outputReport(data)
}