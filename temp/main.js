document.addEventListener('DOMContentLoaded', function() {
    const mainPage = document.querySelector('.main-page');
    const createPostPage = document.querySelector('.create-post-page');
    const postFeed = document.querySelector('.post-feed');
    const createPostButton = document.getElementById('create-post-button');
    const backToMainButton = document.getElementById('back-to-main');
    const logoutButton = document.getElementById('logout-button');
    const goToLoginButton = document.getElementById('go-to-login');
    const goToRegisterButton = document.getElementById('go-to-register');

    function showMainPage() {
        document.querySelector('.auth-container').style.display = 'none';
        mainPage.style.display = 'block';
        createPostPage.style.display = 'none';
        loadPosts();
    }

    function showCreatePostPage() {
        mainPage.style.display = 'none';
        createPostPage.style.display = 'block';
    }

    function showRegisterForm() {
        document.getElementById('login-form').style.display = 'none';
        document.getElementById('register-form').style.display = 'block';
    }

    function showLoginForm() {
        document.getElementById('register-form').style.display = 'none';
        document.getElementById('login-form').style.display = 'block';
    }

    function loadPosts() {
        fetch('/get-posts')
        .then(response => response.json())
        .then(posts => {
            postFeed.innerHTML = ''; // Clear existing posts
            posts.forEach(addPostToFeed);
        })
        .catch(error => {
            console.error('Error loading posts:', error);
        });
    }

    function addPostToFeed(post) {
        const postElement = document.createElement('div');
        postElement.className = 'post';
        postElement.dataset.id = post.id; // Add data-id attribute for identification
        postElement.innerHTML = `
            <h3>${post.author}</h3>
            <p>${post.content}</p>
            <small>Category: ${post.category}</small>
            <div class="comments"></div>
            <form class="comment-form">
                <input type="text" name="comment" placeholder="Add a comment" required>
                <button type="submit">Comment</button>
            </form>
        `;
        postElement.querySelector('.comment-form').addEventListener('submit', function(e) {
            e.preventDefault();
            const commentContent = this.comment.value;
            addCommentToPost(post.id, commentContent);
            this.reset();
        });
        postFeed.appendChild(postElement);
    }

function addCommentToPost(postId, content) {
    fetch('/add-comment', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ postID: postId, commentCont: content }),
    })
    .then(response => response.text()) // Read response as text first
    .then(text => {
        console.log('Server response:', text); // Log raw response text
        try {
            const comment = JSON.parse(text); // Try parsing the JSON
            console.log('Parsed comment:', comment); // Log parsed comment object

            // Check if the response has the expected structure
            if (!comment || !comment.CommentContent) {
                console.error('Invalid comment response:', comment);
                return;
            }

            const post = document.querySelector(`.post[data-id="${postId}"]`);
            if (post) {
                const commentElement = document.createElement('div');
                commentElement.className = 'comment';
                commentElement.textContent = comment.CommentContent;
                // Append comment to the post's comments section
                post.querySelector('.comments').appendChild(commentElement);
            } else {
                console.error('Post not found:', postId);
            }
        } catch (error) {
            console.error('Failed to parse JSON response:', error);
        }
    })
    .catch(error => {
        console.error('Error adding comment:', error);
    });
}
 
    
    document.getElementById('post-form').addEventListener('submit', function(e) {
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
            console.log(message); // Log success message
            this.reset();
            showMainPage(); // Go back to main page after creating post
            loadPosts(); // Reload all posts to show the new post
        })
        .catch(error => {
            console.error('Error:', error);
            alert('Failed to create post: ' + error.message);
        });
    });

    logoutButton.addEventListener('click', function() {
        fetch('/logout', { method: 'POST' })
        .then(response => {
            if (response.ok) {
                window.location.reload();
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

    // Expose loginSuccess function to be called from regToLog.js
    window.loginSuccess = showMainPage;
});
