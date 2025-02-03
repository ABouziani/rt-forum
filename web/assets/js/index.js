window.addEventListener('resize', () => {
    if (document.body.clientWidth > 600 && document.querySelector('.mobile-nav')) {
        document.querySelector('.mobile-nav').style.display = 'none';
    }
})


function throttle(fn, delay) {
    let last = 0;
    return function () {
        const now = +new Date();
        if (now - last > delay) {
            fn.apply(null, arguments);
            last = now;
        }
    };
}

const addcomment = throttle(addcomm, 5000)

let cats = ["", "Technology", "Health", "Travel", "Education", "Entertainment"]

function postreaction(postId, reaction) {
    const logerror = document.getElementById("errorlogin" + postId);
    logerror.innerText = ``;

    fetch("/post/postreaction", {
        method: "POST",
        headers: {
            "Content-Type": "application/x-www-form-urlencoded",
        },
        body: `reaction=${reaction}&post_id=${postId}`
    })
        .then(response => {
            if (response.ok) {
                return response.json();
            } else {
                if (response.status === 401) {
                    refetchLogin('/login');
                } else if (response.status === 400) {
                    writeError(logerror, "red", `Bad request!`, 1000);
                } else if (response.status === 500) {
                    writeError(logerror, "red", `Try again later!`, 1000);
                }
            }
        })
        .then(data => {
            if (data) {
                document.getElementById("likescount" + postId).innerHTML = `<i class="fa-regular fa-thumbs-up"></i>${data.likesCount}`;
                document.getElementById("dislikescount" + postId).innerHTML = `<i class="fa-regular fa-thumbs-down"></i>${data.dislikesCount}`;
            }
        })
        .catch(error => {
            console.error(error);
        });
}

function commentreaction(commentid, reaction) {
    const logerror = document.getElementById("commenterrorlogin" + commentid);
    logerror.innerText = ``;

    fetch("/post/commentreaction", {
        method: "POST",
        headers: {
            "Content-Type": "application/x-www-form-urlencoded",
        },
        body: `reaction=${reaction}&comment_id=${commentid}`
    })
        .then(response => {
            if (response.ok) {
                return response.json();
            } else {
                if (response.status === 401) {
                    refetchLogin('/login');
                } else if (response.status === 400) {
                    writeError(logerror, "red", `bad request!`, 1000);
                } else if (response.status === 500) {
                    writeError(logerror, "red", `Try again later!`, 1000);
                }
            }
        })
        .then(data => {
            if (data) {
                document.getElementById("commentlikescount" + commentid).innerHTML = `<i class="fa-regular fa-thumbs-up"></i>${data.commentlikesCount}`;
                document.getElementById("commentdislikescount" + commentid).innerHTML = `<i class="fa-regular fa-thumbs-down"></i>${data.commentdislikesCount}`;
            }
        })
        .catch(error => {
            console.error(error);
        });
}


function addcomm(postId) {
    const content = document.getElementById("comment-content");
    const logerror = document.getElementById("errorlogin" + postId);

    if (!content.value) {
        writeError(logerror, "red", 'Please fill in Comment field.', 3000);
        return;
    }

    if (content.value.length > 500) {
        writeError(logerror, "red", 'Comment is too long. Please keep it under 500 characters.', 3000);
        return;
    }

    fetch("/post/addcommentREQ", {
        method: "POST",
        headers: {
            "Content-Type": "application/x-www-form-urlencoded",
        },
        body: `postid=${postId}&comment=${encodeURIComponent(content.value)}`
    })
        .then(response => {
            if (response.ok) {
                return response.json();
            } else if (response.status === 400) {
                writeError(logerror, "red", `Invalid comment!`, 1000);
            } else if (response.status === 401) {
                refetchLogin('/login');
            } else {
                writeError(logerror, "red", `Cannot add comment now, try again later!`, 1000);
            }
        })
        .then(response => {
            if (response) {
                const comment = document.createElement("div");
                comment.innerHTML = `
                <div class="comment">
                    <div class="comment-header">
                        <p class="comment-user">${response.username}</p>
                        <span></span>
                        <p class="comment-time">${response.created_at}</p>
                    </div>
                    <div class="comment-body">
                        <p class="comment-content">${response.content}</p>
                    </div>
                    <div class="comment-footer">
                        <button id="commentlikescount${response.ID}" onclick="commentreaction('${response.ID}','like')"
                            class="comment-like"><i class="fa-regular fa-thumbs-up"></i>${response.likes}</button>
                        <button id="commentdislikescount${response.ID}" onclick="commentreaction('${response.ID}','dislike')"
                            class="comment-dislike"><i class="fa-regular fa-thumbs-down"></i>${response.dislikes}</button>
                    </div>
                    <span style="color:red" id="commenterrorlogin${response.ID}"></span>
                </div>
            `;
                document.getElementsByClassName("comments")[0].prepend(comment);
                document.getElementsByClassName("post-comments")[0].innerHTML = `<i class="fa-regular fa-comment"></i>${response.commentscount}`;
                content.value = "";
            }
        })
        .catch(error => {
            console.error(error);
        });
}


function CreatPost() {
    const title = document.querySelector(".create-post-title");
    const content = document.querySelector(".content");
    const categories = document.querySelector(".selected-categories");
    const logerror = document.querySelector(".errorarea");

    if (title.value.trim() === "" || content.value.trim() === "" || categories.childElementCount === 0) {
        writeError(logerror, "red", 'No empty entries allowed!', 3000);
        return;
    }

    if (title.value.length > 100) {
        writeError(logerror, "red", 'Title is too long. Please keep it under 100 characters.', 3000);
        return;
    }

    if (content.value.length > 3000) {
        writeError(logerror, "red", 'Content is too long. Please keep it under 3000 characters.', 3000);
        return;
    }

    let cateris = [];
    Array.from(categories.getElementsByTagName('input')).forEach((x) => {
        cateris.push(x.value);
    });

    fetch("/post/createpost", {
        method: "POST",
        headers: {
            "Content-Type": "application/x-www-form-urlencoded",
        },
        body: `title=${encodeURIComponent(title.value)}&content=${encodeURIComponent(content.value)}&categories=${cateris}`
    })
        .then(response => {
            if (response.ok) {
                const btn = document.getElementById("create-post-btn");
                document.getElementById("publish-post-icon").style.display = "none";
                document.getElementById("publish-post-circle").style.display = "inline-block";
                btn.disabled = true;
                btn.style.background = "grey";
                btn.style.cursor = "not-allowed";

                writeError(logerror, "green", 'Post created successfully, redirect to home page in 2s ...', 2000);
                setTimeout(() => {
                    refetch('/')
                }, 2000);
            } else if (response.status === 401) {
                refetchLogin("/login")
            } else if (response.status === 400) {
                writeError(logerror, "red", 'Bad request!', 1500);
            } else {
                writeError(logerror, "red", 'Error: check your entries and try again!', 1500);
            }
        })
        .catch(error => {
            console.error('Error:', error);
        });
}

function register() {
    const email = document.querySelector("#email");
    const username = document.querySelector("#username");
    const password = document.querySelector("#password");
    const passConfirm = document.querySelector("#password-confirmation");
    const logerror = document.querySelector(".errorarea");

    if (username.value.length < 4 || username.value.includes(" ")) {
        writeError(logerror, "red", "Username too short! or has space", 1500);
        return;
    }

    if (password.value.length < 6) {
        writeError(logerror, "red", "Password too short!", 1500);
        return;
    }

    if (password.value !== passConfirm.value) {
        writeError(logerror, "red", "Password and password confirmation are not identical", 1500);
        return;
    }

    // Prepare the data to be sent
    const formData = new URLSearchParams();
    formData.append('email', email.value);
    formData.append('username', username.value);
    formData.append('password', password.value);
    formData.append('password-confirmation', passConfirm.value);

    // Send the request using fetch
    fetch('/signup', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: formData.toString(),
    })
        .then(response => {
            if (response.status === 200) {
                writeError(logerror, "green", `User ${username.value} created successfully, redirecting to login page in 2s ...`, 2000);
                setTimeout(async () => {
                    refetchLogin('/login')
                }, 2000);
            } else if (response.status === 302) {
                writeError(logerror, "green", 'You are already logged in, redirecting to home page in 2s...', 2000);
                setTimeout(async () => {
                    document.documentElement.innerHTML = await response.text()
                }, 2000);
            } else if (response.status === 400) {
                writeError(logerror, "red", 'Error: Verify your data and try again!', 1500);
            } else if (response.status === 304) {
                writeError(logerror, "red", 'User already exists!', 1500);
            } else {
                writeError(logerror, "red", 'Cannot create user, try again later!', 1500);
            }
        })
        .catch(() => {
            writeError(logerror, "red", 'Network error. Please try again later!', 1500);
        });
}




function login() {
    const username = document.querySelector("#username");
    const password = document.querySelector("#password");
    const logerror = document.querySelector(".errorarea");

    if (username.value.length < 4) {
        writeError(logerror, "red", "Username too short!", 1500);
        return;
    }
    if (password.value.length < 6) {
        writeError(logerror, "red", "Password too short!", 1500);
        return;
    }

    const formData = new URLSearchParams();
    formData.append('username', username.value);
    formData.append('password', password.value);

    fetch('/signin', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: formData.toString(),
    })
        .then(async response => {
            if (response.status === 200) {


                writeError(logerror, "green", `Login in successfully, redirect to home page in 2s ...`, 2000);

                document.documentElement.innerHTML = await response.text()

            } else if (response.status === 302) {
                writeError(logerror, "green", 'You are already logged in, redirect to home page in 2s...', 2000);
                setTimeout(async () => {
                    document.documentElement.innerHTML = await response.text()
                }, 2000);
            } else if (response.status === 400) {
                writeError(logerror, "red", 'Error: verify your data and try again!', 1500);
            } else if (response.status === 404) {
                writeError(logerror, "red", 'User not found!', 1500);
            } else if (response.status === 401) {
                writeError(logerror, "red", 'Invalid username or password!', 1500);
            } else {
                writeError(logerror, "red", 'Cannot log you in now, try again later!', 1500);
            }
        })
        .catch(() => {
            writeError(logerror, "red", 'Network error, please try again later!', 1500);
        });
}


function logout() {

    fetch('/logout', {
        method: 'POST',
    })
        .then(async response => {
            if (response.status === 200) {
                document.documentElement.innerHTML = await response.text()
            } else {
                console.log("errrrrror");

            }
        })
        .catch(error => {
            writeError(logerror, "red", 'Network error, please try again later!', 1500);
        });
}

const displayMobileNav = (e) => {
    const nav = document.querySelector('.mobile-nav')
    nav.style.display = 'block'
}

const closeMobileNav = (e) => {
    const nav = document.querySelector('.mobile-nav')
    nav.style.display = 'none'
}

function writeError(targetDiv, color, errormsg, delay) {
    targetDiv.innerText = errormsg
    targetDiv.style.color = color
    setTimeout(() => {
        targetDiv.innerText = ''
    }, delay)
}

let data = true
async function refetch(request) {


    let re = true
    await fetch(request).then(resp => {
        let redirect = resp.headers.get('Location') == '/login';
        if (resp.ok && !redirect) {
            return resp.text()

        } else if (resp.status == 401 || redirect) {
            refetchLogin('/login')
            return
        }
    })
        .then(html => {
            data = true
            let dom = new DOMParser().parseFromString(html, 'text/html')
            if (document.querySelector('.container') && dom.querySelector('.container')) {
                document.querySelector('.container').innerHTML = dom.querySelector('.container').innerHTML
                if (document.querySelector('.next')) {
                    document.querySelector('.next').setAttribute('name', request)
                    document.querySelector('.back').setAttribute('name', request)
                }
            }

        })
        .catch(err => {
            console.log(err);

            data = false
            re = false
        }

        )
    switch (true) {
        case request.includes("/category"):
            document.querySelector('.currentPage').innerText = cats[request[10]];
            break;
        case request.includes("mycreatedposts"):
            document.querySelector('.currentPage').innerText = "My Posts";
            break;
        case request.includes("mylikedposts"):
            document.querySelector('.currentPage').innerText = "My Liked Posts";
            break;
        default:
            document.querySelector('.currentPage').innerText = "Home";
            break;
    }
    return re
}

var PageID = 1

async function pagination(dir) {

    if (dir === "next" && data) {
        let path = document.querySelector('.next').name
        let index = path.indexOf('?')
        if (index != -1) {
            path = path.slice(0, index)
            console.log(path);
        } else {
            PageID = 1
        }

        // const page = +document.querySelector(".currentpage").innerText + 1
        PageID++
        let er = await refetch(`${path}?PageID=${PageID}`)
        if (!er) {
            PageID--
        }


    }

    if (dir === "back" && PageID > 1) {
        let path = document.querySelector('.back').name
        // const page = +document.querySelector(".currentpage").innerText - 1
        let index = path.indexOf('?')
        if (index != -1) {
            path = path.slice(0, index)
            console.log(path);
        }
        PageID--
        console.log(path);

        refetch(`${path}?PageID=${PageID}`)
    }
    console.log(PageID);
}


function selectCat(e) {
    // Parse the value as JSON to extract id and label
    const selectedValue = JSON.parse(e.target.value);
    const { id, label } = selectedValue;

    // create the elemenet for the category
    const span = document.createElement('span');
    span.textContent = label;
    span.classList.add('selected-category');

    // Add a remove button to the span
    const removeBtn = document.createElement('span');
    removeBtn.textContent = '×';
    removeBtn.classList.add('remove-category');
    removeBtn.addEventListener('click', () => {
        span.remove();
        input.remove();
        // Re-enable the corresponding option in the select
        Array.from(e.target.options).find(option => {
            try {
                const optionValue = JSON.parse(option.value);
                return optionValue.id === id;
            } catch {
                return false;
            }
        }).disabled = false;
    });

    span.appendChild(removeBtn);

    // create hidden input to hold the id of selected category
    const input = document.createElement('input')
    input.type = 'hidden';
    input.value = id
    input.name = 'categories'

    // add the elements (span and hidden input) 
    // at the first position of the categories container
    const categoriesContainer = document.querySelector('.selected-categories');
    categoriesContainer.append(input, span);

    // disable the option selected in the select
    e.target.options[e.target.selectedIndex].disabled = true;

    // Reset the select 
    e.target.selectedIndex = 0;

}


async function refetchLogin(request) {
    fetch(request).then(resp => resp.text())
        .then(html => {
            document.documentElement.innerHTML = html
        })
}