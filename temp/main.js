document.addEventListener('DOMContentLoaded', function() {
    const mainPage = document.querySelector('.main-page');
    const createPostPage = document.querySelector('.create-post-page');
    const postFeed = document.querySelector('.post-feed');
    const createPostButton = document.getElementById('create-post-button');
    const backToMainButton = document.getElementById('back-to-main');

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
            body: JSON.stringify({ postId, content }),
        })
        .then(response => response.json())
        .then(comment => {
            const post = document.querySelector(`.post[data-id="${postId}"]`);
            const commentElement = document.createElement('div');
            commentElement.className = 'comment';
            commentElement.textContent = comment.content;
            post.querySelector('.comments').appendChild(commentElement);
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

    document.getElementById('logout-button').addEventListener('click', function() {
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

    // Expose loginSuccess function to be called from regToLog.js
    window.loginSuccess = showMainPage;
});