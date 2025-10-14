// frontend/js/app.js
class InterviewApp {
    constructor() {
        this.questions = [];
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

        // Закрытие модального окна при клике вне его
        window.addEventListener('click', (e) => {
            const modal = document.getElementById('questionModal');
            if (e.target === modal) {
                this.closeModal();
            }
        });
    }

    async loadQuestions() {
        try {
            const response = await fetch('/api/v1/questions');
            this.questions = await response.json();
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
        div.innerHTML = `
            <div class="question-header">
                <div class="question-title">${this.escapeHtml(question.title)}</div>
                <div class="tags">
                    ${question.tags.map(tag => 
                        `<span class="tag tag-${tag.toLowerCase()}">${tag}</span>`
                    ).join('')}
                </div>
            </div>
            <div class="question-answer">
                <div class="answer-content">${this.escapeHtml(question.answer)}</div>
            </div>
        `;

        // Обработчик клика для раскрытия/скрытия ответа
        div.addEventListener('click', (e) => {
            if (!e.target.classList.contains('tag')) {
                const answer = div.querySelector('.question-answer');
                answer.classList.toggle('expanded');
            }
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

    openModal() {
        document.getElementById('questionModal').style.display = 'block';
    }

    closeModal() {
        document.getElementById('questionModal').style.display = 'none';
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