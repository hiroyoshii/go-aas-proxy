import http from 'k6/http'
import { check, sleep } from 'k6'
import {htmlReport} from "https://raw.githubusercontent.com/benc-uk/k6-reporter/main/dist/bundle.js";
import {jUnit, textSummary} from "https://jslib.k6.io/k6-summary/0.0.1/index.js";

export const options = {
  duration: '10s',
  vus: 50,
  thresholds: {
    http_req_failed: ['rate<0.01'], // http errors should be less than 1%
    http_req_duration: ['p(95)<500'], // 95 percent of response times must be below 500ms
  },
};

export default function () {
  let res = http.get('http://localhost:8081/shells')

  check(res, { 'success login': (r) => r.status === 200 })

  sleep(0.3)
}

// レポート出力設定
export function handleSummary(data) {
  console.log('Preparing the end-of-test summary...');

  return {
      stdout: textSummary(data, { indent: " ", enableColors: true }),
      "output/summary.html": htmlReport(data),
      'output/summary.xml': jUnit(data), // but also transform it and save it as a JUnit XML...
      'output/summary.json': JSON.stringify(data), // and a JSON with all the details...
  }
}