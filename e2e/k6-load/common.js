import {htmlReport} from "https://raw.githubusercontent.com/benc-uk/k6-reporter/main/dist/bundle.js";
import {jUnit, textSummary} from "https://jslib.k6.io/k6-summary/0.0.1/index.js";

export const thresholds =  {
    http_req_failed: ['rate<0.01'], // http errors should be less than 1%
    http_req_duration: ['p(95)<500'], // 95 percent of response times must be below 500ms
};

// レポート出力設定
export function outputReport(data) {
    return {
    stdout: textSummary(data, { indent: " ", enableColors: true }),
    "output/summary.html": htmlReport(data),
    'output/summary.xml': jUnit(data), // but also transform it and save it as a JUnit XML...
    'output/summary.json': JSON.stringify(data), // and a JSON with all the details...
    }
  }