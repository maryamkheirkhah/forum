
// comment button 

function Post_comment() {
  var Comment_box = document.getElementById("comment_box");
  if (Comment_box.style.display == "none") {
    Comment_box.style.display = "block";
  } else {
    Comment_box.style.display = "none"
  }
}

// Post Like button
const postValue = document.getElementById("post_like_value")
const likeButton = document.getElementById("likeButton")
const dislikeButton = document.getElementById("dislikeButton")

function checkPostValue() {
  const pValue = postValue.innerHTML;
  console.log(pValue)
  if (pValue === "1") {
    //likeButton.style.pointerEvents = "none";
    likeButton.classList.toggle("active")
    //dislikeButton.style.pointerEvents = "auto";
    dislikeButton.style.opacity = 1;
  } else if (pValue === "-1") {
    //  dislikeButton.style.pointerEvents = "none";
    dislikeButton.classList.toggle("active")
    // likeButton.style.pointerEvents = "auto";
    likeButton.style.opacity = 1;
  }
}


checkPostValue()

