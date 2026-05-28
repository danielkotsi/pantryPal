/**
 * Planner Page Handler
 * Fetches weekly plan from backend, falls back to mock data.
 * Fetches recipe ingredients on demand when user expands a meal.
 */

const MEAL_SECTIONS = ['breakfast', 'lunch', 'dinner', 'snacks'];
const MEAL_LABELS = { breakfast: 'Breakfast', lunch: 'Lunch', dinner: 'Dinner', snacks: 'Snacks' };
const MEAL_ICONS = { breakfast: '🌅', lunch: '☀️', dinner: '🌙', snacks: '🍿' };
const DAY_NAMES = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
const DAY_SHORT = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];

class PlannerPageHandler {
    constructor() {
        this.view = 'week';
        this.currentDate = new Date();
        this.weekData = null;
        this.recipeCache = new Map();
        this.setupEventListeners();
    }

    setupEventListeners() {
        document.addEventListener('click', (e) => {
            const prevBtn = e.target.closest('#prevPeriod');
            const nextBtn = e.target.closest('#nextPeriod');
            const todayBtn = e.target.closest('#todayBtn');
            const toggleBtn = e.target.closest('.btn-toggle');
            const expandBtn = e.target.closest('.btn-expand');

            if (prevBtn) this.navigatePeriod(-1);
            if (nextBtn) this.navigatePeriod(1);
            if (todayBtn) this.goToToday();
            if (toggleBtn) this.toggleView(toggleBtn.dataset.view);
            if (expandBtn) this.handleExpand(expandBtn);
        });
    }

    handleExpand(btn) {
        const details = btn.closest('.meal-expand').querySelector('.meal-details');
        if (!details) return;
        const isVisible = details.style.display !== 'none';
        details.style.display = isVisible ? 'none' : 'block';
        btn.textContent = isVisible ? 'Details' : 'Hide';
        if (!isVisible && !details.dataset.loaded) {
            details.dataset.loaded = '1';
            const recipeId = btn.closest('.meal-card').dataset.recipeId;
            if (recipeId) {
                this.loadRecipeDetails(recipeId, details);
            }
        }
    }

    async loadRecipeDetails(recipeId, detailsEl) {
        try {
            const recipe = await api.getRecipe(recipeId);
            this.recipeCache.set(recipeId, recipe);
            if (recipe.ingredients && recipe.ingredients.length > 0) {
                const list = detailsEl.querySelector('.ingredient-list');
                if (list) {
                    list.innerHTML = recipe.ingredients.map(ing =>
                        `<li><span class="ing-qty">${ing.quantity}${ing.unit}</span> ${this.escapeHtml(ing.description)}</li>`
                    ).join('');
                }
            }
            if (recipe.instructions) {
                const instrEl = detailsEl.querySelector('.recipe-instructions');
                if (instrEl) {
                    instrEl.textContent = recipe.instructions.substring(0, 200);
                }
            }
        } catch (err) {
            // Backend not available; keep default ingredient list
        }
    }

    init() {
        this.loadWeek();
    }

    getWeekStart(date) {
        const d = new Date(date);
        const day = d.getDay();
        const diff = d.getDate() - day + (day === 0 ? -6 : 1);
        d.setDate(diff);
        d.setHours(0, 0, 0, 0);
        return d;
    }

    getMonthStart(date) {
        return new Date(date.getFullYear(), date.getMonth(), 1);
    }

    navigatePeriod(direction) {
        if (this.view === 'week') {
            this.currentDate.setDate(this.currentDate.getDate() + (direction * 7));
        } else {
            this.currentDate.setMonth(this.currentDate.getMonth() + direction);
        }
        this.loadWeek();
    }

    goToToday() {
        this.currentDate = new Date();
        this.loadWeek();
    }

    toggleView(view) {
        this.view = view;
        document.querySelectorAll('.btn-toggle').forEach(btn => {
            btn.classList.toggle('active', btn.dataset.view === view);
        });
        this.loadWeek();
    }

    loadWeek() {
        if (this.view === 'month') {
            this.loadMonthData();
        } else {
            this.loadWeekData();
        }
    }

    async loadWeekData() {
        const weekStart = this.getWeekStart(this.currentDate);
        const startStr = this.formatDate(weekStart);

        try {
            const backendData = await api.getWeekPlan(startStr);
            if (backendData && backendData.days && backendData.days.length > 0) {
                this.weekData = this.mapBackendWeek(backendData, weekStart);
                this.render();
                return;
            }
        } catch (err) {
            // Backend unavailable
        }

        // Fallback to mock data
        this.weekData = this.generateWeekData(weekStart);
        this.render();
    }

    async loadMonthData() {
        const start = this.getMonthStart(this.currentDate);
        const year = start.getFullYear();
        const month = start.getMonth();
        const totalDays = new Date(year, month + 1, 0).getDate();

        // Try to fetch backend data for each week in month
        const backendDays = new Map();
        try {
            for (let w = 0; w < 6; w++) {
                const weekDate = new Date(start);
                weekDate.setDate(weekDate.getDate() + (w * 7));
                const startStr = this.formatDate(weekDate);
                const data = await api.getWeekPlan(startStr);
                if (data && data.days) {
                    data.days.forEach(day => {
                        backendDays.set(day.date, day);
                    });
                }
            }
        } catch (err) {
            // Backend unavailable
        }

        const days = [];
        let monthKcal = 0, monthProtein = 0, monthCarbs = 0, monthFat = 0;

        for (let i = 0; i < totalDays; i++) {
            const dayDate = new Date(start);
            dayDate.setDate(start.getDate() + i);
            const dateStr = this.formatDate(dayDate);
            const bd = backendDays.get(dateStr);

            const day = {
                date: dateStr,
                meals: { breakfast: null, lunch: null, dinner: null, snacks: null },
                totals: { kcal: 0, protein: 0, carbs: 0, fat: 0 }
            };

            if (bd && bd.sections) {
                MEAL_SECTIONS.forEach(section => {
                    const mealResp = bd.sections[section];
                    if (mealResp) {
                        const meal = this.mapPlanMeal(mealResp);
                        day.meals[section] = meal;
                        day.totals.kcal += meal.kcal;
                        day.totals.protein += meal.protein;
                        day.totals.carbs += meal.carbs;
                        day.totals.fat += meal.fat;
                    }
                });
            }

            if (!bd || !bd.sections || !bd.sections.breakfast) {
                MEAL_SECTIONS.forEach((section, idx) => {
                    if (!day.meals[section]) {
                        const meal = this.generateMeal(dayDate, section, idx + i);
                        day.meals[section] = meal;
                        day.totals.kcal += meal.kcal;
                        day.totals.protein += meal.protein;
                        day.totals.carbs += meal.carbs;
                        day.totals.fat += meal.fat;
                    }
                });
            }

            day.totals.kcal = Math.round(day.totals.kcal);
            day.totals.protein = Math.round(day.totals.protein);
            day.totals.carbs = Math.round(day.totals.carbs);
            day.totals.fat = Math.round(day.totals.fat);

            monthKcal += day.totals.kcal;
            monthProtein += day.totals.protein;
            monthCarbs += day.totals.carbs;
            monthFat += day.totals.fat;

            days.push(day);
        }

        this.weekData = {
            startDate: this.formatDate(start),
            days,
            weekTotals: { kcal: monthKcal, protein: Math.round(monthProtein), carbs: Math.round(monthCarbs), fat: Math.round(monthFat) }
        };
        this.render();
    }

    mapBackendWeek(backendData, weekStart) {
        const backendDays = {};
        if (backendData.days) {
            backendData.days.forEach(day => {
                backendDays[day.date] = day;
            });
        }

        const days = [];
        let weekKcal = 0, weekProtein = 0, weekCarbs = 0, weekFat = 0;

        for (let i = 0; i < 7; i++) {
            const dayDate = new Date(weekStart);
            dayDate.setDate(weekStart.getDate() + i);
            const dateStr = this.formatDate(dayDate);
            const bd = backendDays[dateStr];

            const day = {
                date: dateStr,
                dayName: DAY_NAMES[dayDate.getDay()],
                dayShort: DAY_SHORT[dayDate.getDay()],
                meals: {},
                totals: { kcal: 0, protein: 0, carbs: 0, fat: 0 }
            };

            let dayKcal = 0, dayProtein = 0, dayCarbs = 0, dayFat = 0;

            MEAL_SECTIONS.forEach((section, idx) => {
                let meal = null;

                // Try backend first
                if (bd && bd.sections && bd.sections[section]) {
                    meal = this.mapPlanMeal(bd.sections[section]);
                }

                // Fall back to mock
                if (!meal) {
                    meal = this.generateMeal(dayDate, section, idx + i);
                }

                day.meals[section] = meal;
                dayKcal += meal.kcal;
                dayProtein += meal.protein;
                dayCarbs += meal.carbs;
                dayFat += meal.fat;
            });

            day.totals = {
                kcal: Math.round(dayKcal),
                protein: Math.round(dayProtein),
                carbs: Math.round(dayCarbs),
                fat: Math.round(dayFat)
            };

            weekKcal += day.totals.kcal;
            weekProtein += day.totals.protein;
            weekCarbs += day.totals.carbs;
            weekFat += day.totals.fat;

            days.push(day);
        }

        return {
            startDate: this.formatDate(weekStart),
            days,
            weekTotals: { kcal: weekKcal, protein: Math.round(weekProtein), carbs: Math.round(weekCarbs), fat: Math.round(weekFat) }
        };
    }

    mapPlanMeal(mealResp) {
        return {
            id: mealResp.id,
            recipeId: mealResp.recipeId,
            name: mealResp.recipeName,
            servings: mealResp.servings,
            kcal: Math.round(mealResp.macros ? mealResp.macros.calories : 0),
            protein: Math.round(mealResp.macros ? mealResp.macros.proteinG : 0),
            carbs: Math.round(mealResp.macros ? mealResp.macros.carbsG : 0),
            fat: Math.round(mealResp.macros ? mealResp.macros.fatG : 0),
            estimatedCostCents: mealResp.estimatedCostCents,
            isConsumed: mealResp.isConsumed,
            ingredients: [{ name: this.getDefaultIngredient(mealResp.recipeName || ''), qty: 1, unit: 'serving' }]
        };
    }

    getDefaultIngredient(recipeName) {
        return recipeName || 'Prepared meal';
    }

    generateWeekData(weekStart) {
        const days = [];

        for (let i = 0; i < 7; i++) {
            const dayDate = new Date(weekStart);
            dayDate.setDate(weekStart.getDate() + i);

            const meals = {};
            let dayKcal = 0, dayProtein = 0, dayCarbs = 0, dayFat = 0;

            MEAL_SECTIONS.forEach((section, idx) => {
                const meal = this.generateMeal(dayDate, section, idx + i);
                meals[section] = meal;
                dayKcal += meal.kcal;
                dayProtein += meal.protein;
                dayCarbs += meal.carbs;
                dayFat += meal.fat;
            });

            days.push({
                date: this.formatDate(dayDate),
                dayName: DAY_NAMES[dayDate.getDay()],
                dayShort: DAY_SHORT[dayDate.getDay()],
                meals,
                totals: { kcal: Math.round(dayKcal), protein: Math.round(dayProtein), carbs: Math.round(dayCarbs), fat: Math.round(dayFat) }
            });
        }

        let weekKcal = 0, weekProtein = 0, weekCarbs = 0, weekFat = 0;
        days.forEach(d => {
            weekKcal += d.totals.kcal;
            weekProtein += d.totals.protein;
            weekCarbs += d.totals.carbs;
            weekFat += d.totals.fat;
        });

        return {
            startDate: this.formatDate(weekStart),
            days,
            weekTotals: { kcal: weekKcal, protein: Math.round(weekProtein), carbs: Math.round(weekCarbs), fat: Math.round(weekFat) }
        };
    }

    mealsForSection(section, idx) {
        const meals = {
            breakfast: [
                { name: 'Oatmeal with Berries', ingredients: [{ name: 'Rolled Oats', qty: 80, unit: 'g' }, { name: 'Mixed Berries', qty: 100, unit: 'g' }, { name: 'Honey', qty: 15, unit: 'g' }], kcal: 350, protein: 12, carbs: 58, fat: 7 },
                { name: 'Scrambled Eggs & Toast', ingredients: [{ name: 'Eggs', qty: 3, unit: 'pc' }, { name: 'Whole Wheat Bread', qty: 2, unit: 'slice' }, { name: 'Butter', qty: 10, unit: 'g' }], kcal: 420, protein: 24, carbs: 32, fat: 22 },
                { name: 'Greek Yogurt Parfait', ingredients: [{ name: 'Greek Yogurt', qty: 200, unit: 'ml' }, { name: 'Granola', qty: 40, unit: 'g' }, { name: 'Banana', qty: 1, unit: 'pc' }], kcal: 380, protein: 20, carbs: 48, fat: 10 },
                { name: 'Green Smoothie', ingredients: [{ name: 'Spinach', qty: 60, unit: 'g' }, { name: 'Apple', qty: 1, unit: 'pc' }, { name: 'Almond Milk', qty: 250, unit: 'ml' }, { name: 'Protein Powder', qty: 30, unit: 'g' }], kcal: 310, protein: 28, carbs: 35, fat: 6 },
                { name: 'Avocado Toast', ingredients: [{ name: 'Sourdough Bread', qty: 2, unit: 'slice' }, { name: 'Avocado', qty: 1, unit: 'pc' }, { name: 'Cherry Tomatoes', qty: 50, unit: 'g' }], kcal: 365, protein: 10, carbs: 30, fat: 22 },
                { name: 'Pancakes & Syrup', ingredients: [{ name: 'Pancake Mix', qty: 100, unit: 'g' }, { name: 'Maple Syrup', qty: 30, unit: 'ml' }, { name: 'Blueberries', qty: 50, unit: 'g' }], kcal: 450, protein: 10, carbs: 72, fat: 12 },
                { name: 'Chia Pudding', ingredients: [{ name: 'Chia Seeds', qty: 30, unit: 'g' }, { name: 'Coconut Milk', qty: 200, unit: 'ml' }, { name: 'Mango', qty: 80, unit: 'g' }], kcal: 340, protein: 10, carbs: 28, fat: 20 }
            ],
            lunch: [
                { name: 'Grilled Chicken Salad', ingredients: [{ name: 'Chicken Breast', qty: 150, unit: 'g' }, { name: 'Mixed Greens', qty: 100, unit: 'g' }, { name: 'Olive Oil', qty: 15, unit: 'ml' }], kcal: 480, protein: 42, carbs: 12, fat: 28 },
                { name: 'Quinoa Bowl', ingredients: [{ name: 'Quinoa', qty: 120, unit: 'g' }, { name: 'Black Beans', qty: 80, unit: 'g' }, { name: 'Corn', qty: 60, unit: 'g' }, { name: 'Avocado', qty: 0.5, unit: 'pc' }], kcal: 520, protein: 18, carbs: 68, fat: 18 },
                { name: 'Turkey Sandwich', ingredients: [{ name: 'Turkey Breast', qty: 100, unit: 'g' }, { name: 'Sourdough Bread', qty: 2, unit: 'slice' }, { name: 'Lettuce & Tomato', qty: 30, unit: 'g' }, { name: 'Mustard', qty: 10, unit: 'g' }], kcal: 410, protein: 32, carbs: 38, fat: 12 },
                { name: 'Salmon & Rice', ingredients: [{ name: 'Salmon Fillet', qty: 150, unit: 'g' }, { name: 'Brown Rice', qty: 150, unit: 'g' }, { name: 'Asparagus', qty: 80, unit: 'g' }], kcal: 560, protein: 40, carbs: 45, fat: 22 },
                { name: 'Vegetable Stir-Fry', ingredients: [{ name: 'Tofu', qty: 120, unit: 'g' }, { name: 'Broccoli', qty: 80, unit: 'g' }, { name: 'Bell Peppers', qty: 60, unit: 'g' }, { name: 'Soy Sauce', qty: 15, unit: 'ml' }], kcal: 380, protein: 22, carbs: 28, fat: 18 },
                { name: 'Lentil Soup', ingredients: [{ name: 'Red Lentils', qty: 100, unit: 'g' }, { name: 'Carrots', qty: 50, unit: 'g' }, { name: 'Celery', qty: 30, unit: 'g' }, { name: 'Bread', qty: 1, unit: 'slice' }], kcal: 440, protein: 24, carbs: 62, fat: 8 },
                { name: 'Tuna Wrap', ingredients: [{ name: 'Canned Tuna', qty: 100, unit: 'g' }, { name: 'Whole Wheat Wrap', qty: 1, unit: 'pc' }, { name: 'Mixed Greens', qty: 40, unit: 'g' }, { name: 'Greek Yogurt', qty: 30, unit: 'ml' }], kcal: 390, protein: 34, carbs: 30, fat: 14 }
            ],
            dinner: [
                { name: 'Pasta Bolognese', ingredients: [{ name: 'Spaghetti', qty: 150, unit: 'g' }, { name: 'Ground Beef', qty: 120, unit: 'g' }, { name: 'Tomato Sauce', qty: 100, unit: 'ml' }, { name: 'Parmesan', qty: 15, unit: 'g' }], kcal: 680, protein: 36, carbs: 72, fat: 26 },
                { name: 'Chicken Tikka Masala', ingredients: [{ name: 'Chicken Thighs', qty: 150, unit: 'g' }, { name: 'Basmati Rice', qty: 120, unit: 'g' }, { name: 'Tikka Sauce', qty: 100, unit: 'ml' }], kcal: 620, protein: 40, carbs: 58, fat: 24 },
                { name: 'Beef Stir-Fry', ingredients: [{ name: 'Beef Strips', qty: 130, unit: 'g' }, { name: 'Noodles', qty: 150, unit: 'g' }, { name: 'Mixed Vegetables', qty: 100, unit: 'g' }], kcal: 590, protein: 34, carbs: 62, fat: 20 },
                { name: 'Baked Cod & Potatoes', ingredients: [{ name: 'Cod Fillet', qty: 160, unit: 'g' }, { name: 'Potatoes', qty: 200, unit: 'g' }, { name: 'Green Beans', qty: 80, unit: 'g' }], kcal: 510, protein: 38, carbs: 48, fat: 14 },
                { name: 'Vegetable Curry', ingredients: [{ name: 'Chickpeas', qty: 150, unit: 'g' }, { name: 'Coconut Milk', qty: 100, unit: 'ml' }, { name: 'Rice', qty: 120, unit: 'g' }, { name: 'Spinach', qty: 50, unit: 'g' }], kcal: 550, protein: 20, carbs: 68, fat: 22 },
                { name: 'Grilled Steak & Salad', ingredients: [{ name: 'Sirloin Steak', qty: 180, unit: 'g' }, { name: 'Sweet Potato', qty: 150, unit: 'g' }, { name: 'Arugula', qty: 50, unit: 'g' }], kcal: 640, protein: 48, carbs: 32, fat: 34 },
                { name: 'Shrimp Pasta', ingredients: [{ name: 'Penne', qty: 130, unit: 'g' }, { name: 'Shrimp', qty: 120, unit: 'g' }, { name: 'Garlic Butter', qty: 20, unit: 'g' }, { name: 'Parsley', qty: 5, unit: 'g' }], kcal: 580, protein: 34, carbs: 58, fat: 22 }
            ],
            snacks: [
                { name: 'Apple & Almond Butter', ingredients: [{ name: 'Apple', qty: 1, unit: 'pc' }, { name: 'Almond Butter', qty: 20, unit: 'g' }], kcal: 220, protein: 6, carbs: 28, fat: 12 },
                { name: 'Trail Mix', ingredients: [{ name: 'Mixed Nuts', qty: 30, unit: 'g' }, { name: 'Dried Cranberries', qty: 15, unit: 'g' }, { name: 'Dark Chocolate', qty: 10, unit: 'g' }], kcal: 190, protein: 6, carbs: 18, fat: 14 },
                { name: 'Hummus & Veggies', ingredients: [{ name: 'Hummus', qty: 60, unit: 'g' }, { name: 'Carrot Sticks', qty: 50, unit: 'g' }, { name: 'Cucumber', qty: 50, unit: 'g' }], kcal: 160, protein: 6, carbs: 14, fat: 10 },
                { name: 'Protein Shake', ingredients: [{ name: 'Protein Powder', qty: 30, unit: 'g' }, { name: 'Almond Milk', qty: 250, unit: 'ml' }, { name: 'Banana', qty: 0.5, unit: 'pc' }], kcal: 240, protein: 28, carbs: 22, fat: 4 },
                { name: 'Rice Cakes & Avocado', ingredients: [{ name: 'Rice Cakes', qty: 2, unit: 'pc' }, { name: 'Avocado', qty: 0.5, unit: 'pc' }, { name: 'Everything Seasoning', qty: 2, unit: 'g' }], kcal: 180, protein: 4, carbs: 20, fat: 10 },
                { name: 'Greek Yogurt & Honey', ingredients: [{ name: 'Greek Yogurt', qty: 150, unit: 'ml' }, { name: 'Honey', qty: 10, unit: 'g' }, { name: 'Walnuts', qty: 15, unit: 'g' }], kcal: 210, protein: 14, carbs: 22, fat: 8 },
                { name: 'Energy Balls', ingredients: [{ name: 'Dates', qty: 40, unit: 'g' }, { name: 'Oats', qty: 20, unit: 'g' }, { name: 'Cocoa Powder', qty: 5, unit: 'g' }], kcal: 170, protein: 4, carbs: 30, fat: 6 }
            ]
        };
        return meals[section][idx % 7];
    }

    generateMeal(date, section, idx) {
        const mock = this.mealsForSection(section, idx);
        return { ...mock, id: null, recipeId: null };
    }

    formatDate(d) {
        return d.toISOString().split('T')[0];
    }

    render() {
        if (!this.weekData) return;

        if (this.view === 'week') {
            this.renderWeekView();
        } else {
            this.renderMonthView();
        }
        this.renderPeriodTotals();
        this.renderPeriodLabel();
    }

    renderPeriodLabel() {
        const label = document.getElementById('periodLabel');
        if (!label) return;

        if (this.view === 'week') {
            const start = new Date(this.weekData.startDate);
            const end = new Date(start);
            end.setDate(start.getDate() + 6);
            const opts = { month: 'long', day: 'numeric', year: 'numeric' };
            label.textContent = `${start.toLocaleDateString('en-US', opts)} — ${end.toLocaleDateString('en-US', opts)}`;
        } else {
            const opts = { month: 'long', year: 'numeric' };
            label.textContent = this.currentDate.toLocaleDateString('en-US', opts);
        }
    }

    renderWeekView() {
        const container = document.getElementById('plannerContent');
        if (!container) return;

        let html = '<div class="week-grid">';

        this.weekData.days.forEach(day => {
            const isToday = this.isToday(day.date);
            html += `
                <div class="day-column ${isToday ? 'today' : ''}">
                    <div class="day-header">
                        <span class="day-name">${day.dayShort}</span>
                        <span class="day-date">${new Date(day.date).getDate()}</span>
                    </div>
                    <div class="day-meals">
            `;

            MEAL_SECTIONS.forEach(section => {
                const meal = day.meals[section];
                html += this.renderMealCard(section, meal);
            });

            html += `
                    </div>
                    <div class="day-totals">
                        <div class="macro kcal"><span class="macro-val">${day.totals.kcal}</span> kcal</div>
                        <div class="macro-row">
                            <span class="macro p">${day.totals.protein}g</span>
                            <span class="macro c">${day.totals.carbs}g</span>
                            <span class="macro f">${day.totals.fat}g</span>
                        </div>
                    </div>
                </div>
            `;
        });

        html += '</div>';
        container.innerHTML = html;
    }

    renderMealCard(section, meal) {
        if (!meal) {
            return `
                <div class="meal-card empty" data-section="${section}">
                    <div class="meal-header">
                        <span class="meal-icon">${MEAL_ICONS[section]}</span>
                        <span class="meal-label">${MEAL_LABELS[section]}</span>
                    </div>
                    <div class="meal-name empty-text">No meal planned</div>
                </div>
            `;
        }

        const ingredients = meal.ingredients.map(ing =>
            `${ing.qty}${ing.unit} ${ing.name}`
        ).join(', ');

        return `
            <div class="meal-card" data-section="${section}" data-recipe-id="${meal.recipeId || ''}">
                <div class="meal-header">
                    <span class="meal-icon">${MEAL_ICONS[section]}</span>
                    <span class="meal-label">${MEAL_LABELS[section]}</span>
                </div>
                <div class="meal-name" title="${this.escapeHtml(meal.name)}">${this.escapeHtml(meal.name)}</div>
                <div class="meal-ingredients" title="${this.escapeHtml(ingredients)}">
                    <span class="ing-label">Ingredients:</span>
                    ${this.escapeHtml(ingredients.length > 50 ? ingredients.substring(0, 50) + '...' : ingredients)}
                </div>
                <div class="meal-macros">
                    <span class="macro kcal">${meal.kcal} kcal</span>
                    <span class="macro p">P: ${meal.protein}g</span>
                    <span class="macro c">C: ${meal.carbs}g</span>
                    <span class="macro f">F: ${meal.fat}g</span>
                </div>
                <div class="meal-expand" data-section="${section}">
                    <button class="btn btn-small btn-expand">Details</button>
                    <div class="meal-details" style="display:none;">
                        <ul class="ingredient-list">
                            ${meal.ingredients.map(ing => `
                                <li><span class="ing-qty">${ing.qty}${ing.unit}</span> ${this.escapeHtml(ing.name)}</li>
                            `).join('')}
                        </ul>
                        ${meal.recipeId ? '<p class="recipe-instructions" style="font-size:0.65rem;color:var(--text-secondary);margin-top:0.3rem;"></p>' : ''}
                    </div>
                </div>
            </div>
        `;
    }

    renderMonthView() {
        const container = document.getElementById('plannerContent');
        if (!container || !this.weekData) return;

        const year = this.currentDate.getFullYear();
        const month = this.currentDate.getMonth();
        const firstDay = new Date(year, month, 1);
        const lastDay = new Date(year, month + 1, 0);
        const startPad = firstDay.getDay();
        const totalDays = lastDay.getDate();

        const dayMap = {};
        this.weekData.days.forEach(d => { dayMap[d.date] = d; });

        let html = '<div class="month-grid">';
        html += '<div class="month-header-row">';
        DAY_SHORT.forEach(d => { html += `<div class="month-header-cell">${d}</div>`; });
        html += '</div>';
        html += '<div class="month-body">';

        for (let i = 0; i < startPad; i++) {
            html += '<div class="month-cell empty"></div>';
        }

        for (let day = 1; day <= totalDays; day++) {
            const dateStr = `${year}-${String(month + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}`;
            const isToday = this.isToday(dateStr);
            const dayData = dayMap[dateStr];

            html += `<div class="month-cell ${isToday ? 'today' : ''}">`;
            html += `<div class="month-day-num">${day}</div>`;

            if (dayData) {
                MEAL_SECTIONS.forEach(section => {
                    const meal = dayData.meals[section];
                    if (meal) {
                        html += `
                            <div class="month-meal" title="${this.escapeHtml(meal.name)}">
                                <span class="month-meal-icon">${MEAL_ICONS[section]}</span>
                                <span class="month-meal-name">${this.escapeHtml(meal.name)}</span>
                            </div>
                        `;
                    }
                });
                html += `<div class="month-day-kcal">${dayData.totals.kcal} kcal</div>`;
            }

            html += '</div>';
        }

        html += '</div></div>';
        container.innerHTML = html;
    }

    renderPeriodTotals() {
        const container = document.getElementById('weekTotals');
        if (!container || !this.weekData) return;

        const t = this.weekData.weekTotals;
        const isMonth = this.view === 'month';
        const daysInPeriod = this.weekData.days.length;
        const dailyAvg = daysInPeriod > 0 ? Math.round(t.kcal / daysInPeriod) : 0;

        container.innerHTML = `
            <div class="totals-bar">
                <div class="totals-title">${isMonth ? 'Month' : 'Week'} Summary</div>
                <div class="totals-macros">
                    <div class="total-item">
                        <span class="total-val">${t.kcal.toLocaleString()}</span>
                        <span class="total-label">Total kcal</span>
                    </div>
                    <div class="total-item">
                        <span class="total-val">${dailyAvg.toLocaleString()}</span>
                        <span class="total-label">Daily avg</span>
                    </div>
                    <div class="total-item">
                        <span class="total-val">${t.protein}g</span>
                        <span class="total-label">Protein</span>
                    </div>
                    <div class="total-item">
                        <span class="total-val">${t.carbs}g</span>
                        <span class="total-label">Carbs</span>
                    </div>
                    <div class="total-item">
                        <span class="total-val">${t.fat}g</span>
                        <span class="total-label">Fat</span>
                    </div>
                </div>
                ${t.kcal > 0 ? `
                <div class="totals-bar-chart">
                    <div class="bar bar-protein" style="width:${(t.protein * 4 / t.kcal * 100).toFixed(1)}%"></div>
                    <div class="bar bar-carbs" style="width:${(t.carbs * 4 / t.kcal * 100).toFixed(1)}%"></div>
                    <div class="bar bar-fat" style="width:${(t.fat * 9 / t.kcal * 100).toFixed(1)}%"></div>
                </div>
                <div class="totals-bar-legend">
                    <span><span class="dot protein"></span> Protein ${(t.protein * 4 / t.kcal * 100).toFixed(0)}%</span>
                    <span><span class="dot carbs"></span> Carbs ${(t.carbs * 4 / t.kcal * 100).toFixed(0)}%</span>
                    <span><span class="dot fat"></span> Fat ${(t.fat * 9 / t.kcal * 100).toFixed(0)}%</span>
                </div>
                ` : ''}
            </div>
        `;
    }

    isToday(dateStr) {
        const today = new Date();
        return dateStr === this.formatDate(today);
    }

    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
}

const plannerPageHandler = new PlannerPageHandler();
