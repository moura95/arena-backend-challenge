import http from 'k6/http';
import { check, sleep } from 'k6';

// Test configuration
export const options = {
    stages: [
        { duration: '10s', target: 50 },    // Ramp up to 50 users
        { duration: '20s', target: 100 },   // Ramp up to 100 users
        { duration: '20s', target: 250 },   // Ramp up to 250 users
        { duration: '30s', target: 500 },   // Ramp up to 500 users
        { duration: '30s', target: 500 },   // Stay at 500 users
        { duration: '20s', target: 1000 },  // Spike to 1000 users
        { duration: '10s', target: 0 },     // Ramp down to 0 users
    ],
    thresholds: {
        http_req_duration: ['p(95)<100'], // 95% of requests must complete below 100ms
        http_req_failed: ['rate<0.50'],   // Error rate must be less than 50% (accounting for 404s)
    },
};

const BASE_URL = 'http://localhost:8080';

// Test data - mix of valid and invalid IPs
const testIPs = [
    // IPs that DEFINITELY exist in the database (from CSV)
    { ip: '1.0.0.100', expectedStatus: [200], description: 'US - Los Angeles (in DB)' },
    { ip: '1.0.1.50', expectedStatus: [200], description: 'China - Fuzhou (in DB)' },
    { ip: '1.0.16.200', expectedStatus: [200], description: 'Australia - Melbourne (in DB)' },
    { ip: '1.0.64.1', expectedStatus: [200], description: 'Japan - Tokyo (in DB)' },
    { ip: '1.0.128.100', expectedStatus: [200], description: 'Thailand (in DB)' },
    { ip: '1.1.0.1', expectedStatus: [200], description: 'China - Guangzhou (in DB)' },

    // Popular IPs that might or might not be in DB
    { ip: '8.8.8.8', expectedStatus: [200, 404], description: 'Google DNS' },
    { ip: '1.1.1.1', expectedStatus: [200, 404], description: 'Cloudflare DNS' },

    // Private IPs (likely 404)
    { ip: '192.168.1.1', expectedStatus: [404], description: 'Private IP' },
    { ip: '10.0.0.1', expectedStatus: [404], description: 'Private IP' },

    // Invalid IPs that should return 400
    { ip: 'invalid.ip', expectedStatus: [400], description: 'Invalid format' },
    { ip: '256.256.256.256', expectedStatus: [400], description: 'Out of range' },
    { ip: '1.2.3', expectedStatus: [400], description: 'Missing octet' },
    { ip: '', expectedStatus: [400], description: 'Empty IP' },
];

export default function () {
    // Select a random IP from test data
    const testCase = testIPs[Math.floor(Math.random() * testIPs.length)];

    // Make request
    const url = testCase.ip
        ? `${BASE_URL}/ip/location?ip=${testCase.ip}`
        : `${BASE_URL}/ip/location`;

    const response = http.get(url);

    // Validations
    const isValidStatus = Array.isArray(testCase.expectedStatus)
        ? testCase.expectedStatus.includes(response.status)
        : response.status === testCase.expectedStatus;

    check(response, {
        'status is correct': () => isValidStatus,
        'response time < 100ms': (r) => r.timings.duration < 100,
        'content-type is JSON': (r) => r.headers['Content-Type'] === 'application/json',
    });

    // Additional checks for successful responses (200)
    if (response.status === 200) {
        check(response, {
            'has country field': (r) => {
                try {
                    const body = JSON.parse(r.body);
                    return body.country !== undefined;
                } catch (e) {
                    return false;
                }
            },
            'has countryCode field': (r) => {
                try {
                    const body = JSON.parse(r.body);
                    return body.countryCode !== undefined;
                } catch (e) {
                    return false;
                }
            },
        });
    }

    // Additional checks for error responses (400 or 404)
    if (response.status >= 400) {
        check(response, {
            'has error field': (r) => {
                try {
                    const body = JSON.parse(r.body);
                    return body.error !== undefined;
                } catch (e) {
                    return false;
                }
            },
        });
    }

    // Small sleep between requests (100-300ms)
    sleep(Math.random() * 0.2 + 0.1);
}

// Smoke test - runs before the main test
export function setup() {
    console.log('Running smoke test...');

    const response = http.get(`${BASE_URL}/ip/location?ip=8.8.8.8`);

    if (response.status !== 200) {
        throw new Error(`Smoke test failed: Server returned status ${response.status}`);
    }

    console.log('Smoke test passed! Server is ready.');
}