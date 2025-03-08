// Функція для форматування чисел з певною кількістю десяткових знаків
function formatNumber(number, decimals = 2) {
    return parseFloat(number).toFixed(decimals);
}

// Обробка форми вибору кабелів
const cableSelectionForm = document.getElementById('cable-selection-form');
const cableResult = document.getElementById('cable-result');

cableSelectionForm.addEventListener('submit', async function(event) {
    event.preventDefault();
    
    const formData = {
        sn: parseFloat(document.getElementById('sn').value),
        un: parseFloat(document.getElementById('un').value),
        cosFi: parseFloat(document.getElementById('cosFi').value),
        length: parseFloat(document.getElementById('length').value),
        installMethod: document.getElementById('installMethod').value,
        cableCount: parseInt(document.getElementById('cableCount').value),
        ambientTemp: parseFloat(document.getElementById('ambientTemp').value),
        soilThermalRes: parseFloat(document.getElementById('soilThermalRes').value)
    };
    
    try {
        const response = await fetch('/api/cable-selection', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(formData)
        });
        
        if (!response.ok) {
            throw new Error('Помилка на сервері');
        }
        
        const result = await response.json();
        
        // Відображення результатів
        document.getElementById('cable-type').textContent = result.cableType || 'Не знайдено';
        document.getElementById('cable-section').textContent = result.cableSection ? formatNumber(result.cableSection, 0) : 'Не знайдено';
        document.getElementById('rated-current').textContent = formatNumber(result.ratedCurrent, 1);
        document.getElementById('permissible-load').textContent = result.permissibleLoad ? formatNumber(result.permissibleLoad, 1) : 'Не знайдено';
        document.getElementById('voltage-drop').textContent = formatNumber(result.voltageDropPerc, 2);
        document.getElementById('cable-valid').textContent = result.isValid ? 'Кабель відповідає вимогам' : 'Кабель не відповідає вимогам';
        document.getElementById('cable-valid').className = result.isValid ? 'text-success' : 'text-danger';
        
        // Показ блоку результатів
        cableResult.style.display = 'block';
        cableResult.classList.add('fade-in');
    } catch (error) {
        console.error('Помилка:', error);
        alert('Виникла помилка при розрахунку. Будь ласка, спробуйте ще раз.');
    }
});

// Обробка форми розрахунку струмів КЗ на шинах ГПП
const scCurrentsForm = document.getElementById('sc-currents-form');
const scCurrentsResult = document.getElementById('sc-currents-result');

scCurrentsForm.addEventListener('submit', async function(event) {
    event.preventDefault();
    
    const formData = {
        systemVoltage: parseFloat(document.getElementById('systemVoltage').value),
        systemPower: parseFloat(document.getElementById('systemPower').value),
        xc: parseFloat(document.getElementById('xc').value),
        rc: parseFloat(document.getElementById('rc').value),
        l1: parseFloat(document.getElementById('l1').value),
        x1: parseFloat(document.getElementById('x1').value),
        r1: parseFloat(document.getElementById('r1').value),
        l2: parseFloat(document.getElementById('l2').value),
        x2: parseFloat(document.getElementById('x2').value),
        r2: parseFloat(document.getElementById('r2').value)
    };
    
    try {
        const response = await fetch('/api/sc-currents', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(formData)
        });
        
        if (!response.ok) {
            throw new Error('Помилка на сервері');
        }
        
        const result = await response.json();
        
        // Відображення результатів
        document.getElementById('three-phase-current').textContent = formatNumber(result.threePhaseCurrent, 2);
        document.getElementById('single-phase-current').textContent = formatNumber(result.singlePhaseCurrent, 2);
        document.getElementById('impulse-current').textContent = formatNumber(result.impulseCurrentValue, 2);
        document.getElementById('thermal-impulse').textContent = formatNumber(result.thermalImpulse, 2);
        document.getElementById('thermal-stability').textContent = result.thermalStability ? 'Забезпечено' : 'Не забезпечено';
        document.getElementById('thermal-stability').className = result.thermalStability ? 'text-success' : 'text-danger';
        document.getElementById('thermal-stability-min').textContent = formatNumber(result.thermalStabilityMin, 2);
        document.getElementById('dynamic-stability-force').textContent = formatNumber(result.dynamicStabilityForce, 2);
        document.getElementById('dynamic-stability').textContent = result.dynamicStability ? 'Забезпечено' : 'Не забезпечено';
        document.getElementById('dynamic-stability').className = result.dynamicStability ? 'text-success' : 'text-danger';
        
        // Показ блоку результатів
        scCurrentsResult.style.display = 'block';
        scCurrentsResult.classList.add('fade-in');
    } catch (error) {
        console.error('Помилка:', error);
        alert('Виникла помилка при розрахунку. Будь ласка, спробуйте ще раз.');
    }
});

// Обробка форми розрахунку струмів КЗ для підстанції ХПнЕМ
const emScCurrentsForm = document.getElementById('em-sc-currents-form');
const emScCurrentsResult = document.getElementById('em-sc-currents-result');

emScCurrentsForm.addEventListener('submit', async function(event) {
    event.preventDefault();
    
    const formData = {
        systemVoltage: parseFloat(document.getElementById('emSystemVoltage').value),
        normalModeSc: parseFloat(document.getElementById('normalModeSc').value),
        minModeSc: parseFloat(document.getElementById('minModeSc').value),
        emergencyModeSc: parseFloat(document.getElementById('emergencyModeSc').value)
    };
    
    try {
        const response = await fetch('/api/em-sc-currents', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(formData)
        });
        
        if (!response.ok) {
            throw new Error('Помилка на сервері');
        }
        
        const result = await response.json();
        
        // Відображення результатів
        document.getElementById('normal-mode-three-phase').textContent = formatNumber(result.normalModeThreePhaseSC, 2);
        document.getElementById('normal-mode-single-phase').textContent = formatNumber(result.normalModeSinglePhaseSC, 2);
        document.getElementById('min-mode-three-phase').textContent = formatNumber(result.minModeThreePhaseSC, 2);
        document.getElementById('min-mode-single-phase').textContent = formatNumber(result.minModeSinglePhaseSC, 2);
        document.getElementById('emergency-mode-three-phase').textContent = formatNumber(result.emergencyModeThreePhaseSC, 2);
        document.getElementById('emergency-mode-single-phase').textContent = formatNumber(result.emergencyModeSinglePhaseSC, 2);
        
        // Показ блоку результатів
        emScCurrentsResult.style.display = 'block';
        emScCurrentsResult.classList.add('fade-in');
    } catch (error) {
        console.error('Помилка:', error);
        alert('Виникла помилка при розрахунку. Будь ласка, спробуйте ще раз.');
    }
});

// Додаємо обробник для зміни вкладок, щоб приховати результати
document.querySelectorAll('.nav-link').forEach(tab => {
    tab.addEventListener('click', function() {
        // Приховуємо всі блоки результатів
        cableResult.style.display = 'none';
        scCurrentsResult.style.display = 'none';
        emScCurrentsResult.style.display = 'none';
    });
});

// Ініціалізація контрольних прикладів (для демонстрації)
document.addEventListener('DOMContentLoaded', function() {
    // Приклад 1: Вибір кабелів
    const cableExample = document.createElement('button');
    cableExample.textContent = 'Контрольний приклад';
    cableExample.className = 'btn btn-outline-secondary ms-2';
    cableExample.addEventListener('click', function() {
        document.getElementById('sn').value = '1000'; // 1000 кВА
        document.getElementById('un').value = '10'; // 10 кВ
        document.getElementById('cosFi').value = '0.9';
        document.getElementById('length').value = '0.5'; // 0.5 км
        document.getElementById('installMethod').value = 'ground'; // У землі
        document.getElementById('cableCount').value = '1';
        document.getElementById('ambientTemp').value = '25'; // 25°C
        document.getElementById('soilThermalRes').value = '1.2';
    });
    document.querySelector('#cable-selection-form button[type="submit"]').after(cableExample);
    
    // Приклад 2: Розрахунок струмів КЗ на шинах ГПП
    const scExample = document.createElement('button');
    scExample.textContent = 'Контрольний приклад';
    scExample.className = 'btn btn-outline-secondary ms-2';
    scExample.addEventListener('click', function() {
        document.getElementById('systemVoltage').value = '10'; // 10 кВ
        document.getElementById('systemPower').value = '500'; // 500 МВА
        document.getElementById('xc').value = '0.5'; // 0.5 Ом
        document.getElementById('rc').value = '0.1'; // 0.1 Ом
        document.getElementById('l1').value = '5'; // 5 км
        document.getElementById('x1').value = '0.4'; // 0.4 Ом/км
        document.getElementById('r1').value = '0.2'; // 0.2 Ом/км
        document.getElementById('l2').value = '3'; // 3 км
        document.getElementById('x2').value = '0.35'; // 0.35 Ом/км
        document.getElementById('r2').value = '0.15'; // 0.15 Ом/км
    });
    document.querySelector('#sc-currents-form button[type="submit"]').after(scExample);
    
    // Приклад 3: Розрахунок струмів КЗ для підстанції ХПнЕМ
    const emExample = document.createElement('button');
    emExample.textContent = 'Контрольний приклад';
    emExample.className = 'btn btn-outline-secondary ms-2';
    emExample.addEventListener('click', function() {
        document.getElementById('emSystemVoltage').value = '10'; // 10 кВ
        document.getElementById('normalModeSc').value = '400'; // 400 МВА
        document.getElementById('minModeSc').value = '300'; // 300 МВА
        document.getElementById('emergencyModeSc').value = '250'; // 250 МВА
    });
    document.querySelector('#em-sc-currents-form button[type="submit"]').after(emExample);
});