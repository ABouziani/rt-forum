package utils

var HtmlTemplates = map[string]string{
	"header": `{{define "header"}}<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <link rel="stylesheet" href="/assets/css/app.css">
    <link rel="stylesheet" href="/assets/css/navbar.css">
    <link rel="stylesheet" href="/assets/css/login.css">
    <link rel="stylesheet" href="/assets/css/chat.css">
    <link rel="stylesheet" href="/assets/css/post.css">
    <link rel="stylesheet" href="/assets/css/error.css">
    <link rel="shortcut icon" href="/assets/images/favicon.ico" type="image/x-icon">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0-beta3/css/all.min.css">
</head>

<body>
    <header>
        <a class="z01forum" onclick="refetch('/')">
            <p>
                <img class="img-01" src="/assets/images/01.png" alt="01">
                <span class="forum-title">forum</span>
            </p>
        </a>
        {{ if .IsAuthenticated}}
        <div class="header-user">
            <span class="header-username"><i class="fa-regular fa-user"></i>{{.UserName}}</span>

            <button onclick="logout()" class="logout-link">
                <i class="fa-solid fa-right-from-bracket"></i>
            </button>

        </div>
        {{end}}
    </header>
    {{end}}`,

	"footer": `{{define "footer"}}</div>
<div class="footer-container">
    <div class="hidden-footer"></div>
    <footer>
        <div>
            <span>01</span>Forum copyright © 2025
        </div>
        <div class="footer-team">
            <a href="https://github.com/ABouziani"><img src="/assets/images/github.png" alt="github-logo" /> Abdelhamid
                Bouziani</a> |
            <a href="https://github.com/2001basta"><img src="/assets/images/github.png" alt="github-logo" /> Youssef
                Basta</a> |
        </div>
    </footer>
</div>
<script src="/assets/js/index.js"></script>
<script src="/assets/js/ws.js"></script>
</body>

</html>
{{end}}`,

	"navbar": `{{define "navbar"}}<nav>
    <ul class="nav-list">
        <li><a onclick="refetch('/')"><i class="fa-solid fa-house"></i>Home</a></li>
        {{ if .IsAuthenticated}}
        <li><a onclick="refetch('/mycreatedposts')"><i class="fa-regular fa-star"></i></i>My Posts</a></li>
        <li><a onclick="refetch('/mylikedposts')"><i class="fa-regular fa-heart"></i></i>Liked Posts</a></li>
        {{end}}
        <li>
            <span class="categories-title"><i class="fa-solid fa-list"></i>Categories</span>
            {{if .Categories}}
            <ul class="categories-list">
                {{range .Categories}}
                <li><a onclick="refetch('/category/{{.ID}}')">#{{.Label}}</a></li>
                {{end}}
            </ul>
            {{else}}
            <p class="no-categories">No categories available.</p>
            {{end}}
        </li>
        <li><a><i class="fa-regular fa-comments"></i></i>Chat</a></li>
        <ul id="chat-section">
    
        </ul>
    </ul>
</nav>

<!-- mobile nav -->

<nav class="mobile-nav">
    <button class="close-nav" onclick="closeMobileNav()">
        <i class="fa-solid fa-xmark"></i>
    </button>
    <ul class="nav-list">
        <li><a onclick="refetch('/')"><i class="fa-solid fa-house"></i>Home</a></li>
        {{ if .IsAuthenticated}}
        <li><a onclick="refetch('/mycreatedposts')"><i class="fa-regular fa-star"></i></i>My Posts</a></li>
        <li><a onclick="refetch('/mylikedposts')"><i class="fa-regular fa-heart"></i></i>Liked Posts</a></li>
        {{end}}
        <li>
            <span class="categories-title"><i class="fa-solid fa-list"></i>Categories</span>
            {{if .Categories}}
            <ul class="categories-list">
                {{range .Categories}}
                <li><a onclick="refetch('/category/{{.ID}}','cat')">#{{.Label}} ({{.PostsCount}})</a></li>
                {{end}}
            </ul>
            {{else}}
            <p class="no-categories">No categories available.</p>
            {{end}}
        </li>
        <li><a><i class="fa-regular fa-star"></i></i>Chat</a></li>
        <ul id="chat-mobile">
            
        </ul>
    </ul>
</nav>
{{end}}`,

	"error": `
    {{ template "header" . }}
<div class="error-page">
    <div class="error-card">
        <h1>{{.Data.Code}}</h1>
        <p>{{.Data.Message}}</p>
        <a href="/"><i class="fa-solid fa-circle-left"></i> Back Home</a>
    </div>
</div>
{{ template "footer" }}

`,

	"home": `

{{ template "header" . }}
 {{ template "navbar" . }}
<div class="container">
    <div class="alert-content">
        <p id="alert-message"></p>
    </div>
    <div class="posts">
        <p class="currentPage">Home</p>
        <div class="posts-header">
            <button class="nav-button" onclick="displayMobileNav()">
                <i class="fa-solid fa-bars"></i>
            </button>
            <a onclick="refetch('/post/create')" class="create-post-link">
                <i class="fa-solid fa-plus"></i>
                Create post
            </a>
        </div>
        {{if .Data}}
        {{range .Data}}
        <div class="post">
            <div class="post-body">
                <a onclick="refetch('/post/{{.ID}}')" style="cursor: pointer;" class="post-title">{{.Title}}</a>
                <div class="post-header">
                    <p class="post-user">{{.UserName}} </p>
                    <span></span>
                    <p class="post-time" data-timestamp="{{.CreatedAt}}">{{.CreatedAt}}</p>
                </div>
                <pre class="post-content" id="post-content-home">{{.Content}} </pre>
                <div class="post-categories">
                    {{range .Categories}}
                    <span class="post-category">#{{.}}</span>
                    {{end}}
                </div>
            </div>
            <div class="post-footer">
                <button id="likescount{{.ID}}" onclick="postreaction('{{.ID}}','like')"
                    class="post-like post-footer-hover"><i class="fa-regular fa-thumbs-up"></i>{{.Likes}}</button>
                <button id="dislikescount{{.ID}}" onclick="postreaction('{{.ID}}','dislike')"
                    class="post-dislike post-footer-hover"><i
                        class="fa-regular fa-thumbs-down"></i>{{.Dislikes}}</button>
                <a onclick="refetch('/post/{{.ID}}')" class="post-comments post-footer-hover">
                    <i class="fa-regular fa-comment"></i>{{.Comments}}
                </a>
            </div>
            <span style="color:red" id="errorlogin{{.ID}}"></span>
        </div>
        {{end}}
        {{else}}
        <p class="no-posts">No posts available to display !</p>
        {{end}}
    </div>
    <div class="pagination">
        <a onclick="pagination('back', '{{if .Data}}true{{end}}')"  class="back" name="/">&laquo;
            Back</a>
        <span class="currentpage">1</span>
        <a onclick="pagination('next', '{{if .Data}}true{{end}}')" class="next" name="/">Next
            &raquo;</a>
    </div>
</div>
</div>
{{ template "footer" }}
`,

	"login": `
    {{ template "header" . }}
<div class="login-page">
    <div class="login-card">
        <h1>Log in</h1>
        <div class="login-form">
            <input type="text" name="username" id="username" class="login-input" placeholder="email or username">
            <input type="password" name="password" id="password" class="login-input" placeholder="********">
            <span class="errorarea" style="color: rgb(255, 0, 0);"></span>
            <button onclick="login()" class="login-submit">Log in<i class="fa-solid fa-right-to-bracket"></i></button>
        </div>
        <span class="form-line"></span>
        <p class="form-second-option">New to 01Forum? <a style="cursor: pointer; color: blue;" onclick="refetchLogin('/register')">Register</a></p>
    </div>
</div>
{{ template "footer" }}
`,

	"post-form": `
    {{ template "header" . }}
 {{ template "navbar" . }}
<div class="container">
    <div class="alert-content">
        <p id="alert-message"></p>
    </div>
    <div class="create-post">
        <h1>Create a New Post</h1>
        <div class="create-post-fields">
            <label>Title* <span class="max-char">(max: 100 char)</span></label>
            <input name="title" class="create-post-title" placeholder="Enter your post title here..." />
        </div>
        <div class="create-post-fields">
            <label>Categories*</label>
            <div class="create-post-categories">
                <div class="selected-categories"></div>
                <select onchange="selectCat(event)" id="categories-select">
                    <option selected disabled>Select a category</option>
                    {{range .Categories}}
                    <option value='{"id":"{{.ID}}","label":"{{.Label}}"}'>{{.Label}}</option>
                    {{end}}
                </select>
            </div>
        </div>
        <div class="create-post-fields">
            <label>Content* <span class="max-char">(max: 1000 char)</span></label>
            <textarea class="content" name="content" placeholder="What's on your mind?"></textarea>
        </div>
        <span class="errorarea"></span>
        <button id="create-post-btn" onclick="CreatPost()">
            Publish Post
            <i class="fa-regular fa-calendar-plus" id="publish-post-icon"></i>
            <i class="fa-solid fa-circle-notch fa-spin" id="publish-post-circle"></i>
        </button>
    </div>
</div>
{{ template "footer" }}
`,

	"post": `
    {{ template "header" . }}
 {{ template "navbar" . }}
<div class="container">
    <div class="alert-content">
        <p id="alert-message"></p>
    </div>
    <div class="post-detail">
        <div class="post">
            <div class="post-body">
                <p class="post-title">{{.Data.Post.Title}} </p>
                <div class="post-header">
                    <p class="post-user">{{.Data.Post.UserName}} </p>
                    <span></span>
                    <p class="post-time" data-timestamp="{{.Data.Post.CreatedAt}}">{{.Data.Post.CreatedAt}}</p>
                </div>
                <pre class="post-content">{{.Data.Post.Content}} </pre>
                <div class="post-categories">
                    {{range .Data.Post.Categories}}
                    <span class="post-category">#{{.}}</span>
                    {{end}}
                </div>
            </div>
            <div class="post-footer">
                <button id="likescount{{.Data.Post.ID}}" onclick="postreaction('{{.Data.Post.ID}}','like')"
                    class="post-like post-footer-hover"><i
                        class="fa-regular fa-thumbs-up"></i>{{.Data.Post.Likes}}</button>
                <button id="dislikescount{{.Data.Post.ID}}" onclick="postreaction('{{.Data.Post.ID}}','dislike')"
                    class="post-dislike post-footer-hover"><i
                        class="fa-regular fa-thumbs-down"></i>{{.Data.Post.Dislikes}}</button>
                <span class="post-comments"><i class="fa-regular fa-comment"></i>{{.Data.Post.Comments}}</span>
            </div>
            <span style="color:red; border: none;" id="errorlogin{{.Data.Post.ID}}"></span>

        </div>
        <div class="comment-add">
            <textarea name="postid" hidden>{{.Data.Post.ID}}</textarea>
            <textarea id="comment-content" name="comment" placeholder="Add a comment..." required></textarea>
            <button onclick="addcomment('{{.Data.Post.ID}}')" >Comment</button>
        </div>
        <h2 style="padding-left: 10px;">Comments: </h2>
        <div class="comments">
            {{range .Data.Comments}}
            <div class="comment">
                <div class="comment-header">
                    <p class="comment-user">{{.UserName}}</p>
                    <span></span>
                    <p class="comment-time" data-timestamp="{{.CreatedAt}}">{{.CreatedAt}}</p>
                </div>
                <div class="comment-body">
                    <pre class="comment-content">{{.Content}} </pre>
                </div>
                <div class="comment-footer">
                    <button id="commentlikescount{{.ID}}" onclick="commentreaction('{{.ID}}','like')"
                        class="comment-like"><i class="fa-regular fa-thumbs-up"></i>{{.Likes}}</button>
                    <button id="commentdislikescount{{.ID}}" onclick="commentreaction('{{.ID}}','dislike')"
                        class="comment-dislike"><i class="fa-regular fa-thumbs-down"></i>{{.Dislikes}}</button>
                </div>
                <span style="color:red" id="commenterrorlogin{{.ID}}"></span>
            </div>
            {{end}}
        </div>
    </div>
</div>
{{ template "footer" }}
`,

	"register": `
{{ template "header" . }}
<div class="register-page">
    <div class="register-card">
        <h1>Register</h1>
        <div method="post" class="register-form">
            <input id="firstname" type="text" name="firstname" class="register-input" placeholder="First Name">
            <input id="lastname" type="text" name="lastname" class="register-input" placeholder="Last Name">
            <input id="age" type="text" name="age" class="register-input" placeholder="Age (>=18)">
            <input id="email" type="text" name="email" class="register-input" placeholder="email ">
            <input id="username" type="text" name="username" class="register-input" placeholder="username">
            <input id="password" type="password" name="password" class="register-input" placeholder="********">
            <input id="password-confirmation" type="password" name="password-confirmation" class="register-input" placeholder="********">
            <div id="gender" >
                <input id="male" type="radio" name="gender" checked value="male">
                <label for="male">Male</label>
                <input id="female" type="radio" name="gender" value="female">
                <label for="female">Female</label>
            </div>
            <div class="errorarea"></div>
            <button onclick="register()" type="submit" class="register-submit">Register<i class="fa-solid fa-user-plus"></i></button>
        </div>
        <span class="form-line"></span>
        <p class="form-second-option">Already registered? <a style="cursor: pointer; color: blue;" onclick="refetchLogin('/login')">Log in</a></p>
    </div>
</div>
{{ template "footer" }}
`,
}
