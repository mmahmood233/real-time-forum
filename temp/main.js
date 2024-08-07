document.addEventListener('DOMContentLoaded', function () {
    const authContainer = document.querySelector('.auth-container');
    const mainPage = document.querySelector('.main-page');
    const createPostPage = document.querySelector('.create-post-page');
    const postFeed = document.querySelector('.post-feed');
    const chatArea = document.querySelector('.chat-area');
    const registerForm = document.getElementById('register-form');
    const loginFormElement = document.getElementById('login-form-element');
    const createPostButton = document.getElementById('create-post-button');
    const backToMainButton = document.getElementById('back-to-main');
    const logoutButton = document.getElementById('logout-button');
    const goToLoginButton = document.getElementById('go-to-login');
    const goToRegisterButton = document.getElementById('go-to-register');
    let socket;
    let currentChatUser = null;

    function hideAllSections() {
        authContainer.style.display = 'none';
        mainPage.style.display = 'none';
        createPostPage.style.display = 'none';
    }

    function showRegisterForm() {
        hideAllSections();
        authContainer.style.display = 'block';
        registerForm.style.display = 'block';
        loginFormElement.style.display = 'none';
    }

    function showLoginForm() {
        hideAllSections();
        authContainer.style.display = 'block';
        loginFormElement.style.display = 'block';
        registerForm.style.display = 'none';
    }

    function showMainPage() {
        console.log('Showing main page');
        hideAllSections();
        const mainPage = document.querySelector('.main-page');
        if (mainPage) {
            mainPage.style.display = 'block';
            console.log('Main page displayed');
           
            if (!document.querySelector('.post-feed')) {
                console.error('Post feed element not found, creating it');
                const postFeed = document.createElement('div');
                postFeed.className = 'post-feed';
                mainPage.querySelector('.left-column').appendChild(postFeed);
            }
           
            loadPosts();
            loadChatArea();
            connectWebSocket();
        } else {
            console.error('Main page element not found');
        }
    }

    function showCreatePostPage() {
        hideAllSections();
        createPostPage.style.display = 'block';
    }

    function loadPosts() {
        console.log('Loading posts...');
        fetch('/get-posts')
            .then(response => response.text())
            .then(html => {
                console.log('Received posts HTML:', html);
                const postFeed = document.querySelector('.post-feed');
                if (postFeed) {
                    postFeed.innerHTML = html;
                    console.log('Posts updated');
                    document.querySelectorAll('.post').forEach(post => {
                        const commentsSection = post.querySelector('.comments');
                        const commentForm = post.querySelector('.comment-form');
                        if (commentsSection) {
                            commentsSection.style.display = 'none';
                        }
                        if (commentForm) {
                            commentForm.style.display = 'none';
                        }
                        post.addEventListener('click', handlePostClick);
                    });
                    addLikeDislikeListeners();
                } else {
                    console.error('Post feed element (.post-feed) not found in the DOM');
                    console.log('Current DOM structure:', document.body.innerHTML);
                }
            })
            .catch(error => {
                console.error('Error loading posts:', error);
            });
    }
    
    function handlePostClick(e) {
        if (e.target.classList.contains('like-post') || e.target.classList.contains('dislike-post')) {
            return; // Don't toggle comments when like/dislike buttons are clicked
        }
        const commentsSection = this.querySelector('.comments');
        const commentForm = this.querySelector('.comment-form');
        if (commentsSection && commentForm) {
            const isHidden = commentsSection.style.display === 'none';
            commentsSection.style.display = isHidden ? 'block' : 'none';
            commentForm.style.display = isHidden ? 'block' : 'none';
        }
    }
    

    function addPostEventListeners() {
        document.querySelectorAll('.view-comments').forEach(button => {
            button.addEventListener('click', handleViewComments);
        });
        document.querySelectorAll('.comment-form').forEach(form => {
            form.addEventListener('submit', handleCommentSubmit);
        });
        addLikeDislikeListeners();
    }

    function handleViewComments(e) {
        const postId = e.target.dataset.postId;
        const commentsSection = e.target.closest('.post').querySelector('.comments');
        const commentForm = e.target.closest('.post').querySelector('.comment-form');

        if (commentsSection.style.display === 'none') {
            loadComments(postId, commentsSection);
            commentsSection.style.display = 'block';
            commentForm.style.display = 'block';
        } else {
            commentsSection.style.display = 'none';
            commentForm.style.display = 'none';
        }
    }

    function loadComments(postId, commentsSection) {
        fetch(`/get-comments?postID=${postId}`)
            .then(response => response.text())
            .then(html => {
                commentsSection.innerHTML = html;
                addLikeDislikeListeners();
            })
            .catch(error => console.error('Error loading comments:', error));
    }

    function loadChatArea() {
        console.log('Loading chat area...');
        fetch('/get-chat-area', {
            method: 'GET',
            credentials: 'include'
        })
        .then(response => response.text())
        .then(html => {
            console.log('Received chat HTML:', html);
            const chatArea = document.querySelector('.chat-area');
            if (chatArea) {
                chatArea.innerHTML = html;
                console.log('Chat area updated');
                chatArea.style.display = 'block';
                chatArea.querySelectorAll('li').forEach(userItem => {
                    userItem.addEventListener('click', function() {
                        const userId = this.dataset.userId;
                        console.log('User clicked:', userId);
                        loadMessageHistory(userId);
                    });
                });
            } else {
                console.error('Chat area element not found');
            }
        })
        .catch(error => {
            console.error('Error loading chat area:', error);
        });
    }

    function showChatWindow(userId) {
        const chatWindow = document.querySelector('.chat-window');
        chatWindow.style.display = 'block';
        loadMessageHistory(userId);
    }

    function handleCommentSubmit(e) {
        e.preventDefault();
        const postElement = e.target.closest('.post');
        if (!postElement) {
            console.error('Could not find parent post element');
            return;
        }
        const postId = postElement.dataset.id;
        if (!postId) {
            console.error('Could not find post ID');
            return;
        }
        const commentContent = e.target.comment.value;
        addCommentToPost(postId, commentContent);
        e.target.reset();
    }

    function addCommentToPost(postId, content) {
        const formData = new FormData();
        formData.append('commentCont', content);
   
        fetch(`/add-comment?postID=${postId}`, {
            method: 'POST',
            body: formData,
        })
        .then(response => {
            if (!response.ok) {
                return response.text().then(text => {
                    throw new Error(`Failed to add comment: ${text}`);
                });
            }
            return response.text();
        })
        .then(text => {
            console.log('Server response:', text);
            loadComments(postId, document.querySelector(`.post[data-id="${postId}"] .comments`));
        })
        .catch(error => {
            console.error('Error adding comment:', error);
            alert(`Failed to add comment: ${error.message}`);
        });
    }

    function addLikeDislikeListeners() {
        document.querySelectorAll('.like-post').forEach(button => {
            button.addEventListener('click', handleLikePost);
        });
        document.querySelectorAll('.dislike-post').forEach(button => {
            button.addEventListener('click', handleDislikePost);
        });
        document.querySelectorAll('.like-comment').forEach(button => {
            button.addEventListener('click', handleLikeComment);
        });
        document.querySelectorAll('.dislike-comment').forEach(button => {
            button.addEventListener('click', handleDislikeComment);
        });
    }

    function handleLikePost(e) {
        const postId = e.target.dataset.postId;
        fetch(`/like-post?postID=${postId}`, { method: 'POST' })
            .then(response => response.json())
            .then(data => {
                console.log('Post liked:', data);
                loadPosts();
            })
            .catch(error => console.error('Error:', error));
    }

    function handleDislikePost(e) {
        const postId = e.target.dataset.postId;
        fetch(`/dislike-post?postID=${postId}`, { method: 'POST' })
            .then(response => response.json())
            .then(data => {
                console.log('Post disliked:', data);
                loadPosts();
            })
            .catch(error => console.error('Error:', error));
    }

    function handleLikeComment(e) {
        const commentId = e.target.dataset.commentId;
        fetch(`/like-comment?commentID=${commentId}`, { method: 'POST' })
            .then(response => response.json())
            .then(data => {
                console.log('Comment liked:', data);
                loadPosts();
            })
            .catch(error => console.error('Error:', error));
    }

    function handleDislikeComment(e) {
        const commentId = e.target.dataset.commentId;
        fetch(`/dislike-comment?commentID=${commentId}`, { method: 'POST' })
            .then(response => response.json())
            .then(data => {
                console.log('Comment disliked:', data);
                loadPosts();
            })
            .catch(error => console.error('Error:', error));
    }

    document.getElementById('post-form').addEventListener('submit', function (e) {
        e.preventDefault();
        const content = this.content.value;
        const category = this.category.value;
   
        const formData = new FormData();
        formData.append('postCont', content);
        formData.append('catCont', category);
   
        fetch('/create-post', {
            method: 'POST',
            body: formData,
        })
        .then(response => {
            if (!response.ok) {
                return response.text().then(text => { throw new Error(text) });
            }
            return response.text();
        })
        .then(message => {
            console.log(message);
            this.reset();
            showMainPage();
        })
        .catch(error => {
            console.error('Error:', error);
            alert('Failed to create post: ' + error.message);
        });
    });

    logoutButton.addEventListener('click', function () {
        fetch('/logout', {
            method: 'POST',
            redirect: 'follow'
        })
        .then(response => {
            if (response.ok) {
                window.location.href = response.url;
            } else {
                throw new Error('Logout failed');
            }
        })
        .catch(error => {
            console.error('Error:', error);
            alert('Logout failed. Please try again.');
        });
    });

    function connectWebSocket() {
        socket = new WebSocket('ws://' + window.location.host + '/ws');
        socket.onmessage = function(event) {
            const [senderId, content] = event.data.split(':', 2);
            if (senderId === currentChatUser) {
                appendMessage(content, 'received');
            }
        };
    }
   
    function loadMessageHistory(userId) {
        currentChatUser = userId;
        fetch(`/get-messages?userId=${userId}`)
            .then(response => response.text())
            .then(html => {
                const chatMessages = document.getElementById('chat-messages');
                chatMessages.innerHTML = html;
                chatMessages.scrollTop = chatMessages.scrollHeight;
                document.getElementById('chat-form').style.display = 'flex';
            })
            .catch(error => {
                console.error('Error loading message history:', error);
            });
    }
   
    function appendMessage(content, type) {
        const messageHistory = document.getElementById('chat-messages');
        const messageElement = document.createElement('div');
        messageElement.className = `message ${type}`;
        messageElement.innerHTML = `<span class="content">${content}</span>`;
        messageHistory.appendChild(messageElement);
        messageHistory.scrollTop = messageHistory.scrollHeight;
    }
   
    document.getElementById('chat-form').addEventListener('submit', function(e) {
        e.preventDefault();
        const input = document.getElementById('chat-input');
        if (currentChatUser && input.value.trim()) {
            const message = `${currentChatUser}:${input.value}`;
            socket.send(message);
            appendMessage(input.value, 'sent');
            input.value = '';
        }
    });

    createPostButton.addEventListener('click', showCreatePostPage);
    backToMainButton.addEventListener('click', showMainPage);
    goToLoginButton.addEventListener('click', showLoginForm);
    goToRegisterButton.addEventListener('click', showRegisterForm);

    // Initial page load
    showLoginForm();
    connectWebSocket();

    window.loginSuccess = showMainPage;
});
