document.addEventListener("DOMContentLoaded", function() {
    const token = getTokenFromCookie();
    if (token) {
        showLoggedInNav();
         LoadPosts();
    } else {
        showLoggedOutNav();
    }
    // Add event listeners to the navigation links
    document.getElementById("nav-login").addEventListener("click", function(event) {
        event.preventDefault();
        showSection("login-section");
    });
    document.getElementById("nav-register").addEventListener("click", function(event) {
        event.preventDefault();
        showSection("register-section");
    });
    document.getElementById("nav-create-posts").addEventListener("click", function(event) {
        event.preventDefault();
        showSection("create-post-section");
    });
    document.getElementById("nav-posts").addEventListener("click", function(event) {
        event.preventDefault();
        showSection("posts-section");
        LoadPosts();
    });
    document.getElementById("nav-private-messages").addEventListener("click", function(event) {
        event.preventDefault();
        showSection("private-messages-section");
    });
    document.getElementById("nav-logout").addEventListener("click", function(event) {
        event.preventDefault();
        logoutUser();
    });
    const googleLoginLink = document.getElementById("google-login");
    if (googleLoginLink) {
        googleLoginLink.addEventListener("click", function(event) {
            event.preventDefault();
            window.location.href = "/auth/google/login-form";
        });
    }
    const githubLoginLink = document.getElementById("github-login");
    if (githubLoginLink) {
        githubLoginLink.addEventListener("click", function(event) {
            event.preventDefault();
            window.location.href = "/auth/github/login-form";
        });
    }
    
    document.getElementById("login-form").addEventListener("submit", function(event) {
        event.preventDefault();
        const identifier = document.getElementById("identifier").value;
        const password = document.getElementById("password").value;
        const errorMessageElement = document.getElementById("login-error");
        console.log(identifier, password)
        if (validateLoginForm(identifier, password, errorMessageElement)) {
            login(identifier, password, errorMessageElement);
        }
    });
    const formulaires = document.getElementById("register-form");
    if (formulaires) {
        formulaires.addEventListener("submit", function(event){
            event.preventDefault();
            const formData = new FormData(formulaires);
            const form = {
                username: formData.get("username"),
                age: parseInt(formData.get("age")),
                gender: formData.get("gender"),
                firstName: formData.get("firstName"),
                lastName: formData.get("lastName"),
                email: formData.get("email"),
                password: formData.get("password")
            };
            console.log("form", form.gender)
            const errorMessageElement = document.getElementById("register-error");
            if (validateRegisterForm(form, errorMessageElement)) {
                register(form, errorMessageElement);
            }
        })
    }
    const createPostForm = document.getElementById("create-post-form");
    if (createPostForm) {
        createPostForm.addEventListener("submit", function(event) {
            event.preventDefault();
            const formData = new FormData(createPostForm);
            const errorMessageElement = document.getElementById("create-post-error");
            const post = {
                title: formData.get("title"),
                category: formData.get("category"),
                content: formData.get("content"),
                image_path: formData.get("image").name
            };
            if (validateCreatePostForm(post, errorMessageElement)) {
                createPost(post, errorMessageElement);
            }
        });
    }

    function getTokenFromCookie() {
        const cookies = document.cookie.split(";").map(cookie => cookie.trim());
        for (const cookie of cookies) {
            if (cookie.startsWith("token=")) {
                return cookie.substring("token=".length);
            }
        }
        return null;
    }
    function showLoggedOutNav() {
        document.querySelectorAll("nav ul li a").forEach(link => {
            const id = link.getAttribute("id");
            if (id === "nav-login" || id === "nav-register") {
                link.classList.remove("hidden");
            } else {
                link.classList.add("hidden");
            }
        });
    }
    
    function showLoggedInNav() {
        document.querySelectorAll("nav ul li a").forEach(link => {
            const id = link.getAttribute("id");
            if (id !== "nav-login" && id !== "nav-register") {
                link.classList.remove("hidden");
            } else {
                link.classList.add("hidden");
            }
        });
    }
    
    function showSection(sectionId) {
        document.querySelectorAll("section").forEach(section => {
            if (section.id === sectionId) {
                section.classList.remove("hidden");
            } else {
                section.classList.add("hidden");
            }
        });
    }
    function displayErrorMessage(message, element) {
        element.textContent = message;
        element.classList.remove("hidden");
    }
    function validateEmail(email) {
        const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        return re.test(String(email).toLowerCase());
    }
    function validateLoginForm(identifier, password, errorMessageElement) {
        if (identifier.length === 0 || password.length === 0) {
            displayErrorMessage("Please complete all fields", errorMessageElement);
            return false;
        }
        if (identifier.length < 6 || identifier.length > 16 || password.length < 6 || password.length > 16) {
            displayErrorMessage("Identifier and password must be between 6 and 16 characters", errorMessageElement);
            return false;
        }
        return true;
    }
    function validateRegisterForm(formData, errorMessageElement) {
        if (formData.username.length === 0) {
            displayErrorMessage("Le nom d'utilisateur est requis.", errorMessageElement);
            return false;
        }
        if (formData.username.length < 4 || formData.username.length > 16) {
            displayErrorMessage("Le nom d'utilisateur doit contenir entre 4 et 16 caractères.", errorMessageElement);
            return false;
        }
    
        if (isNaN(formData.age) || formData.age <= 0) {
            displayErrorMessage("L'âge doit être un nombre supérieur à zéro.", errorMessageElement);
            return false;
        }
    
        if (!["homme", "femme"].includes(formData.gender.toLowerCase())) {
            displayErrorMessage("Le genre doit être 'Homme' ou 'Femme'.", errorMessageElement);
            return false;
        }
    
        if (formData.firstName.length === 0) {
            displayErrorMessage("Le prénom est requis.", errorMessageElement);
            return false;
        }
        if (formData.firstName.length > 16) {
            displayErrorMessage("Le prénom ne doit pas dépasser 16 caractères.", errorMessageElement);
            return false;
        }
    
        if (formData.lastName.length === 0) {
            displayErrorMessage("Le nom de famille est requis.", errorMessageElement);
            return false;
        }
        if (formData.lastName.length > 16) {
            displayErrorMessage("Le nom de famille ne doit pas dépasser 16 caractères.", errorMessageElement);
            return false;
        }
    
        if (!validateEmail(formData.email)) {
            displayErrorMessage("Veuillez entrer un email valide.", errorMessageElement);
            return false;
        }
    
        if (formData.password.length < 6 || formData.password.length > 16) {
            displayErrorMessage("Le mot de passe doit contenir entre 6 et 16 caractères.", errorMessageElement);
            return false;
        }
    
        return true;
    }
    
    function validateCreatePostForm(post, errorMessageElement) {
        let errorMessage = "";
        if (!post.title || post.title.trim().length === 0) {
            errorMessage += "Title is required. ";
        }
        if (!post.content || post.content.trim().length === 0) {
            errorMessage += "Content is required. ";
        }
        if (!post.category || post.category.trim().length === 0) {
            errorMessage += "Category is required. ";
        } else {
            const validCategories = ["All", "HTML/CSS", "JavaScript", "Java", "C#", "C++", "Python", "PHP", "Ruby"];
            if (!validCategories.includes(post.category)) {
                errorMessage += "Invalid category. ";
            }
        }
        if (errorMessage) {
            displayErrorMessage(errorMessage, errorMessageElement);
            return false;
        }
        return true;
    }
    function register(formData, errorMessageElement) {
        fetch("/register-form", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(formData)
        })
        .then(response => {
            if (!response.ok) {
                return response.text().then(text => { throw new Error(text) });
            }
            return response.json();
        })
        .then(data => {
            console.log("Registration successful:", data);
            showSection("login-section");
        })
        .catch(error => {
            console.error("Register failed", error);
            displayErrorMessage("Invalid registration credentials", errorMessageElement);
        });
    }
    function login(identifier, password, errorMessageElement) {
        fetch("/login-form", {
            method: "POST",
            headers: {
                "Content-Type": "application/x-www-form-urlencoded",
            },
            body: `identifier=${encodeURIComponent(identifier)}&password=${encodeURIComponent(password)}`
        })
        .then(response => {
            if (!response.ok) {
                return response.json().then(errData => {
                    throw new Error(errData.message || "Login failed");
                });
            }
            return response.json();
        })
        .then(data => {
            console.log("Login successful:", data);
            document.cookie = `token=${data.token}; expires=${new Date(Date.now() + 24 * 60 * 60 * 1000).toUTCString()}; path=/; Secure; HttpOnly`;
            showLoggedInNav();
            showSection("posts-section");
            LoadPosts();
        })
        .catch(error => {
           console.log("Login failed:", error);
            displayErrorMessage("Incorrect password or username", errorMessageElement);
        });
    }

    function displayPosts(posts) {
        const postsSection = document.getElementById("posts-section");
        postsSection.innerHTML = "";
        postsSection.classList.remove("hidden");
    
        if (posts.length === 0) {
            postsSection.innerHTML = "<p>No posts available</p>";
        } else {
            posts.forEach(post => {
                const postElement = document.createElement("div");
                postElement.classList.add("postClass");
    
                postElement.innerHTML = `
                    <div class="post-header">
                        <h3>${post.title}</h3>
                    </div>
                    <p><strong>Category:</strong> ${post.category}</p>
                    <p>${post.content}</p>
                    <p><strong>By:</strong> ${post.username} <strong>On:</strong> ${new Date(post.created_at).toLocaleDateString()}</p>
                    ${post.image_path ? `<img src="${post.image_path}" alt="Post Image" class="post-image">` : ""}
    
                    <div class="post-actions">
                        <button class="likeButton" data-post-id="${post.id}">Like</button>
                        <button class="unlikeButton" data-post-id="${post.id}" disabled>Unlike</button>
                        <span class="likeCount" id="like-count-${post.id}">0 Likes</span>
                    </div>
    
                    <button class="show-comments-button" data-post-id="${post.id}">Show Comments</button>
                    <div class="comments hidden" id="comments-container-${post.id}"></div>
                `;
    
                postsSection.appendChild(postElement);
    
                const likeButton = postElement.querySelector('.likeButton');
                const unlikeButton = postElement.querySelector('.unlikeButton');
                const likeCount = postElement.querySelector(`#like-count-${post.id}`);
    
                let likes = 0;
    
                likeButton.addEventListener('click', function() {
                    likes++;
                    likeCount.textContent = `${likes} Likes`;
                    likeButton.disabled = true;
                    unlikeButton.disabled = false;
                });
    
                unlikeButton.addEventListener('click', function() {
                    likes--;
                    likeCount.textContent = `${likes} Likes`;
                    likeButton.disabled = false;
                    unlikeButton.disabled = true;
                });
    
                const showCommentButton = postElement.querySelector(".show-comments-button");
                showCommentButton.addEventListener("click", function () {
                    const commentsContainer = postElement.querySelector(`#comments-container-${post.id}`);
                    if (commentsContainer.classList.contains('hidden')) {
                        LoadComments(post.id);
                        commentsContainer.classList.remove('hidden');
                        showCommentButton.textContent = "Hide Comments";
                    } else {
                        commentsContainer.classList.add('hidden');
                        showCommentButton.textContent = "Show Comments";
                    }
                });
            });
        }
    }

    function displayComments(postID, comments) {
        const commentsContainer = document.getElementById(`comments-container-${postID}`);
        commentsContainer.innerHTML = ""; 
    
        if (!comments.length) {
            commentsContainer.innerHTML = "<p>No comments available</p>";
        } else if (Array.isArray(comments)) {
            comments.forEach(comment => {
                const commentElement = document.createElement("div");
                commentElement.classList.add("comment");
                commentElement.innerHTML = `
                    <p>${comment.content}</p>
                    <p><strong>By:</strong> ${comment.username_post} <strong>On:</strong> ${new Date(comment.created_at).toLocaleDateString()}</p>
                `;
                commentsContainer.appendChild(commentElement);
            });
    
            const commentForm = document.createElement("div");
            commentForm.innerHTML = `
                <textarea id="comment-text-${postID}" placeholder="Write a comment..."></textarea>
                <button class="submit-comment" data-post-id="${postID}">Submit Comment</button>
                <p class="error-message" id="create-comment-error-${postID}" style="color:red;"></p>
            `;
            commentsContainer.appendChild(commentForm);
    
            const commentButton = commentForm.querySelector('.submit-comment');
            commentButton.addEventListener('click', function() {
                const commentText = commentForm.querySelector(`#comment-text-${postID}`).value.trim();
                const errorMessageElement = commentForm.querySelector(`#create-comment-error-${postID}`);
    
                if (commentText) {
                    AddComment(postID, commentText, errorMessageElement);
                }
            });
        } else {
            console.error('Expected comments to be an array:', comments);
        }
    }
    async function LoadPosts(){
        try {
            const response = await fetch('/list-posts-form');
            if (!response.ok){
                throw new Error(`Failed to fetch posts: ${response.status} ${response.statusText}`);
            }
            const text = await response.text();
            let posts;
            try {
                posts = JSON.parse(text);
            } catch (error) {
                throw new Error('Failed to parse JSON: ' + error.message);
            }
            displayPosts(posts);
        } catch (error) {
            console.error('Error loading posts:', error);
        }
    }

    async function LoadComments(postID) {
        try {
            const response = await fetch(`/list-comment-form?post_id=${postID}`);
            if (!response.ok) {
                throw new Error(`Failed to fetch comments: ${response.status}`);
            }
            const data = await response.json();
            displayComments(postID, data.comments);
        } catch (error) {
            console.error('Error loading comments:', error);
        }
    }
    
    function AddComment(postID, commentText, errorMessageElement) {
        const comment = {
            post_id: parseInt(postID, 10), 
            content: commentText
        };
    
        fetch("/create-comment-form", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(comment) 
        })
        .then(response => {
            if (!response.ok) {
                return response.text().then(text => { throw new Error(text) });
            }
            return response.json();
        })
        .then(newComment => {
            document.getElementById(`comment-text-${postID}`).value = ""; 
            if (errorMessageElement) errorMessageElement.classList.add("hidden");
    
            console.log("Create comment successful:", newComment);
            LoadComments(postID); 
        })
        .catch(error => {
            console.error("Create comment failed:", error);
            displayErrorMessage("Failed to create comment", errorMessageElement);
        });
    }
    
      
    function createPost(post, errorMessageElement) {
        fetch("/create-post-form", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(post)
        })
        .then(response => {
            if (!response.ok) {
                return response.text().then(text => { throw new Error(text) });
            }
            return response.json();
        })
        .then(data => {
            console.log("Create post successful:", data);
            document.getElementById("create-post-form").reset();
            errorMessageElement.classList.add("hidden");
            showSection("posts-section");
        })
        .catch(error => {
            console.error("Create post failed:", error);
            displayErrorMessage("Failed to create post", errorMessageElement);
        });
    }



    function logoutUser() {
        fetch("/logout", {
            method: "POST"
        })
        .then(response => {
            if (response.ok) {
                document.cookie = "token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
                showLoggedOutNav();
                showSection("login-section");
            } else {
                alert('Logout failed. Please try again.');
            }
        })
        .catch(error => {
            console.error('Error logging out:', error);
        });
    }
    
});