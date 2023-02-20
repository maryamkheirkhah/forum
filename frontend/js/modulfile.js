// Hiding stuff (check userName)
const userName = document.getElementById("UserName")
const USER = document.getElementById("USER")
const GUEST = document.getElementById("GUEST")
export function checkUser() {
  const userValue = userName.innerHTML;

  if (userValue === "guest") {
    let Post_link = document.getElementById("post_link")
    Post_link.style.display = "none"
    USER.style.display = "none";
    USER.remove()
  } else {
    GUEST.style.display = "none";
    GUEST.remove()
  }
}

const Comment_box = document.getElementById("comment_box")
const commentBox = document.getElementById("commentBox")
const likeButton = document.getElementById("likeButton")
const dislikeButton = document.getElementById("dislikeButton")
const commentReaction = document.getElementById("commentReaction")
export function checkUser1() {
  const userValue = userName.innerHTML;
  if (userValue === "guest") {
    likeButton.disabled = true
    dislikeButton.disabled = true
    USER.style.display = "none";
    Comment_box.remove()
    commentBox.remove()
    USER.remove()
  } else {
    likeButton.disabled = false
    dislikeButton.disabled = false
    GUEST.style.display = "none";
    GUEST.remove()
  }
}
export function checkUser3() {
  const userValue = userName.innerHTML;

  if (userValue === "guest") {
    USER.style.display = "none";
    USER.remove()
  } else {
    GUEST.style.display = "none";
    GUEST.remove()
  }
}


// check comment status
export function checkCommentValue(id) {

  let CommentId = document.getElementById("CommentId" + id)
  let likeButton = document.getElementById("like" + id)
  let dislikeButton = document.getElementById("dislike" + id)
  let ReactionBox = document.getElementById("Reaction")

  let commentValue = CommentId.innerHTML
  const userValue = userName.innerHTML;
  
  if (commentValue === "1") {
    //likeButton.style.pointerEvents = "none";
    likeButton.classList.toggle("active")
    //dislikeButton.style.pointerEvents = "auto";
    dislikeButton.style.opacity = 1;
  } else if (commentValue === "-1") {
    //  dislikeButton.style.pointerEvents = "none";
    dislikeButton.classList.toggle("active")
    // likeButton.style.pointerEvents = "auto";
    likeButton.style.opacity = 1;
  }
  if (userValue === "guest"){
    ReactionBox.remove()
  }
}

export function alertMessage(str){
  if (str.length != 0){
  window.onload = function() {
    alert(str);
  };
}
}
