class InterviewApp {
    constructor() {
        this.questions = [];
        this.availableTags = [
            'golang', 'sql', 'database', 'concurrency', 'algorithms',
            'data-structures', 'net', 'linux', 'docker', 'git'
        ];
        this.init();
    }

    init() {
        this.bindEvents();
        this.loadQuestions();
    }

    bindEvents() {
        // Кнопка добавления вопроса
        document.getElementById('addQuestionBtn').addEventListener('click', () => {
            this.openModal();
        });

        // Закрытие модального окна
        document.querySelector('.close').addEventListener('click', () => {
            this.closeModal();
        });

        // Форма добавления вопроса
        document.getElementById('questionForm').addEventListener('submit', (e) => {
            e.preventDefault();
            this.createQuestion();
        });

        // Обработчик ввода тегов для показа подсказок
        document.getElementById('questionTags').addEventListener('input', (e) => {
            this.showTagSuggestions(e.target.value);
        });

        // Обработчик фокуса на поле тегов
        document.getElementById('questionTags').addEventListener('focus', (e) => {
            this.showTagSuggestions(e.target.value);
        });

        // Закрытие модального окна при клике вне его
        window.addEventListener('click', (e) => {
            const modal = document.getElementById('questionModal');
            if (e.target === modal) {
                this.closeModal();
                this.hideTagSuggestions();
            }
        });
    }

    // Показать подсказки тегов
    showTagSuggestions(inputValue) {
        this.hideTagSuggestions();
        
        if (!inputValue.trim()) {
            // Если поле пустое, показываем все теги
            this.renderAllAvailableTags();
            return;
        }

        const inputTags = inputValue.split(',').map(tag => tag.trim().toLowerCase());
        const currentTag = inputTags[inputTags.length - 1];
        
        if (!currentTag) {
            this.renderAllAvailableTags();
            return;
        }

        const filteredTags = this.availableTags.filter(tag => 
            tag.toLowerCase().includes(currentTag.toLowerCase()) && 
            !inputTags.includes(tag)
        );

        this.renderTagSuggestions(filteredTags, currentTag);
    }

     // Скрыть подсказки тегов
    hideTagSuggestions() {
        const existingSuggestions = document.getElementById('tagsSuggestions');
        if (existingSuggestions) {
            existingSuggestions.remove();
        }
    }

    // Показать все доступные теги
    renderAllAvailableTags() {
        this.hideTagSuggestions();
        
        const tagsInput = document.getElementById('questionTags');
        const suggestionsDiv = document.createElement('div');
        suggestionsDiv.id = 'tagsSuggestions';
        suggestionsDiv.className = 'tags-suggestions';
        
        suggestionsDiv.innerHTML = `
            <div class="suggestions-title">Возможные теги:</div>
            <div class="available-tags">
                ${this.availableTags.map(tag => 
                    `<span class="available-tag tag-${tag}" onclick="interviewApp.addTagToInput('${tag}')">${tag}</span>`
                ).join('')}
            </div>
        `;
        
        tagsInput.parentNode.insertBefore(suggestionsDiv, tagsInput.nextSibling);
    }

    // Показать отфильтрованные теги
    renderTagSuggestions(tags, currentTag) {
        if (tags.length === 0) return;
        
        this.hideTagSuggestions();
        
        const tagsInput = document.getElementById('questionTags');
        const suggestionsDiv = document.createElement('div');
        suggestionsDiv.id = 'tagsSuggestions';
        suggestionsDiv.className = 'tags-suggestions';
        
        suggestionsDiv.innerHTML = `
            <div class="suggestions-title">Возможные теги:</div>
            <div class="available-tags">
                ${tags.map(tag => 
                    `<span class="available-tag tag-${tag}" onclick="interviewApp.addTagToInput('${tag}')">${tag}</span>`
                ).join('')}
            </div>
        `;
        
        tagsInput.parentNode.insertBefore(suggestionsDiv, tagsInput.nextSibling);
    }

    // Добавить тег в поле ввода
    addTagToInput(tag) {
        const tagsInput = document.getElementById('questionTags');
        const currentTags = tagsInput.value.split(',').map(t => t.trim()).filter(t => t);
        
        // Убираем текущий частично введенный тег
        if (currentTags.length > 0) {
            currentTags.pop();
        }
        
        // Добавляем выбранный тег
        currentTags.push(tag);
        tagsInput.value = currentTags.join(', ');
        
        // Фокусируемся на поле и скрываем подсказки
        tagsInput.focus();
        this.hideTagSuggestions();
        
        // Показываем подсказки снова для возможного добавления следующего тега
        setTimeout(() => this.showTagSuggestions(tagsInput.value), 100);
    }

    async loadQuestions() {
        try {
            console.log('Начинаю загрузку вопросов...');
            const response = await fetch('/api/v1/questions');
            console.log('Статус ответа:', response.status);
            
            const data = await response.json();
            console.log('Получены вопросы:', data);
            
            this.questions = data;
            this.renderQuestions();
        } catch (error) {
            console.error('Ошибка загрузки вопросов:', error);
        }
    }

    renderQuestions() {
        const container = document.getElementById('questionsList');
        container.innerHTML = '';

        this.questions.forEach(question => {
            const questionElement = this.createQuestionElement(question);
            container.appendChild(questionElement);
        });
    }

    createQuestionElement(question) {
        const div = document.createElement('div');
        div.className = 'question-item';
        
        const id = question.id;
        if (!id) {
            console.error('Question has no ID:', question);
        }

        // Безопасное получение значений
        const title = question.title || 'Без названия';
        const answer = question.answer || 'Нет ответа';
        const tags = question.tags || [];
            
        div.innerHTML = `
            <div class="question-header">
                <div class="question-title">${this.escapeHtml(title)}</div>
                <div class="question-actions">
                    <div class="tags">
                        ${tags.map(tag => 
                            `<span class="tag tag-${String(tag).toLowerCase()}">${this.escapeHtml(tag)}</span>`
                        ).join('')}
                    </div>
                    <button class="delete-btn" data-id="${question.ID}">🗑️</button>
                </div>
            </div>
            <div class="question-answer">
                <div class="answer-content">${this.escapeHtml(answer)}</div>
            </div>
        `;

        // Обработчик клика для раскрытия/скрытия ответа
        div.addEventListener('click', (e) => {
            if (!e.target.classList.contains('tag')) {
                const answer = div.querySelector('.question-answer');
                answer.classList.toggle('expanded');
            }
        });

        // Обработчик удаления
        const deleteBtn = div.querySelector('.delete-btn');
        deleteBtn.addEventListener('click', (e) => {
            e.stopPropagation(); // Предотвращаем раскрытие вопроса
            this.deleteQuestion(id);
        });

        return div;
    }

    async createQuestion() {
        const title = document.getElementById('questionTitle').value;
        const answer = document.getElementById('questionAnswer').value;
        const tags = document.getElementById('questionTags').value
            .split(',')
            .map(tag => tag.trim())
            .filter(tag => tag);

        try {
            const response = await fetch('/api/v1/questions', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ title, answer, tags }),
            });

            if (response.ok) {
                this.closeModal();
                this.loadQuestions(); // Перезагружаем список
                document.getElementById('questionForm').reset();
            } else {
                alert('Ошибка при создании вопроса');
            }
        } catch (error) {
            console.error('Ошибка:', error);
            alert('Ошибка при создании вопроса');
        }
    }

    async deleteQuestion(id) {
        if (!confirm('Вы уверены, что хотите удалить этот вопрос?')) {
            return;
        }

        try {
            const response = await fetch(`/api/v1/questions/del/${id}`, {
                method: 'DELETE',
            });

            if (response.ok) {
                // ИСПРАВЛЕНО: используем lowercase id
                this.questions = this.questions.filter(q => q.id !== id);
                this.renderQuestions();
            } else {
                const error = await response.json();
                alert(`Ошибка при удалении: ${error.error}`);
            }
        } catch (error) {
            console.error('Ошибка:', error);
            alert('Ошибка при удалении вопроса');
        }
    }

    openModal() {
        document.getElementById('questionModal').style.display = 'block';
        // Показываем подсказки при открытии модального окна
        setTimeout(() => this.showTagSuggestions(''), 100);
    }

    closeModal() {
        document.getElementById('questionModal').style.display = 'none';
        this.hideTagSuggestions();
    }

    escapeHtml(unsafe) {
        return unsafe
            .replace(/&/g, "&amp;")
            .replace(/</g, "&lt;")
            .replace(/>/g, "&gt;")
            .replace(/"/g, "&quot;")
            .replace(/'/g, "&#039;");
    }
}

// Инициализация приложения
document.addEventListener('DOMContentLoaded', () => {
    new InterviewApp();
});