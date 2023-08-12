import http from 'k6/http'
import { check, sleep } from 'k6'
import {thresholds, outputReport} from "./common.js"

export const options = {
  scenarios: {
    my_api_test_1: {
      executor: 'constant-arrival-rate',
      rate: 90,
      timeUnit: '1m', // 90 iterations per minute, i.e. 1.5 RPS
      duration: '5m',
      preAllocatedVUs: 10, // the size of the VU (i.e. worker) pool for this scenario
      tags: { test_type: 'api' }, // different extra metric tags for this scenario
      env: { MY_CROC_ID: '1' }, // and we can specify extra environment variables as well!
      exec: 'apitest', // this scenario is executing different code than the one above!
    },
    my_api_test_2: {
      executor: 'ramping-arrival-rate',
      startTime: '30s', // the ramping API test starts a little later
      startRate: 50,
      timeUnit: '1s', // we start at 50 iterations per second
      stages: [
        { target: 200, duration: '30s' }, // go from 50 to 200 iters/s in the first 30 seconds
        { target: 200, duration: '3m30s' }, // hold at 200 iters/s for 3.5 minutes
        { target: 0, duration: '30s' }, // ramp down back to 0 iters/s over the last 30 second
      ],
      preAllocatedVUs: 50, // how large the initial pool of VUs would be
      maxVUs: 100, // if the preAllocatedVUs are not enough, we can initialize more
      tags: { test_type: 'api' }, // different extra metric tags for this scenario
      env: { MY_CROC_ID: '2' }, // same function, different environment variables
      exec: 'apitest', // same function as the scenario above, but with different env vars
    },
  },
  discardResponseBodies: true,
  thresholds: thresholds,
};

export function apitest() {
  let res = http.get('http://localhost:8081/shells')

  check(res, { 'success': (r) => r.status === 200 })

  sleep(0.3)
}


// レポート出力設定
export function handleSummary(data) {
  console.log('Preparing the end-of-test summary...');

  return outputReport(data)
}