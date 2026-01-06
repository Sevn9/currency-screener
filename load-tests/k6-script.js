import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

const BASE_URL = 'http://localhost:8080'; // Порт Gateway

/**
 * Вспомогательная функция для получения даты в формате YYYY-MM-DD
 * @param {number} daysOffset - Смещение относительно текущей даты (например, 0 для сегодня, 1 для завтра, -1 для вчера)
 * @returns {string} Дата в формате YYYY-MM-DD
 */
function getFormattedDate(daysOffset) {
    const date = new Date();
    date.setDate(date.getDate() + daysOffset);
    return date.toISOString().split('T')[0];
}

export let options = {
    // Сценарий нагрузки
    stages: [
        { duration: '10s', target: 5 },  // Разгон до 5 пользователей
        { duration: '20s', target: 5 },  // Держать нагрузку
        { duration: '10s', target: 0 },  // Спад
    ],
    thresholds: {
        http_req_duration: ['p(95)<500'], // 95% HTTP запросов быстрее 500мс
    },
};

export default function () {
    const uniqueUser = `user_${randomString(5)}`;
    const password = 'password123';

    // Даты для запроса курсов (Rate)
    const dateTo = getFormattedDate(0); // Сегодня
    const dateFrom = getFormattedDate(-1); // Вчера

    // testing (GATEWAY)
    group('HTTP Gateway Flow', () => {

        // ping
        let resPing = http.get(`${BASE_URL}/ping`);
        check(resPing, { 'Ping status 200': (r) => r.status === 200 });

        // registration
        let registerPayload = JSON.stringify({ 
            username: uniqueUser,
            password: password 
        });
     
        let resReg = http.post(`${BASE_URL}/api/v1/register`, registerPayload, {
            headers: { 'Content-Type': 'application/json' },
        });

        console.log(`VU ${__VU}: Register status: ${resReg.status}, Response: ${resReg.body.substring(0, 100)}`);
        check(resReg, { 
            'Register successful': (r) => r.status === 201 || r.status === 200 
        });

        // login
        let loginPayload = JSON.stringify({ 
            username: uniqueUser,
            password: password 
        });
        
        let resLogin = http.post(`${BASE_URL}/api/v1/login`, loginPayload, {
            headers: { 'Content-Type': 'application/json' },
        });
        
        console.log(`VU ${__VU}: Login status: ${resLogin.status}, Response: ${resLogin.body.substring(0, 100)}`);
        
        let isLoginSuccessful = resLogin.status === 200;
        let authToken = null;
        
        if (isLoginSuccessful) {
            try {
                const loginResponse = JSON.parse(resLogin.body);
                authToken = loginResponse.token; 
                console.log(`VU ${__VU}: Got token: ${authToken ? 'Yes' : 'No'}`);
            } catch (e) {
                console.error(`VU ${__VU}: Failed to parse login response or token: ${e}`);
                isLoginSuccessful = false;
            }
        }
        
        check(resLogin, { 'Login successful': (r) => r.status === 200 });

        // get rate
        if (isLoginSuccessful && authToken) {
            let rateUrl = `${BASE_URL}/api/v1/rate?currency=eur&date_from=${dateFrom}&date_to=${dateTo}`;
            
            let resRate = http.get(rateUrl, {
                headers: {
                    'Authorization': `Bearer ${authToken}`
                }
            });
            
            console.log(`VU ${__VU}: Rate status: ${resRate.status}, Body: ${resRate.body.substring(0, 100)}`);
            
            check(resRate, { 
                'Get Rates HTTP 200': (r) => r.status === 200,
                'Response time < 500ms': (r) => r.timings.duration < 500
            });
        }
        
        // logout (not implement ex)
        if (authToken) {
            let resLogout = http.post(`${BASE_URL}/api/v1/logout`, null, {
                headers: {
                    'Authorization': `Bearer ${authToken}`
                }
            });
            
            console.log(`VU ${__VU}: Logout status: ${resLogout.status}, Body: ${resLogout.body.substring(0, 100)}`);
            
            check(resLogout, {
                'Logout successful': (r) => r.status === 200
            });
        }
    });

    sleep(0.5);
}