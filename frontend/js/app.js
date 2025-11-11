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
        document.getElementById('addQuestionBtn').addEventListener('click', () => {
            this.openModal();
        });

        document.querySelectorAll('.close').forEach(closeBtn => {
            closeBtn.addEventListener('click', (e) => {
                const modal = e.target.closest('.modal');
                if (modal.id === 'questionModal') {
                    this.closeModal();
                } else if (modal.id === 'editQuestionModal') {
                    this.closeEditModal();
                }
            });
        });

        document.getElementById('questionForm').addEventListener('submit', (e) => {
            e.preventDefault();
            this.createQuestion();
        });

        document.getElementById('editQuestionForm').addEventListener('submit', (e) => {
            e.preventDefault();
            this.updateQuestion();
        });

        document.getElementById('questionTags').addEventListener('input', (e) => {
            this.showTagSuggestions(e.target.value);
        });

        document.getElementById('editQuestionTags').addEventListener('input', (e) => {
            this.showEditTagSuggestions(e.target.value);
        });

        window.addEventListener('click', (e) => {
            if (e.target.id === 'questionModal') {
                this.closeModal();
            } else if (e.target.id === 'editQuestionModal') {
                this.closeEditModal();
            }
        });
    }

    showTagSuggestions(inputValue) {
        this.hideTagSuggestions();
        
        if (!inputValue.trim()) {
            this.renderAllAvailableTags('tagsSuggestions', 'questionTags');
            return;
        }

        const inputTags = inputValue.split(',').map(tag => tag.trim().toLowerCase());
        const currentTag = inputTags[inputTags.length - 1];
        
        if (!currentTag) {
            this.renderAllAvailableTags('tagsSuggestions', 'questionTags');
            return;
        }

        const filteredTags = this.availableTags.filter(tag => 
            tag.toLowerCase().includes(currentTag.toLowerCase()) && 
            !inputTags.includes(tag)
        );

        this.renderTagSuggestions(filteredTags, 'tagsSuggestions', 'questionTags');
    }

    showEditTagSuggestions(inputValue) {
        this.hideEditTagSuggestions();
        
        if (!inputValue.trim()) {
            this.renderAllAvailableTags('editTagsSuggestions', 'editQuestionTags');
            return;
        }

        const inputTags = inputValue.split(',').map(tag => tag.trim().toLowerCase());
        const currentTag = inputTags[inputTags.length - 1];
        
        if (!currentTag) {
            this.renderAllAvailableTags('editTagsSuggestions', 'editQuestionTags');
            return;
        }

        const filteredTags = this.availableTags.filter(tag => 
            tag.toLowerCase().includes(currentTag.toLowerCase()) && 
            !inputTags.includes(tag)
        );

        this.renderTagSuggestions(filteredTags, 'editTagsSuggestions', 'editQuestionTags');
    }

    hideTagSuggestions() {
        const existingSuggestions = document.getElementById('tagsSuggestions');
        if (existingSuggestions) {
            existingSuggestions.remove();
        }
    }

    hideEditTagSuggestions() {
        const existingSuggestions = document.getElementById('editTagsSuggestions');
        if (existingSuggestions) {
            existingSuggestions.remove();
        }
    }

    renderAllAvailableTags(suggestionsId, tagsInputId) {
        const tagsInput = document.getElementById(tagsInputId);
        const suggestionsDiv = document.createElement('div');
        suggestionsDiv.id = suggestionsId;
        suggestionsDiv.className = 'tags-suggestions';
        
        suggestionsDiv.innerHTML = `
            <div class="suggestions-title">Возможные теги:</div>
            <div class="available-tags">
                ${this.availableTags.map(tag => 
                    `<span class="available-tag tag-${tag}" onclick="interviewApp.addTagToInput('${tag}', '${tagsInputId}')">${tag}</span>`
                ).join('')}
            </div>
        `;
        
        tagsInput.parentNode.insertBefore(suggestionsDiv, tagsInput.nextSibling);
    }

    renderTagSuggestions(tags, suggestionsId, tagsInputId) {
        if (tags.length === 0) return;
        
        const tagsInput = document.getElementById(tagsInputId);
        const suggestionsDiv = document.createElement('div');
        suggestionsDiv.id = suggestionsId;
        suggestionsDiv.className = 'tags-suggestions';
        
        suggestionsDiv.innerHTML = `
            <div class="suggestions-title">Возможные теги:</div>
            <div class="available-tags">
                ${tags.map(tag => 
                    `<span class="available-tag tag-${tag}" onclick="interviewApp.addTagToInput('${tag}', '${tagsInputId}')">${tag}</span>`
                ).join('')}
            </div>
        `;
        
        tagsInput.parentNode.insertBefore(suggestionsDiv, tagsInput.nextSibling);
    }

    addTagToInput(tag, inputId) {
        const tagsInput = document.getElementById(inputId);
        const currentTags = tagsInput.value.split(',').map(t => t.trim()).filter(t => t);
        
        if (currentTags.length > 0) {
            currentTags.pop();
        }
        
        currentTags.push(tag);
        tagsInput.value = currentTags.join(', ');
        tagsInput.focus();
        
        if (inputId === 'questionTags') {
            this.hideTagSuggestions();
            setTimeout(() => this.showTagSuggestions(tagsInput.value), 100);
        } else {
            this.hideEditTagSuggestions();
            setTimeout(() => this.showEditTagSuggestions(tagsInput.value), 100);
        }
    }

    async loadQuestions() {
        try {
            const response = await fetch('/api/v1/questions');
            const data = await response.json();
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
        const title = question.title || 'Без названия';
        const answer = question.answer || 'Нет ответа';
        const tags = question.tags || [];
        const formattedAnswer = this.formatAnswerText(answer);
            
        div.innerHTML = `
            <div class="question-header">
                <div class="question-title">${this.escapeHtml(title)}</div>
                <div class="question-actions">
                    <div class="tags">
                        ${tags.map(tag => 
                            `<span class="tag tag-${String(tag).toLowerCase()}">${this.escapeHtml(tag)}</span>`
                        ).join('')}
                    </div>
                    <button class="edit-btn" data-id="${id}">✏️</button>
                    <button class="delete-btn" data-id="${id}">🗑️</button>
                </div>
            </div>
            <div class="question-answer">
                <div class="answer-content">${formattedAnswer}</div>
            </div>
        `;

        div.addEventListener('click', (e) => {
            if (!e.target.classList.contains('tag') && 
                !e.target.classList.contains('delete-btn') &&
                !e.target.classList.contains('edit-btn')) {
                const answer = div.querySelector('.question-answer');
                answer.classList.toggle('expanded');
            }
        });

        div.querySelector('.edit-btn').addEventListener('click', (e) => {
            e.stopPropagation();
            this.openEditModal(question);
        });

        div.querySelector('.delete-btn').addEventListener('click', (e) => {
            e.stopPropagation();
            this.deleteQuestion(id);
        });

        return div;
    }

    formatAnswerText(text) {
        if (!text) return 'Нет ответа';
        const safeText = this.escapeHtml(text);
        return safeText.replace(/\n/g, '<br>');
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
                this.loadQuestions();
                document.getElementById('questionForm').reset();
            } else {
                alert('Ошибка при создании вопроса');
            }
        } catch (error) {
            console.error('Ошибка:', error);
            alert('Ошибка при создании вопроса');
        }
    }

    openEditModal(question) {
        document.getElementById('editQuestionId').value = question.id;
        document.getElementById('editQuestionTitle').value = question.title;
        document.getElementById('editQuestionAnswer').value = question.answer;
        document.getElementById('editQuestionTags').value = question.tags ? question.tags.join(', ') : '';
        document.getElementById('editQuestionModal').style.display = 'block';
        setTimeout(() => this.showEditTagSuggestions(''), 100);
    }

    closeEditModal() {
        document.getElementById('editQuestionModal').style.display = 'none';
        this.hideEditTagSuggestions();
    }

    async updateQuestion() {
        const id = document.getElementById('editQuestionId').value;
        const title = document.getElementById('editQuestionTitle').value;
        const answer = document.getElementById('editQuestionAnswer').value;
        const tags = document.getElementById('editQuestionTags').value
            .split(',')
            .map(tag => tag.trim())
            .filter(tag => tag);

        try {
            const response = await fetch(`/api/v1/questions/up/${id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ title, answer, tags }),
            });

            if (response.ok) {
                this.closeEditModal();
                this.loadQuestions();
            } else {
                const error = await response.json();
                alert(`Ошибка при обновлении: ${error.error}`);
            }
        } catch (error) {
            console.error('Ошибка:', error);
            alert('Ошибка при обновлении вопроса');
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

document.addEventListener('DOMContentLoaded', () => {
    window.interviewApp = new InterviewApp();
});