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

    function hideAllSections() {
        authContainer.style.display = 'none';
        mainPage.style.display = 'none';
        createPostPage.style.display = 'none';
        registerForm.style.display = 'none';
        loginFormElement.style.display = 'none';
    }

    function showRegisterForm() {
        hideAllSections();
        authContainer.style.display = 'block';
        registerForm.style.display = 'block';
    }

    function showLoginForm() {
        hideAllSections();
        authContainer.style.display = 'block';
        loginFormElement.style.display = 'block';
    }

    function showMainPage() {
        hideAllSections();
        mainPage.style.display = 'block';
        loadPosts();
        loadChatArea();
    }

    function showCreatePostPage() {
        hideAllSections();
        createPostPage.style.display = 'block';
    }

    function loadPosts() {
        fetch('/get-posts')
            .then(response => response.text())
            .then(html => {
                postFeed.innerHTML = html;
                document.querySelectorAll('.comment-form').forEach(form => {
                    form.addEventListener('submit', handleCommentSubmit);
                });
                addLikeDislikeListeners();
            })
            .catch(error => {
                console.error('Error loading posts:', error);
            });
    }

    function loadChatArea() {
        fetch('/get-chat-area', {
            method: 'GET',
            credentials: 'include' // This ensures the session cookie is sent with the request
        })
        .then(response => response.text())
        .then(html => {
            const chatArea = document.querySelector('.chat-area');
            if (chatArea) {
                chatArea.innerHTML = html;
            } else {
                console.error('Chat area element not found');
            }
        })
        .catch(error => {
            console.error('Error loading chat area:', error);
        });
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
            loadPosts();
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

    createPostButton.addEventListener('click', showCreatePostPage);
    backToMainButton.addEventListener('click', showMainPage);
    goToLoginButton.addEventListener('click', showLoginForm);
    goToRegisterButton.addEventListener('click', showRegisterForm);

    setInterval(loadChatArea, 10); // Refresh chat area every 30 seconds

    // Initial page load
    showLoginForm();

    window.loginSuccess = showMainPage;
});
