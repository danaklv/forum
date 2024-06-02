document.addEventListener('DOMContentLoaded', function() {


const wrapper = document.querySelector('.wrapper');
const loginLink = document.querySelector('.login-link');
const registerLink = document.querySelector('.register-link');
const btnLogin = document.querySelector('.loginBtn');
const iconClose = document.querySelector('.icon-close');
const LogButton = document.querySelector('.form-box.login .btn')
const RegButton = document.querySelector('.form-box.register .btn')
const RegForm = document.getElementById('reg-form')
const LoginForm = document.getElementById('login-form')
const ExitButton = document.getElementsByClassName('exitBtn')


registerLink.addEventListener('click', ()=> {
    wrapper.classList.add('active');
});

loginLink.addEventListener('click', ()=> {
    wrapper.classList.remove('active');
});

iconClose.addEventListener('click', ()=> {
    window.location.href = "/";
});


document.addEventListener('DOMContentLoaded', function() {
    console.log("Inline script executed");

    const regForm = document.getElementById("reg-form");
    console.log("Registration form:", regForm);

    if (regForm) {
        regForm.addEventListener("submit", function(event) {
            console.log("Form submitted");

            event.preventDefault();
            event.stopPropagation(); // остановить дальнейшую обработку событий

            // Disable the submit button to prevent multiple submissions
            const submitButton = this.querySelector('button[type="submit"]');
            submitButton.disabled = true;

            var formData = new FormData(this);
            formData.forEach((value, key) => {
                console.log(key + ": " + value);
            });

            var xhr = new XMLHttpRequest();
            xhr.open("POST", "/register", true);

            xhr.onreadystatechange = function() {
                if (xhr.readyState === XMLHttpRequest.DONE) {
                    submitButton.disabled = false; // Re-enable the submit button

                    if (xhr.status === 200) {
                        try {
                            var response = JSON.parse(xhr.responseText);
                            if (response.success) {
                                window.location.href = "/";
                            } else {
                                document.getElementById("err-message").innerText = response.error;
                            }
                        } catch (e) {
                            console.error("Error parsing JSON response: " + e);
                            document.getElementById("err-message").innerText = "An unexpected error occurred";
                        }
                    } else {
                        console.error("Error during request: " + xhr.status);
                        document.getElementById("err-message").innerText = "An error occurred: " + xhr.status;
                    }
                }
            };
            xhr.setRequestHeader("X-Requested-With", "XMLHttpRequest");
            xhr.send(formData);
        }, { once: true });
    }
});


});

document.addEventListener('DOMContentLoaded', function() {
    document.querySelectorAll('.like-button').forEach(button => {
        button.addEventListener('click', function(event) {
            event.preventDefault();
            if (!this.classList.contains('loading')) {
                var postID = this.getAttribute('data-post-id');
                sendLikeRequest(postID, this);
                this.classList.add('loading');
            }
        });
    });

    document.querySelectorAll('.dislike-button').forEach(button => {
        button.addEventListener('click', function(event) {
            event.preventDefault();
            if (!this.classList.contains('loading')) {
                var postID = this.getAttribute('data-post-id');
                sendDislikeRequest(postID, this);
                this.classList.add('loading');
            }
        });
    });
});

function sendLikeRequest(postID, button) {
    var xhr = new XMLHttpRequest();
    xhr.open("POST", "/like", true);
    xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    xhr.onreadystatechange = function() {
        if (xhr.readyState === XMLHttpRequest.DONE) {
            if (xhr.status === 200) {
                var response = JSON.parse(xhr.responseText);
                document.querySelector(`.like-count[data-post-id='${postID}']`).innerText = response.likes;
                button.classList.toggle('liked', !button.classList.contains('liked'));
                document.querySelector(`.dislike-count[data-post-id='${postID}']`).innerText = response.dislikes;
   
                
                if (!button.classList.contains('liked')) {
                    document.querySelector(`.dislike-button[data-post-id='${postID}']`).classList.remove('disliked');
                }
            } else if (xhr.status == 401) {
                window.location.href = "/notauthenticated";
            } else {
                console.error("Error during request: " + xhr.status);
            }
            button.classList.remove('loading');
        }
    };
    xhr.send(`post_id=${postID}`);
}

function sendDislikeRequest(postID, button) {
    var xhr = new XMLHttpRequest();
    xhr.open("POST", "/dislike", true);
    xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
    xhr.onreadystatechange = function() {
        if (xhr.readyState === XMLHttpRequest.DONE) {
            if (xhr.status === 200) {
                var response = JSON.parse(xhr.responseText);
                document.querySelector(`.like-count[data-post-id='${postID}']`).innerText = response.likes;
                document.querySelector(`.dislike-count[data-post-id='${postID}']`).innerText = response.dislikes;
                button.classList.toggle('disliked', !button.classList.contains('disliked'));
                if (!button.classList.contains('disliked')) {
                    document.querySelector(`.like-button[data-post-id='${postID}']`).classList.remove('liked');
                }
            } else if (xhr.status == 401) {
                window.location.href = "/notauthenticated";
            } else {
                console.error("Error during request: " + xhr.status);
            }
            button.classList.remove('loading');
        }
    };
    xhr.send(`post_id=${postID}`);
}