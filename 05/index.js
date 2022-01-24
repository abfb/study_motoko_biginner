import { hello } from "../../declarations/hello";

// 发布自己的消息
async function post() {
  let error = document.getElementById("error");
  error.innerText = "";
  let post_button = document.getElementById("post");
  post_button.disabled = true;
  let textarea = document.getElementById("message");
  let opt = document.getElementById("otp").value;
  let text = textarea.value;
  try {
    await hello.post(opt, text);
    textarea.value = "";
  } catch (err) {
    console.log(err);
    error.innerText = "Post Failed!";
  }
  post_button.disabled = false;
}

// 加载自己的消息列表
var num_posts  = 0;
async function load_posts() {
  let posts = await hello.posts(new Date().getTime() * 1000000);
  if (num_posts == posts.length) {
    return;
  }

  let posts_section = document.getElementById("posts");
  posts_section.replaceChildren([]);
  num_posts = posts.length;
  for (var i = 0; i < posts.length; i++) {
    let post = document.createElement("p");
    post.innerText = posts[i].text;
    posts_section.appendChild(post);
  }
}

// 加载关注者列表
var num_followings  = 0;
async function load_followings() {
  let followings = await hello.allFollowings();
  if (num_followings == followings.length) {
    return;
  }

  let followings_section = document.getElementById("followings");
  followings_section.replaceChildren([]);
  num_followings = followings.length;
  for (var i = 0; i < followings.length; i++) {
    let follow = document.createElement("p");
    follow.innerText = followings[i];
    followings_section.appendChild(follow);
  }
}

// 加载关注者消息列表
async function load_msg_canister() {
  let error = document.getElementById("error");
  error.innerText = "";

  let canisterID = document.getElementById("canisterID").value;
  if (canisterID == null) {
    alert("canisterID is null!")
    return;
  }

  try {
    let followingMessages = await hello.followingMessages(canisterID);
    if (followingMessages.length == 0) {
      alert("No Message!")
      return;
    }

    let followingMessages_section = document.getElementById("followingMessages");
    followingMessages_section.replaceChildren([]);
    for (var i = 0; i < followingMessages.length; i++) {
      let followMsg = document.createElement("p");
      followMsg.innerText = "Author: " + followingMessages[i].author 
          + ",\tText: " + followingMessages[i].text 
          + ",\tTime: " + followingMessages[i].time + "\r\n";
      followingMessages_section.appendChild(followMsg);
    }
  } catch (err) {
    console.log(err);
    error.innerText = "Load Message Failed!";
  }
}

// 初始加载
function load() {
  let post_button = document.getElementById("post");
  post_button.onclick = post;
  let loadFollowingMessage = document.getElementById("loadFollowingMessage");
  loadFollowingMessage.onclick = load_msg_canister;
  load_posts();
  setInterval(load_posts, 3000);
  load_followings();
  setInterval(load_followings, 3000);
}

window.onload = load

